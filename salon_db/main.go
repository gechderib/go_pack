package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	// db_postgres_url := os.Getenv("DB_POSTGRES_URL")
	db_postgres_url := LoadConfig().DBURL
	db, err = gorm.Open(postgres.Open(db_postgres_url), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	// db.AutoMigrate(&User{}, &Order{})
}

// loadEnv loads environment variables from .env file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

}

var logger *zap.Logger

func initLogger() {
	var err error

	// logger, err = zap.NewProduction()
	// if err != nil {
	// 	panic("failed to initialize logger")
	// }

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"logs/app.log", "stdout"}
	logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger")
	}

}

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func main() {
	// load environment variables
	loadEnv()

	// config
	config := LoadConfig()

	// initialize database
	initDB()

	// initialize logger
	initLogger()
	defer logger.Sync()

	// for pprof testing
	go func() {
		logger.Info("pprof running on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logger.Error("pprof error", zap.Error(err))
		}
	}()

	// register Prometheus metrics
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)

	// initialize routes
	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	// public routes
	r.Group(func(r chi.Router) {
		r.Use(LoggingMiddleware)
		r.Use(RecoveryMiddleware)
		r.Use(AuthMiddleware)
		r.Use(TimeoutMiddleware)

		r.Get("/users", GetUsers)
		r.Get("/users/{id}", GetUserById)
	})

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(LoggingMiddleware)
		r.Use(RecoveryMiddleware)
		r.Use(AuthMiddleware)

		r.Post("/users", CreateUser)
		r.Delete("/users/{id}", DeleteUser)
		r.Post("/order", CreateOrder)
		r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("This is a test panic")
		})
	})

	// http.ListenAndServe(":"+config.PORT, r) // IGNORE -- Graceful shutdown

	// Graceful shutdown
	serve := http.Server{
		Addr:    ":" + config.PORT,
		Handler: r,
	}

	go func() {
		if err := serve.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serve.Shutdown(ctx); err != nil {
		logger.Fatal("Failed to shutdown server", zap.Error(err))
	}

}
