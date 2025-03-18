package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/StatelyCloud/demo-w/pkg/client"
	"github.com/google/uuid"
)

type server struct {
	client *client.Client
}

type createUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type createResourceRequest struct {
	Name string `json:"name"`
}

type createLeaseRequest struct {
	UserID      string  `json:"userId"`
	ResourceID  string  `json:"resourceId"`
	Reason      string  `json:"reason"`
	DurationHrs float64 `json:"durationHours"`
}

const PORT = "8080"

func main() {
	ctx := context.Background()

	storeStr := os.Getenv("STATELY_STORE_ID")
	if storeStr == "" {
		log.Fatal("STATELY_STORE_ID environment variable is required")
	}
	storeID, err := strconv.ParseUint(storeStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid store ID: %v", err)
	}

	c, err := client.NewClient(ctx, storeID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	s := &server{client: c}

	// Register routes
	http.HandleFunc("/users", s.handleCreateUser)
	http.HandleFunc("/resources", s.handleCreateResource)
	http.HandleFunc("/leases", s.handleCreateLease)
	http.HandleFunc("/users/", s.handleGetUserLeases)
	http.HandleFunc("/resources/", s.handleGetResourceLeases)

	log.Printf("Server starting on port %s", PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.client.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *server) handleCreateResource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resource, err := s.client.CreateResource(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

func (s *server) handleCreateLease(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createLeaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	resourceID, err := uuid.Parse(req.ResourceID)
	if err != nil {
		http.Error(w, "Invalid resource ID format", http.StatusBadRequest)
		return
	}

	lease, err := s.client.CreateLease(r.Context(), userID, resourceID, req.Reason,
		time.Duration(req.DurationHrs*float64(time.Hour)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lease)
}

func (s *server) handleGetUserLeases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := uuid.Parse(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	leases, err := s.client.GetLeasesForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leases)
}

func (s *server) handleGetResourceLeases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resourceID, err := uuid.Parse(r.URL.Path[len("/resources/"):])
	if err != nil {
		http.Error(w, "Invalid resource ID format", http.StatusBadRequest)
		return
	}

	leases, err := s.client.GetLeasesForResource(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leases)
}
