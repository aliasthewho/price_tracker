package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	scraper "github.com/aliasthewho/price_tracker/internal/api/emmsa"
	"github.com/aliasthewho/price_tracker/internal/storage/pantry"
)

func main() {
	// Parse command line flags
	outputFile := flag.String("output", "", "Output JSON file (default: stdout)")
	dateStr := flag.String("date", "", "Date in YYYY-MM-DD format (default: today)")
	enablePantry := flag.Bool("pantry", false, "Enable Pantry storage")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Set up logging
	log.SetFlags(0)
	if !*debug {
		log.SetOutput(io.Discard)
	}

	// Parse date
	var targetDate time.Time
	if *dateStr == "" {
		targetDate = time.Now()
	} else {
		var err error
		targetDate, err = time.Parse("2006-01-02", *dateStr)
		if err != nil {
			log.Fatalf("Invalid date format. Please use YYYY-MM-DD: %v", err)
		}
	}

	// Create a new EMMSA scraper
	s, err := scraper.NewEMMSAScraper()
	if err != nil {
		log.Fatalf("Failed to create scraper: %v", err)
	}
	defer s.Close()

	// Scrape prices
	prices, err := s.ScrapePrices(targetDate)
	if err != nil {
		log.Fatalf("Failed to fetch prices: %v", err)
	}

	// Prepare data for storage
	data := map[string]interface{}{
		"date":    targetDate.Format("2006-01-02"),
		"prices":  prices,
		"fetched": time.Now().Format(time.RFC3339),
	}

	// Store in Pantry if enabled
	if *enablePantry {
		saveToPantry(targetDate, data)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(prices, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal prices to JSON: %v", err)
	}

	// Output results
	if *outputFile != "" {
		// Write to file
		err = os.WriteFile(*outputFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		fmt.Printf("Successfully saved prices to %s\n", *outputFile)
	} else {
		// Print to stdout
		fmt.Println(string(jsonData))
	}
}

func saveToPantry(date time.Time, data interface{}) {
	// Initialize Pantry client
	cfg, err := pantry.NewConfigFromEnv()
	if err != nil {
		log.Printf("Warning: Failed to initialize Pantry client: %v", err)
		return
	}

	manager := pantry.NewBasketManager(cfg)
	basketName := pantry.BasketName(date)

	// Check if basket exists
	exists, err := manager.BasketExists(context.Background(), basketName)
	if err != nil {
		log.Printf("Warning: Failed to check if basket exists: %v", err)
		return
	}

	// Create basket if it doesn't exist
	if !exists {
		if err := manager.CreateBasket(context.Background(), basketName); err != nil {
			log.Printf("Warning: Failed to create basket: %v", err)
			return
		}
		log.Printf("Created new Pantry basket: %s", basketName)
	}

	// Update basket with data
	if err := manager.UpdateBasket(context.Background(), basketName, data); err != nil {
		log.Printf("Warning: Failed to update basket: %v", err)
		return
	}

	log.Printf("Successfully updated Pantry basket: %s", basketName)
}
