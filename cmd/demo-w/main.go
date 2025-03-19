package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	user, err := s.client.CreateUser(r.Context(), req.Name, req.Email)
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

	userID, err := fromStatelyUUID(req.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID format %s", err.Error()), http.StatusBadRequest)
		return
	}

	resourceID, err := fromStatelyUUID(req.ResourceID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid resource ID format %s", err.Error()), http.StatusBadRequest)
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

	userID, err := fromStatelyUUID(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user ID format %s (%s)", err.Error(), r.URL.Path[len("/users/"):]), http.StatusBadRequest)
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

	resourceID, err := fromStatelyUUID(r.URL.Path[len("/resources/"):])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid resource ID format %s (%s)", err.Error(), r.URL.Path[len("/resources/"):]), http.StatusBadRequest)
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

func fromStatelyUUID(b64id string) (uuid.UUID, error) {
	// decode the b64id to a byte slice
	id, err := base64.StdEncoding.DecodeString(b64id)
	if err != nil {
		return uuid.Nil, err
	}
	// convert the byte slice to a UUID
	u, err := uuid.FromBytes(id)
	if err != nil {
		return uuid.Nil, err
	}
	return u, nil
}
