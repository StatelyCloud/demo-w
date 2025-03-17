package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/StatelyCloud/demo-w/pkg/client"
)

func main() {
	ctx := context.Background()
	storeStr := flag.String("store", "", "Stately store ID")
	flag.Parse()
	if *storeStr == "" {
		log.Fatal("store ID is required")
	}
	storeID, err := strconv.ParseUint(*storeStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid store ID: %v", err)
	}

	c, err := client.NewClient(ctx, storeID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a user
	user, err := c.CreateUser(ctx, "john.doe@company.com", "John Doe")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user: %v\n", user)

	// Create a resource
	resource, err := c.CreateResource(ctx, "sensitive-service")
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}
	fmt.Printf("Created resource: %v\n", resource)

	// Create a lease for John Doe to access the sensitive service for 3 hours
	lease, err := c.CreateLease(ctx, user.Id, resource.Id, "Debugging an error", 3*time.Hour)
	if err != nil {
		log.Fatalf("Failed to create lease: %v", err)
	}
	fmt.Printf("Created lease: %v\n", lease)

	// List all leases for John Doe
	leases, err := c.GetLeasesForUser(ctx, user.Id)
	if err != nil {
		log.Fatalf("Failed to list leases: %v", err)
	}
	fmt.Println("Current leases for John Doe:")
	for _, l := range leases {
		fmt.Printf("- %v\n", l)
	}

	// List all leases for sensitive-service
	leases, err = c.GetLeasesForResource(ctx, resource.Id)
	if err != nil {
		log.Fatalf("Failed to list leases: %v", err)
	}
	fmt.Println("Current leases for sensitive-service:")
	for _, l := range leases {
		fmt.Printf("- %v\n", l)
	}
}
