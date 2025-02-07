package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"monitoring-service/internal/blockchain/delegation"
	"monitoring-service/internal/database"
	"monitoring-service/pkg/config"

	"github.com/ethereum/go-ethereum/common"
)

type CheckNFTRequest struct {
	Address string `json:"address"`
	CommissionRate string `json:"commission_rate"`
}

type CheckNFTResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func updateOwnershipClientRegistration(db *database.Database, address string, totalAmount int64, checkNFTInterval int, commissionRate string) error {
	exists, err := db.ClientExists(address)
	if err != nil {
		return err
	}

	// convert commission rate to float64
	commissionRateFloat, err := strconv.ParseFloat(commissionRate, 64)
	if err != nil {
		return err
	}

	totalTime := 0
	operationPoints := database.OperationPointRecord{
		Amount:         totalAmount,
		Timestamp:     time.Now(),
		CommissionRate: commissionRateFloat,
		Time:           0,
	}
	
	if exists {
		client, err := db.GetClient(address)
		if err != nil {
			return err
		}
		
		if time.Since(client.LastHeartbeat) <= time.Duration(checkNFTInterval) * time.Minute {
			totalTime = int(client.TotalTime) + int(time.Since(client.LastHeartbeat).Seconds())
			// update Time of the operationPoints
			operationPoints.Time = int64(time.Since(client.LastHeartbeat).Seconds())
		} else {
			totalTime = int(client.TotalTime)
		}
	}

	return db.RegisterClient(address, operationPoints, nil, int64(totalTime))
}

func updateDelegationClientRegistration(db *database.Database, address string, totalAmount int64, delegationAddress string, commissionRate string) error {
	// convert commission rate to float64
	commissionRateFloat, err := strconv.ParseFloat(commissionRate, 64)
	if err != nil {
		return err
	}

	exists, err := db.ClientExists(address)
	if err != nil {
		return err
	}

	totalTime := int64(0)
	if exists {
		client, err := db.GetClient(address)
		if err != nil {
			return err
		}
		totalTime = client.TotalTime
	}

	delegationPoints := &database.DelegationPointRecord{
		Address:        delegationAddress,
		Amount:         totalAmount,
		Timestamp:      time.Now(),
		CommissionRate: commissionRateFloat,
	}
	return db.RegisterClient(address, database.OperationPointRecord{}, delegationPoints, totalTime)
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

		if req.CommissionRate == "" {
			response.Status = "error"
			response.Message = "Commission rate is required"
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
				if delegationRights == configRights && delegation.Contract == common.HexToAddress(cfg.NFTContractAddr) && delegation.Type == 5 {
					delegations = append(delegations, delegation)
				}
			}

			// get map From address to token id and amount
			tokenIdMap := make(map[string]int64)
			for _, delegation := range delegations {
				tokenIdMap[delegation.From.String()] = tokenIdMap[delegation.From.String()] + delegation.Amount.Int64()
			}
			
			// sum up all the amounts
			var totalAmount int64
			for _, amount := range tokenIdMap {
				totalAmount += amount
			}

			if totalAmount > 0 {
				response.Status = "success"
				response.Message = "Address has NFT or delegation for required NFT"

				if err := updateOwnershipClientRegistration(db, req.Address, totalAmount, cfg.CheckNFTInterval, req.CommissionRate); err != nil {
					http.Error(w, "Failed to update client registration", http.StatusInternalServerError)
					return
				}

				// update delegation client registration
				for key, amount := range tokenIdMap {
					if key == req.Address { continue; }
					if err := updateDelegationClientRegistration(db, key, amount, req.Address, req.CommissionRate); err != nil {
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
