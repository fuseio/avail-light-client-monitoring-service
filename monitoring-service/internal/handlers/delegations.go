package handlers

import (
	"encoding/json"
	"net/http"
	"monitoring-service/internal/database"
)

type GetDelegationsResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Clients []ClientWithDelegations `json:"clients"`
}

type ClientWithDelegations struct {
	ClientInfo *database.ClientInfo `json:"client"`
}

func GetDelegations(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		delegations, err := db.GetFromDelegationsByAddress(address)
		if err != nil {
			http.Error(w, "Failed to fetch delegations", http.StatusInternalServerError)
			return
		}

		var clientsWithDelegations []ClientWithDelegations
		clientDelegationMap := make(map[string]int64) 

		for _, delegation := range delegations {
			clientAddress := delegation.ToAddress
			client, err := db.GetClient(clientAddress)
			if err != nil {
				http.Error(w, "Failed to fetch client", http.StatusInternalServerError)
				return
			}

			if client != nil {
				clientDelegationMap[clientAddress] += delegation.Amount
			}
		}

		for clientAddress, totalDelegatedAmount := range clientDelegationMap {
			client, err := db.GetClient(clientAddress)
			if err != nil {
				http.Error(w, "Failed to fetch client", http.StatusInternalServerError)
				return
			}

			if client != nil {
				client.NFTAmount = totalDelegatedAmount
				clientsWithDelegations = append(clientsWithDelegations, ClientWithDelegations{
					ClientInfo: client,
				})
			}
		}

		response := GetDelegationsResponse{
			Status:  "success",
			Message: "Delegations retrieved successfully",
			Clients: clientsWithDelegations,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
