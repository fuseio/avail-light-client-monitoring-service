package handlers

import (

	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"monitoring-service/internal/blockchain/delegation"
	"monitoring-service/internal/blockchain/nft"
	"monitoring-service/internal/database"
	"monitoring-service/pkg/config"

	"github.com/ethereum/go-ethereum/common"
)

type CheckNFTRequest struct {
	Address              string `json:"address"`
	CommissionRate       string `json:"commission_rate"`
	OperatorName         string `json:"operator_name"`
	RewardCollectorAddress string `json:"reward_collector_address"`
}

type CheckNFTResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func updateOwnershipClientRegistration(db *database.Database, address string, totalAmount int64, checkNFTInterval int, commissionRate string, operatorName string, rewardCollectorAddress string) error {
	exists, err := db.ClientExists(address)
	if err != nil {
		return err
	}

	var commissionRateFloat float64
	var clientRecord *database.ClientInfo

	if exists {
		clientRecord, err = db.GetClient(address)
		if err != nil {
			return err
		}
		// Use the already set commission rate.
		commissionRateFloat = clientRecord.CommissionRate
	} else {
		commissionRateFloat, err = strconv.ParseFloat(commissionRate, 64)
		if err != nil {
			return err
		}
	}

	totalTime := 0
	operationPoints := database.OperationPointRecord{
		Amount:         totalAmount,
		Timestamp:      time.Now(),
		CommissionRate: commissionRateFloat,
		Time:           0,
	}

	if exists {
		totalTime = int(clientRecord.TotalTime) + int(time.Since(clientRecord.LastHeartbeat).Seconds())
		if time.Since(clientRecord.LastHeartbeat) <= time.Duration(checkNFTInterval)*time.Minute {
			operationPoints.Time = int64(time.Since(clientRecord.LastHeartbeat).Seconds())
		} else {
			operationPoints.Time = 0
		}
	}

	return db.RegisterClient(address, operationPoints, int64(totalTime), operatorName, rewardCollectorAddress)
}

func updateDelegationClientRegistration(db *database.Database, address string, totalAmount int64, delegationAddress string, commissionRate string) error {
	// convert commission rate to float64
	commissionRateFloat, err := strconv.ParseFloat(commissionRate, 64)
	if err != nil {
		return err
	}

	delegationPoints := database.DelegationPointRecord{
		Address:        delegationAddress,
		Amount:         totalAmount,
		Timestamp:      time.Now(),
		CommissionRate: commissionRateFloat,
	}
	return db.RegisterDelegation(address, delegationPoints)
}

func CheckNFT(db *database.Database, delegateRegistry *delegation.DelegationCaller, nftChecker *nft.NFTChecker) http.HandlerFunc {
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
			fmt.Println("Validation Error: Address is required")
			response.Status = "error"
			response.Message = "Address is required"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.CommissionRate == "" {
			fmt.Println("Validation Error: Commission rate is required")
			response.Status = "error"
			response.Message = "Commission rate is required"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.OperatorName == "" {
			fmt.Println("Validation Error: Operator name is required")
			response.Status = "error"
			response.Message = "Operator name is required"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.RewardCollectorAddress == "" {
			fmt.Println("Validation Error: Reward collector address is required")
			response.Status = "error"
			response.Message = "Reward collector address is required"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check: commission rate must be between 0 and 10
		commission, err := strconv.ParseFloat(req.CommissionRate, 64)
		if err != nil {
			fmt.Println("Validation Error: Invalid commission rate format")
			response.Status = "error"
			response.Message = "Invalid commission rate format"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		if commission < 0 || commission > 10 {
			fmt.Println("Validation Error: Commission rate must be between 0 and 10")
			response.Status = "error"
			response.Message = "Commission rate must be between 0 and 10"
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
			fmt.Println("Incoming delegations found, processing...")
			var delegations []delegation.IDelegateRegistryDelegation
			delegatorBalances := make(map[string]int64)

			// First collect all delegations and check delegator balances
			for _, delegation := range incommingDelegation {
				delegationRights := [32]byte(delegation.Rights)
				configRights := [32]byte(cfg.Rights)
				if delegationRights == configRights && delegation.Contract == common.HexToAddress(cfg.NFTContractAddr) && delegation.Type == 5 {
					// Check delegator's NFT balance
					balance, err := nftChecker.GetBatchBalance(delegation.From.String(), []*big.Int{delegation.TokenId})
					if err != nil {
						fmt.Printf("Failed to check NFT balance for delegator %s: %v\n", delegation.From.String(), err)
						http.Error(w, "Failed to check NFT balance", http.StatusInternalServerError)
						return
					}

					fmt.Printf("Delegation Check - From: %s, To: %s, TokenID: %s, DelegatedAmount: %d, ActualBalance: %d\n",
						delegation.From.String(),
						delegation.To.String(),
						delegation.TokenId.String(),
						delegation.Amount.Int64(),
						balance[0].Int64())

					if len(balance) > 0 && balance[0].Int64() > 0 {
						delegations = append(delegations, delegation)
						delegatorBalances[delegation.From.String()] += balance[0].Int64()
					}
				}
			}

			// get map From address to token id and amount
			tokenIdMap := make(map[string]int64)
			for _, delegation := range delegations {
				// Only count delegation amount up to delegator's actual balance
				availableBalance := delegatorBalances[delegation.From.String()]
				delegationAmount := delegation.Amount.Int64()
				if delegationAmount > availableBalance {
					delegationAmount = availableBalance
				}
				
				if delegationAmount > 0 {
					tokenIdMap[delegation.To.String()] = tokenIdMap[delegation.To.String()] + delegationAmount
				}
			}

			// sum up all the amounts
			var totalAmount int64
			for _, amount := range tokenIdMap {
				totalAmount += amount
			}

			// Check if this client already exists in the DB
			exists, err := db.ClientExists(req.Address)
			if err != nil {
				http.Error(w, "Failed to check client existence", http.StatusInternalServerError)
				return
			}

			// If client exists OR totalAmount > 0 (new client with non-zero delegation), update the record.
			if exists || totalAmount > 0 {
				if err := updateOwnershipClientRegistration(db, req.Address, totalAmount, cfg.CheckNFTInterval, req.CommissionRate, req.OperatorName, req.RewardCollectorAddress); err != nil {
					http.Error(w, "Failed to update client registration", http.StatusInternalServerError)
					return
				}
				for key, amount := range tokenIdMap {
					if err := updateDelegationClientRegistration(db, req.Address, amount, key, req.CommissionRate); err != nil {
						http.Error(w, "Failed to update delegation registration", http.StatusInternalServerError)
						return
					}
				}
				response.Status = "success"
				response.Message = "Address has NFT or delegation for required NFT"
			} else {
				fmt.Println("No valid NFT or delegation found for the address.")
				response.Status = "error"
				response.Message = "Address does not own or have delegation for required NFT"
			}
		} else {
			// No incoming delegation recorded: skip updating clients collection.
			fmt.Println("No incoming delegations recorded.")
			response.Status = "error"
			response.Message = "Address does not have any incoming delegations"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
