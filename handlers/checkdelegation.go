package handlers

import (
	"avail-light-client-monitoring-service/blockchain"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type CheckDelegationRequest struct {
	Address string `json:"address"`
	TokenID string `json:"token_id"`
}

type CheckDelegationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Details struct {
		HasTokenDelegation    bool   `json:"has_token_delegation"`
		HasContractDelegation bool   `json:"has_contract_delegation"`
		HasWalletDelegation   bool   `json:"has_wallet_delegation"`
		ERC1155Amount         string `json:"erc1155_amount,omitempty"`
	} `json:"details"`
}

func CheckDelegation(nftChecker *blockchain.NFTChecker, delegateRegistry *blockchain.DelegationRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CheckDelegationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response := CheckDelegationResponse{
			Status:  "success",
			Message: "Delegation check completed",
		}

		// Validate required fields
		if req.Address == "" {
			response.Status = "error"
			response.Message = "Address is required"
			sendJSON(w, response)
			return
		}

		if req.TokenID == "" {
			response.Status = "error"
			response.Message = "Token ID is required"
			sendJSON(w, response)
			return
		}

		// Convert TokenID string to big.Int
		tokenID := new(big.Int)
		tokenID, success := tokenID.SetString(req.TokenID, 10)
		if !success {
			response.Status = "error"
			response.Message = "Invalid Token ID format"
			sendJSON(w, response)
			return
		}

		checksumAddr := common.HexToAddress(req.Address)
		contractAddr := nftChecker.GetContractAddress()
		var rights [32]byte // Zero rights for basic delegation check

		// Check token-level delegation
		hasTokenDelegation, err := delegateRegistry.CheckDelegateForToken(checksumAddr, contractAddr, contractAddr, tokenID)
		if err != nil {
			fmt.Printf("Error checking token delegation: %v\n", err)
		}
		response.Details.HasTokenDelegation = hasTokenDelegation

		// Check contract-level delegation
		hasContractDelegation, err := delegateRegistry.CheckDelegateForContract(checksumAddr, contractAddr, contractAddr)
		if err != nil {
			fmt.Printf("Error checking contract delegation: %v\n", err)
		}
		response.Details.HasContractDelegation = hasContractDelegation

		// Check wallet-level delegation
		hasWalletDelegation, err := delegateRegistry.CheckDelegateForAll(checksumAddr, contractAddr)
		if err != nil {
			fmt.Printf("Error checking wallet delegation: %v\n", err)
		}
		response.Details.HasWalletDelegation = hasWalletDelegation

		// Check ERC1155 delegation amount
		amount, err := delegateRegistry.CheckDelegateForERC1155(checksumAddr, contractAddr, contractAddr, tokenID, rights)
		if err != nil {
			fmt.Printf("Error checking ERC1155 delegation: %v\n", err)
		} else if amount != nil {
			response.Details.ERC1155Amount = amount.String()
		}

		sendJSON(w, response)
	}
}

func sendJSON(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
