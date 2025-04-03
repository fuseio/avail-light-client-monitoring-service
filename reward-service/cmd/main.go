package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reward-service/internal/api"
	"reward-service/internal/config"
	"reward-service/internal/database"
	"reward-service/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "[REWARD] ", log.LstdFlags)
	logger.Println("Starting reward service...")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database, logger)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	rewardService := service.NewRewardService(db, logger)
	rewardHandler := api.NewRewardHandler(db, logger)

	router := mux.NewRouter()
	router.Use(api.LoggingMiddleware(logger))
	router.Use(api.RecoveryMiddleware(logger))
	
	rewardHandler.RegisterRoutes(router)
	
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	server := &http.Server{
		Addr:         cfg.HTTP.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Printf("Starting HTTP server on %s", cfg.HTTP.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	logger.Printf("Configuring reward schedule to run at %02d:%02d", cfg.Reward.Hour, cfg.Reward.Minute)
	rewardService.ScheduleRewardsAt(cfg.Reward.Hour, cfg.Reward.Minute)

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
