package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type ClaimRewardRequest struct {
	RewardID string `json:"reward_id"`
	Address  string `json:"address"`
}

type ClaimRewardResponse struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	Points       int64     `json:"points,omitempty"`
	TotalPoints  int64     `json:"total_points,omitempty"`
	ClaimedAt    time.Time `json:"claimed_at,omitempty"`
}

type GetUserRewardsResponse struct {
	Rewards      []database.RewardRecord `json:"rewards"`
	TotalPoints  int64                   `json:"total_points"`
	ClaimedPoints int64                  `json:"claimed_points"`
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

func (h *RewardHandler) ClaimReward(w http.ResponseWriter, r *http.Request) {
	var req ClaimRewardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rewardID, err := primitive.ObjectIDFromHex(req.RewardID)
	if err != nil {
		http.Error(w, "Invalid reward ID", http.StatusBadRequest)
		return
	}

	err = h.db.ClaimReward(rewardID, req.Address)
	if err != nil {
		h.logger.Printf("ERROR: Failed to claim reward: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPoints, err := h.db.GetUserClaimedPoints(req.Address)
	if err != nil {
		h.logger.Printf("ERROR: Failed to get user points: %v", err)
	}

	reward, err := h.db.GetRewardByID(rewardID)
	if err != nil {
		h.logger.Printf("ERROR: Failed to get reward details: %v", err)
	}

	resp := ClaimRewardResponse{
		Success:     true,
		Message:     "Reward claimed successfully",
		Points:      reward.Points,
		TotalPoints: totalPoints,
		ClaimedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RewardHandler) GetUserRewards(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	rewards, err := h.db.GetRewardsByAddress(address)
	if err != nil {
		h.logger.Printf("ERROR: Failed to get rewards: %v", err)
		http.Error(w, "Failed to get rewards", http.StatusInternalServerError)
		return
	}

	claimedPoints, err := h.db.GetUserClaimedPoints(address)
	if err != nil {
		h.logger.Printf("ERROR: Failed to get user claimed points: %v", err)
	}

	var totalPoints int64
	for _, reward := range rewards {
		totalPoints += reward.Points
	}

	resp := GetUserRewardsResponse{
		Rewards:       rewards,
		TotalPoints:   totalPoints,
		ClaimedPoints: claimedPoints,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	router.HandleFunc("/rewards/claim", h.ClaimReward).Methods("POST")
	router.HandleFunc("/rewards/user/{address}", h.GetUserRewards).Methods("GET")
	router.HandleFunc("/claim-all-rewards", h.HandleClaimAllRewards).Methods("POST")
}