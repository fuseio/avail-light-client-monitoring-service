package handlers

import (
	"avail-light-client-monitoring-service/blockchain"
	"avail-light-client-monitoring-service/database"
	"encoding/json"
	"net/http"
)

type CheckNFTRequest struct {
	Address string `json:"address"`
}

type CheckNFTResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func CheckNFT(db *database.Database, nftChecker *blockchain.NFTChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CheckNFTRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		// First check if address is already registered
		exists, err := db.ClientExists(req.Address)
		if err != nil {
			http.Error(w, "Failed to check client registration", http.StatusInternalServerError)
			return
		}

		// Check if address owns NFT
		hasNFT, err := nftChecker.HasNFT(req.Address)
		if err != nil {
			http.Error(w, "Failed to verify NFT ownership", http.StatusInternalServerError)
			return
		}

		if !hasNFT {
			http.Error(w, "Address does not own required NFT", http.StatusForbidden)
			return
		}

		if !exists {
			if err := db.RegisterClient(req.Address); err != nil {
				http.Error(w, "Failed to register client", http.StatusInternalServerError)
				return
			}
		}

		response := CheckNFTResponse{
			Status:  "success",
			Message: "Client is registered and owns required NFT",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
