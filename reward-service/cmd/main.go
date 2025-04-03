package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"reward-service/internal/api"
	"reward-service/internal/database"
	"reward-service/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger := log.New(os.Stdout, "[REWARD] ", log.LstdFlags)
	logger.Println("Starting reward service...")

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("MONGO_DB", "rewards")
	
	client, err := connectMongoDB(mongoURI)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	
	db := database.New(client, dbName, logger)
	
	rewardService := service.NewRewardService(db, logger)
	
	rewardHandler := api.NewRewardHandler(db, logger)
	
	mux := http.NewServeMux()
	
	mux.HandleFunc("/rewards/claim-all", rewardHandler.HandleClaimAllRewards)
	
	addr := getEnv("HTTP_ADDR", ":8081")
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	
	go func() {
		logger.Printf("Starting HTTP server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()
	
	rewardHour := getEnvInt("REWARD_HOUR", 15) 
	rewardMinute := getEnvInt("REWARD_MINUTE", 22) 
	logger.Printf("Configuring reward schedule to run at %02d:%02d", rewardHour, rewardMinute)
	rewardService.ScheduleRewardsAt(rewardHour, rewardMinute)
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	<-stop
	logger.Println("Shutting down...")
	
	rewardService.Stop()
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown error: %v", err)
	}
	
	logger.Println("Server stopped.")
}

func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	return client, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}
