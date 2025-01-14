package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avail-light-client-monitoring-service/blockchain"
	"avail-light-client-monitoring-service/config"
	"avail-light-client-monitoring-service/database"
	"avail-light-client-monitoring-service/handlers"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[Server] ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize NFT checker
	nftChecker, err := blockchain.NewNFTChecker(cfg.RpcURL, cfg.NFTContractAddr)
	if err != nil {
		logger.Fatalf("Failed to initialize NFT checker: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.DBPath, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize server
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: setupRouter(db, nftChecker),
	}

	// Start server
	go func() {
		logger.Printf("Starting server on port %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited properly")
}

func setupRouter(db *database.Database, nftChecker *blockchain.NFTChecker) http.Handler {
	mux := http.NewServeMux()

	// Add health check endpoint
	mux.HandleFunc("/health", logRequest(handlers.HealthCheck))

	// Add register endpoint
	mux.HandleFunc("/check-nft", logRequest(handlers.CheckNFT(db, nftChecker)))

	return mux
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a custom response writer to capture the status code
		lrw := newLoggingResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log the request details
		log.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | IP: %s | User-Agent: %s",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(startTime),
			r.RemoteAddr,
			r.UserAgent(),
		)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
