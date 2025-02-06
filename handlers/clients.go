package handlers

import (
	"avail-light-client-monitoring-service/database"
	"encoding/json"
	"net/http"
)

type GetClientsResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    []database.ClientInfo `json:"data"`
}

type ClientWithHistoryResponse struct {
	Client      *database.ClientInfo        `json:"client"`
	Heartbeats  []database.HeartbeatRecord `json:"heartbeats"`
	Delegations []database.DelegationRecord `json:"delegations"`
}

func GetClients(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if client address is provided
		address := r.URL.Query().Get("address")
		if address != "" {
			client, heartbeats, delegations, err := db.GetClientWithHistory(address)
			if err != nil {
				http.Error(w, "Failed to fetch client", http.StatusInternalServerError)
				return
			}
			
			response := ClientWithHistoryResponse{
				Client:      client,
				Heartbeats:  heartbeats,
				Delegations: delegations,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Get all clients without history
		clients, err := db.GetAllClients()
		if err != nil {
			http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
			return
		}

		response := GetClientsResponse{
			Status:  "success",
			Message: "Clients retrieved successfully",
			Data:    clients,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
