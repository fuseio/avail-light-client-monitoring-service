package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"reward-service/internal/database"
)

type RewardHandler struct {
	db     *database.Database
	logger *log.Logger
}

func NewRewardHandler(db *database.Database, logger *log.Logger) *RewardHandler {
	return &RewardHandler{
		db:     db,
		logger: logger,
	}
}

type ClaimAllRewardsRequest struct {
	Address string `json:"address"`
}

type ClaimAllRewardsResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ClaimedCount int    `json:"claimed_count"`
	TotalPoints  int64  `json:"total_points"`
}

func (h *RewardHandler) HandleClaimAllRewards(w http.ResponseWriter, r *http.Request) {
	var req ClaimAllRewardsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: Failed to parse claim all rewards request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Address == "" {
		h.logger.Printf("ERROR: Missing address in claim all rewards request")
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	claimedCount, totalPoints, err := h.db.ClaimAllRewards(req.Address)
	if err != nil {
		h.logger.Printf("ERROR: Failed to claim all rewards for %s: %v", req.Address, err)
		http.Error(w, fmt.Sprintf("Failed to claim rewards: %v", err), http.StatusInternalServerError)
		return
	}

	response := ClaimAllRewardsResponse{
		Success:      true,
		ClaimedCount: claimedCount,
		TotalPoints:  totalPoints,
	}

	if claimedCount == 0 {
		response.Message = "No unclaimed rewards found"
	} else {
		response.Message = fmt.Sprintf("Successfully claimed %d rewards for a total of %d points", claimedCount, totalPoints)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("ERROR: Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("SUCCESS: Claimed %d rewards with %d points for %s", claimedCount, totalPoints, req.Address)
}

func (h *RewardHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/rewards/claim-all", h.HandleClaimAllRewards).Methods("POST")
}