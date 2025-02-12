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

		// Fetch delegations made by the specified address
		delegations, err := db.GetFromDelegationsByAddress(address)
		if err != nil {
			http.Error(w, "Failed to fetch delegations", http.StatusInternalServerError)
			return
		}

		// Prepare the response
		var clientsWithDelegations []ClientWithDelegations
		clientDelegationMap := make(map[string]int64) // To accumulate delegated amounts

		for _, delegation := range delegations {
			// Fetch the client who received the delegation
			clientAddress := delegation.ToAddress
			client, err := db.GetClient(clientAddress)
			if err != nil {
				http.Error(w, "Failed to fetch client", http.StatusInternalServerError)
				return
			}

			// Only include clients that are running the node
			if client != nil {
				// Accumulate the delegated amount for each client
				clientDelegationMap[clientAddress] += delegation.Amount
			}
		}

		// Now update the NFTAmount for each client based on accumulated amounts
		for clientAddress, totalDelegatedAmount := range clientDelegationMap {
			client, err := db.GetClient(clientAddress)
			if err != nil {
				http.Error(w, "Failed to fetch client", http.StatusInternalServerError)
				return
			}

			if client != nil {
				// Set the NFTAmount to the total delegated amount
				client.NFTAmount = totalDelegatedAmount // Set NFTAmount directly
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
