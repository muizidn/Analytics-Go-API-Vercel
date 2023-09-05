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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Use POST.", http.StatusMethodNotAllowed)
		return
	}

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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
