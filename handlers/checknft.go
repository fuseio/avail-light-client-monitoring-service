package handlers

import (
	"avail-light-client-monitoring-service/blockchain"
	"avail-light-client-monitoring-service/database"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type CheckNFTRequest struct {
	Address string `json:"address"`
	Owner   string `json:"owner"`
	TokenID string `json:"token_id"` // TokenID as string to handle large numbers
}

type CheckNFTResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func CheckNFT(db *database.Database, nftChecker *blockchain.NFTChecker, delegateRegistry *blockchain.DelegationRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := CheckNFTResponse{
			Status:  "success",
			Message: "Client is registered and owns or has delegation for required NFT",
		}

		var req CheckNFTRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.Address == "" {
			response.Status = "error"
			response.Message = "Address is required"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.TokenID == "" {
			response.Status = "error"
			response.Message = "Token ID is required"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Convert TokenID string to big.Int
		tokenID := new(big.Int)
		tokenID, success := tokenID.SetString(req.TokenID, 10)
		if !success {
			response.Status = "error"
			response.Message = "Invalid Token ID format"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// First check if address is already registered
		exists, err := db.ClientExists(req.Address)
		if err != nil {
			http.Error(w, "Failed to check client registration", http.StatusInternalServerError)
			return
		}

		// Check if address owns NFT with specific token ID
		hasNFT, err := nftChecker.HasNFT(req.Address, tokenID)
		if err != nil {
			fmt.Printf("Error checking NFT ownership: %v\n", err)
			http.Error(w, "Failed to verify NFT ownership", http.StatusForbidden)
			return
		}

		if hasNFT {
			response.Status = "success"
			response.Message = "Address owns NFT"
		}

		// If no direct ownership, check delegations
		if !hasNFT {
			checksumAddr := common.HexToAddress(req.Address)
			contractAddr := nftChecker.GetContractAddress()

			if req.Owner == "" {
				response.Status = "error"
				response.Message = "Address does not have NFT, and to check delegation status, Owner Address is required"
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			ownerAddr := common.HexToAddress(req.Owner)

			// Check ERC1155 delegation if still no access
			if !hasNFT {
				rights := common.HexToHash("0x69706c6963656e73650000000000000000000000000000000000000000000000")
				amount, err := delegateRegistry.CheckDelegateForERC1155(checksumAddr, ownerAddr, contractAddr, tokenID, rights)
				if err == nil && amount != nil && amount.Cmp(big.NewInt(0)) > 0 {
					hasNFT = true
					response.Status = "success"
					response.Message = "Address does not have NFT, but has delegation for required NFT"
				}
			}
		}

		if !hasNFT {
			response.Status = "error"
			response.Message = "Address does not own or have delegation for required NFT"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if !exists {
			if err := db.RegisterClient(req.Address, req.TokenID); err != nil {
				http.Error(w, "Failed to register client", http.StatusInternalServerError)
				return
			}
		} else {
			// Update token ID if client already exists
			if err := db.RegisterClient(req.Address, req.TokenID); err != nil {
				http.Error(w, "Failed to update client token ID", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
