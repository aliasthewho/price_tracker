package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// Parse command line flags
	outputFile := flag.String("output", "", "Output JSON file (default: stdout)")
	dateStr := flag.String("date", "", "Date in YYYY-MM-DD format (default: today)")
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
