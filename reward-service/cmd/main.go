package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"reward-service/internal/api"
	"reward-service/internal/database"
	"reward-service/internal/service"
)

func main() {
	logger := log.New(os.Stdout, "[REWARD-SERVICE] ", log.LstdFlags)
	logger.Println("Starting reward service...")

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "lc-monitoring"
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	
	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		logger.Fatalf("Failed to ping MongoDB: %v", err)
	}
	
	logger.Printf("Connected to MongoDB: %s", mongoURI)
	
	// Create database wrapper
	db := client.Database(mongoDB)
	dbWrapper := database.NewDatabase(client, db, logger)
	
	// Create database indexes
	if err := dbWrapper.AddIndexes(); err != nil {
		logger.Printf("Warning: Failed to create database indexes: %v", err)
	}
	
	// Create API handlers
	rewardHandler := api.NewRewardHandler(dbWrapper, logger)
	
	// Set up router
	router := mux.NewRouter()
	rewardHandler.RegisterRoutes(router)
	
	// Add health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Start reward service
	rewardService := service.NewRewardService(db, logger)
	rewardService.ScheduleRewards()
	
	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if port[0] == ':' {
		port = port[1:]
	}
	
	// Create and start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	
	// Start server in a goroutine
	go func() {
		logger.Printf("Server listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	logger.Println("Shutting down...")
	
	// Stop the reward scheduler
	rewardService.Stop()
	
	// Shut down the server
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}
	
	// Close MongoDB connection
	if err := client.Disconnect(ctx); err != nil {
		logger.Fatalf("MongoDB disconnect failed: %v", err)
	}
	
	logger.Println("Server gracefully stopped")
}
