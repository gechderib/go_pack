package main

import (
	"fmt"
	"net/http"
	"time"

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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, 200}

		// logger.Info("incoming request",
		// 	zap.String("method", r.Method),
		// 	zap.String("path", r.URL.Path),
		// )

		next.ServeHTTP(rw, r)
		logger.Info("completed request",
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

		// Optional: check Bearer format
		// "Bearer token"
		// parts := strings.Split(authHeader, " ")
		// if len(parts) != 2 || parts[0] != "Bearer" { ... }

		next.ServeHTTP(w, r)
	})
}
