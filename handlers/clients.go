package handlers

import (
	"avail-light-client-monitoring-service/database"
	"encoding/json"
	"net/http"
	"strings"
)

type GetClientsResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    []database.ClientInfo `json:"data"`
}

type GetClientResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    database.ClientInfo `json:"data"`
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

func GetClient(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get a client from database by address
		address := strings.TrimPrefix(r.URL.Path, "/client/")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		client, err := db.GetClient(address)
		if err != nil {
			http.Error(w, "Failed to fetch client", http.StatusInternalServerError)
			return
		}

		response := GetClientResponse{
			Status:  "success",
			Message: "Client retrieved successfully",
			Data:    client,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}