package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"reward-service/internal/database"
	"github.com/gorilla/mux"
)

type GetClientsResponse struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []database.MonitoringClientInfo `json:"data"`
}

type GetClientRewardsResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []database.RewardRecord `json:"data"`
}

type GetClientPointsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Points  int64  `json:"points"`
}

type GetRewardsSummaryResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type GetLatestRewardsResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []database.RewardRecord `json:"data"`
	Count   int                    `json:"count"`
	Total   int64                  `json:"total_points"`
}

func GetClients(db *database.Database, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("API Request: GET /clients from %s", r.RemoteAddr)
		
		if r.Method != http.MethodGet {
			logger.Printf("Method not allowed: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		users, err := db.GetAllUsers()
		if err != nil {
			logger.Printf("Error fetching clients: %v", err)
			http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
			return
		}

		response := GetClientsResponse{
			Status:  "success",
			Message: "Clients retrieved successfully",
			Data:    users,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		
		duration := time.Since(startTime)
		logger.Printf("API Response: GET /clients - %d users returned in %v", len(users), duration)
	}
}

func GetClientRewards(db *database.Database, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("API Request: GET %s from %s", r.URL.Path, r.RemoteAddr)
		
		if r.Method != http.MethodGet {
			logger.Printf("Method not allowed: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	
		vars := mux.Vars(r)
		address := vars["address"]
		
		if address == "" {
			logger.Printf("Bad request: missing address parameter")
			http.Error(w, "Address parameter is required", http.StatusBadRequest)
			return
		}
		
	
		claimedPoints, err := db.GetUserClaimedPoints(address)
		if err != nil {
			logger.Printf("Error getting user claimed points: %v", err)
		
			claimedPoints = 0
		}
		
	
		response := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
			Points  int64  `json:"points"`
		}{
			Success: true,
			Message: "Points retrieved successfully",
			Points:  claimedPoints,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		
		duration := time.Since(startTime)
		logger.Printf("API Response: GET %s - %d points returned in %v", r.URL.Path, claimedPoints, duration)
	}
}

func GetClientPoints(db *database.Database, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("API Request: GET %s from %s", r.URL.Path, r.RemoteAddr)
		
		if r.Method != http.MethodGet {
			logger.Printf("Method not allowed: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	
		path := strings.TrimPrefix(r.URL.Path, "/clients/")
		path = strings.TrimSuffix(path, "/points")
		address := path

		if address == "" {
			logger.Printf("Error: Missing address in request %s", r.URL.Path)
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		
		logger.Printf("Fetching points for address: %s", address)

	
		claimedPoints, err := db.GetUserClaimedPoints(address)
		if err != nil {
			logger.Printf("Error fetching claimed points for %s: %v", address, err)
		
			claimedPoints = 0
		}

		response := GetClientPointsResponse{
			Status:  "success",
			Message: "Points retrieved successfully",
			Points:  claimedPoints,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		
		duration := time.Since(startTime)
		logger.Printf("API Response: GET %s - %d points returned in %v", r.URL.Path, claimedPoints, duration)
	}
}

func GetRewardsSummary(db *database.Database, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("API Request: GET /rewards/summary from %s", r.RemoteAddr)
		
		if r.Method != http.MethodGet {
			logger.Printf("Method not allowed: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		summary, err := db.GetRewardsSummary()
		if err != nil {
			logger.Printf("Error generating rewards summary: %v", err)
			http.Error(w, "Failed to generate rewards summary", http.StatusInternalServerError)
			return
		}

		response := GetRewardsSummaryResponse{
			Status:  "success",
			Message: "Rewards summary generated successfully",
			Data:    summary,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		
		duration := time.Since(startTime)
		logger.Printf("API Response: GET /rewards/summary - summary returned in %v", duration)
	}
}

func GetLatestRewards(db *database.Database, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("API Request: GET /rewards/latest from %s", r.RemoteAddr)
		
		if r.Method != http.MethodGet {
			logger.Printf("Method not allowed: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
	
		limit := 10
		limitParam := r.URL.Query().Get("limit")
		if limitParam != "" {
			if n, err := strconv.Atoi(limitParam); err == nil && n > 0 {
				limit = n
				if limit > 100 {
					limit = 100
				}
			}
		}
		
		rewards, err := db.GetLatestRewards(limit)
		if err != nil {
			logger.Printf("Error getting latest rewards: %v", err)
			http.Error(w, "Failed to get latest rewards", http.StatusInternalServerError)
			return
		}
		
	
		var totalPoints int64
		for _, reward := range rewards {
			totalPoints += reward.Points
		}
		
		response := GetLatestRewardsResponse{
			Status:  "success",
			Message: "Latest rewards retrieved successfully",
			Data:    rewards,
			Count:   len(rewards),
			Total:   totalPoints,
		}
		
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		
		duration := time.Since(startTime)
		logger.Printf("API Response: GET /rewards/latest - %d rewards returned in %v", len(rewards), duration)
	}
} 