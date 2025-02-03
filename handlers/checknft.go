package handlers

import (
	"avail-light-client-monitoring-service/blockchain/delegation"
	"avail-light-client-monitoring-service/config"
	"avail-light-client-monitoring-service/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type CheckNFTRequest struct {
	Address string `json:"address"`
}

type CheckNFTResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func updateOwnershipClientRegistration(db *database.Database, address string, totalAmount int64, checkNFTInterval int) error {
	exists, err := db.ClientExists(address)
	if err != nil {
		return err
	}

	runtime := 0
	delegationTime := 0
	totalTime := 0

	if exists {
		client, err := db.GetClient(address)
		if err != nil {
			return err
		}
		
		if time.Since(client.LastHeartbeat) <= time.Duration(checkNFTInterval) * time.Minute {
			runtime = int(client.Runtime) + int(totalAmount) * int(time.Since(client.LastHeartbeat).Seconds())
			totalTime = int(client.TotalTime) + int(time.Since(client.LastHeartbeat).Seconds())
		} else {
			runtime = int(client.Runtime)
			totalTime = int(client.TotalTime)
		}
	}

	return db.RegisterClient(address, int64(runtime), int64(delegationTime), int64(totalTime))
}

func updateDelegationClientRegistration(db *database.Database, address string, totalAmount int64, checkNFTInterval int) error {
	exists, err := db.ClientExists(address)
	if err != nil {
		return err
	}

	runtime := 0
	delegationTime := 0
	totalTime := 0

	if exists {
		client, err := db.GetClient(address)
		if err != nil {
			return err
		}
		
		if time.Since(client.LastHeartbeat) <= time.Duration(checkNFTInterval) * time.Minute {
			delegationTime = int(client.DelegationTime) + int(totalAmount) * int(time.Since(client.LastHeartbeat).Seconds())
			totalTime = int(client.TotalTime) + int(time.Since(client.LastHeartbeat).Seconds())
		} else {
			delegationTime = int(client.DelegationTime)
			totalTime = int(client.TotalTime)
		}
	}

	return db.RegisterClient(address, int64(runtime), int64(delegationTime), int64(totalTime))
}

func CheckNFT(db *database.Database, delegateRegistry *delegation.DelegationCaller) http.HandlerFunc {
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

		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			http.Error(w, "Failed to load config", http.StatusInternalServerError)
			return
		}

		incommingDelegation, err := delegateRegistry.GetIncomingDelegations(nil, common.HexToAddress(req.Address))
		if err != nil {
			http.Error(w, "Failed to get incoming delegations", http.StatusInternalServerError)
			return
		}
		if len(incommingDelegation) > 0 {
			// filter by rights and contract address
			var delegations []delegation.IDelegateRegistryDelegation
			for _, delegation := range incommingDelegation {
				delegationRights := [32]byte(delegation.Rights)
				configRights := [32]byte(cfg.Rights)
				if delegationRights == configRights && delegation.Contract == common.HexToAddress(cfg.NFTContractAddr) {
					delegations = append(delegations, delegation)
				}
			}

			// get map From address to token id and amount
			tokenIdMap := make(map[string]struct {
				TokenId string
				Amount  int64
			})
			for _, delegation := range delegations {
				tokenIdMap[delegation.From.String()] = struct {
					TokenId string
					Amount  int64
				}{
					TokenId: delegation.TokenId.String(),
					Amount:  delegation.Amount.Int64(),
				}
			}
			
			// sum up all the amounts
			var totalAmount int64
			for _, amount := range tokenIdMap {
				totalAmount += amount.Amount
			}

			if totalAmount > 0 {
				response.Status = "success"
				response.Message = "Address has NFT or delegation for required NFT"

				if err := updateOwnershipClientRegistration(db, req.Address, totalAmount, cfg.CheckNFTInterval); err != nil {
					http.Error(w, "Failed to update client registration", http.StatusInternalServerError)
					return
				}

				// update delegation client registration
				for key, amount := range tokenIdMap {
					if key == req.Address { continue; }
					if err := updateDelegationClientRegistration(db, key, amount.Amount, cfg.CheckNFTInterval); err != nil {
						http.Error(w, "Failed to update client registration", http.StatusInternalServerError)
						return
					}
				}
			} else {
				response.Status = "error"
				response.Message = "Address does not own or have delegation for required NFT"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
