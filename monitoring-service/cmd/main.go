package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"monitoring-service/internal/blockchain/delegation"
	"monitoring-service/internal/blockchain/nft"
	"monitoring-service/internal/database"
	"monitoring-service/internal/handlers"
	"monitoring-service/pkg/config"
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
	nftChecker, err := nft.NewNFTChecker(cfg.RpcURL, cfg.NFTContractAddr)
	if err != nil {
		logger.Fatalf("Failed to initialize NFT checker: %v", err)
	}

	// Initialize delegation registry
	client, err := ethclient.Dial(cfg.RpcURL)
	if err != nil {
		logger.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	delegateRegistry, err := delegation.NewDelegationCaller(common.HexToAddress(cfg.DelegateContractAddr), client)
	if err != nil {
		logger.Fatalf("Failed to initialize delegation registry: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.MongoURI, cfg.MongoDB, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize server
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: setupRouter(db, nftChecker, delegateRegistry),
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

func setupRouter(db *database.Database, nftChecker *nft.NFTChecker, delegateRegistry *delegation.DelegationCaller) http.Handler {
	mux := http.NewServeMux()

	// Add health check endpoint
	mux.HandleFunc("/health", logRequest(handlers.HealthCheck))
	
	// Wrap delegations endpoint with CORS
	mux.Handle("/delegations", enableCors(logRequest(handlers.GetDelegations(db))))
	
	// Wrap clients endpoint with CORS
	mux.Handle("/clients", enableCors(logRequest(handlers.GetClients(db))))

	// NFT check endpoint
	mux.HandleFunc("/check-nft", logRequest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.CheckNFT(db, delegateRegistry)(w, r)
	}))

	// Delegation check endpoint
	mux.HandleFunc("/check-delegation", logRequest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.CheckDelegation(nftChecker, delegateRegistry)(w, r)
	}))

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

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
