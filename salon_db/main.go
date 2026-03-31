package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
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

	// initialize routes
	r := chi.NewRouter()

	// public routes
	r.Group(func(r chi.Router) {
		r.Use(LoggingMiddleware)
		r.Use(RecoveryMiddleware)

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

	http.ListenAndServe(":"+config.Port, r)
}
