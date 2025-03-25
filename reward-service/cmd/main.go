package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reward-service/internal/database"
	"reward-service/internal/handlers"
	"reward-service/internal/service"
	"reward-service/pkg/config"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[Server] ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}
	
	logger.Printf("Starting reward service with configuration:")
	logger.Printf("MongoDB URI: %s", cfg.MongoURI)
	logger.Printf("MongoDB Database: %s", cfg.MongoDB)
	logger.Printf("Port: %s", cfg.Port)
	logger.Printf("Check NFT Interval: %d minutes", cfg.CheckNFTInterval)

	// Initialize database
	logger.Printf("Connecting to MongoDB...")
	db, err := database.NewDatabase(cfg.MongoURI, cfg.MongoDB, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	logger.Printf("Successfully connected to MongoDB")

	// Initialize reward service
	rewardInterval := time.Duration(cfg.CheckNFTInterval) * time.Minute
	logger.Printf("Initializing reward service with interval: %v", rewardInterval)
	rewardService := service.NewRewardService(db, rewardInterval, logger)
	rewardService.Start()
	defer rewardService.Stop()

	// Initialize server
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: setupRouter(db, logger),
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

func setupRouter(db *database.Database, logger *log.Logger) http.Handler {
	mux := http.NewServeMux()

	// Add health check endpoint
	mux.HandleFunc("/health", logRequest(handlers.HealthCheckWithLogging(logger), logger))

	// Add reward service endpoints
	mux.HandleFunc("/clients", logRequest(handlers.GetClients(db, logger), logger))
	mux.HandleFunc("/rewards/summary", logRequest(handlers.GetRewardsSummary(db, logger), logger))
	mux.HandleFunc("/rewards/latest", logRequest(handlers.GetLatestRewards(db, logger), logger))
	mux.HandleFunc("/clients/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// Handle /clients/{address}/rewards
		if len(path) > 8 && path[len(path)-8:] == "/rewards" {
			logRequest(handlers.GetClientRewards(db, logger), logger)(w, r)
			return
		}
		
		// Handle /clients/{address}/points
		if len(path) > 7 && path[len(path)-7:] == "/points" {
			logRequest(handlers.GetClientPoints(db, logger), logger)(w, r)
			return
		}
		
		// Handle other /clients/ paths
		http.NotFound(w, r)
	})

	return mux
}

func logRequest(next http.HandlerFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a custom response writer to capture the status code
		lrw := newLoggingResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log the request details
		logger.Printf(
			"Request: %s | Path: %s | Status: %d | Duration: %v | IP: %s | User-Agent: %s",
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
