package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aliasthewho/price_tracker/internal/storage/pantry"
)

func main() {
	// Load configuration from environment variables
	cfg, err := pantry.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a new basket manager
	manager := pantry.NewBasketManager(cfg)

	// Create a basket name based on today's date
	basketName := pantry.BasketName(time.Now())

	// Check if basket exists
	exists, err := manager.BasketExists(context.Background(), basketName)
	if err != nil {
		log.Fatalf("Failed to check if basket exists: %v", err)
	}

	if exists {
		fmt.Printf("Basket %s already exists\n", basketName)
	} else {
		// Create a new basket
		if err := manager.CreateBasket(context.Background(), basketName); err != nil {
			log.Fatalf("Failed to create basket: %v", err)
		}
		fmt.Printf("Created new basket: %s\n", basketName)
	}

	// Example data to store
	data := map[string]interface{}{
		"products": []map[string]interface{}{
			{
				"name":  "Example Product",
				"price": 9.99,
			},
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Update the basket with data
	if err := manager.UpdateBasket(context.Background(), basketName, data); err != nil {
		log.Fatalf("Failed to update basket: %v", err)
	}

	fmt.Printf("Successfully updated basket %s with data\n", basketName)

	// List all baskets
	baskets, err := manager.ListBaskets(context.Background())
	if err != nil {
		log.Fatalf("Failed to list baskets: %v", err)
	}

	fmt.Println("\nAll baskets in your pantry:")
	for _, b := range baskets {
		fmt.Printf("- %s\n", b)
	}
}
