package handler

import (
	"analytics-go-api-vercel/repo"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type TrackEventRequest struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

type TrackEventResponse struct {
	Message string `json:"message"`
}

func TrackEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, PATCH, DELETE, POST, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req TrackEventRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	connString := os.Getenv("MONGODB_URI")
	database := os.Getenv("MONGODB_DBNAME")

	repo, err := repo.NewMongoRepository(connString, database, "Test")
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	repo.Track(req.Event, req.Properties)

	response := TrackEventResponse{
		Message: "Event tracked successfully",
	}
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
