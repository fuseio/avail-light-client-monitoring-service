package api

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create response wrapper to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(rw, r)
			
			duration := time.Since(start)
			
			logger.Printf(
				"[%s] %s %s %d %s",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				rw.statusCode,
				duration,
			)
		})
	}
}

func RecoveryMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("PANIC: %v\n%s", err, debug.Stack())
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
