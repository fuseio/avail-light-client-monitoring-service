package handlers

import (
	"avail-light-client-monitoring-service/database"
	"encoding/json"
	"net/http"
)

type GetClientsResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    []database.ClientInfo `json:"data"`
}

func GetClients(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
