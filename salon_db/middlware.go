package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		fmt.Println("➡️", r.Method, r.URL.Path)

// 		next.ServeHTTP(w, r)

// 		fmt.Println("⬅️ Completed in", time.Since(start))

//		})
//	}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		rw := &responseWriter{w, 200}

// 		requestId := uuid.New().String()
// 		ctx := context.WithValue(r.Context(), "request_id", requestId)
// 		r = r.WithContext(ctx)

// 		next.ServeHTTP(rw, r)
// 		// 🔥 IMPORTANT: use updated context AFTER next
// 		userId, _ := r.Context().Value("user_id").(string)

// 		logger.Info("completed request",
// 			zap.String("request_id", requestId),
// 			zap.String("user_id", userId),
// 			zap.String("method", r.Method),
// 			zap.String("path", r.URL.Path),
// 			zap.Int("status", rw.statusCode),
// 			zap.Duration("duration", time.Since(start)),
// 		)
// 	})
// }

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, 200}

		requestId := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request_id", requestId)
		r = r.WithContext(ctx)

		duration := time.Since(start).Seconds()
		// 🔥 Prometheus metrics
		requestCount.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(rw.statusCode),
		).Inc()

		requestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(duration)

		next.ServeHTTP(rw, r)
		// 🔥 IMPORTANT: use updated context AFTER next
		userId, _ := r.Context().Value("user_id").(string)

		logger.Info("completed request",
			zap.String("request_id", requestId),
			zap.String("user_id", userId),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("🚨 Recovered from panic:", err)
				// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				JSONResponse(w, http.StatusInternalServerError, APIResponse{
					Success: false,
					Message: "Internal server Error bro",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			JSONResponse(w, http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Unauthorized",
				Errors:  "Missing Authorization header",
			})
			return
		}
		// consider this is from a token parsing logic
		userId := uuid.New().String()
		ctx := context.WithValue(r.Context(), "user_id", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
