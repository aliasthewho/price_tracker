package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aliasthewho/price_tracker/internal/metrics"
	scraper "github.com/aliasthewho/price_tracker/internal/api/emmsa"
	"github.com/aliasthewho/price_tracker/internal/storage/pantry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Parse command line flags
	outputFile := flag.String("output", "", "Output JSON file (default: stdout)")
	dateStr := flag.String("date", "", "Date in YYYY-MM-DD format (default: today)")
	enablePantry := flag.Bool("pantry", false, "Enable Pantry storage")
	debug := flag.Bool("debug", false, "Enable debug logging")
	metricsAddr := flag.String("metrics-addr", ":2112", "The address to expose Prometheus metrics")
	flag.Parse()

	// Start metrics server in a goroutine
	metricsServer := &http.Server{
		Addr:    *metricsAddr,
		Handler: promhttp.Handler(),
	}

	go func() {
		log.Printf("Starting metrics server on %s", *metricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Set up logging
	if *debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// Parse date
	var date time.Time
	var err error
	if *dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", *dateStr)
		if err != nil {
			log.Fatalf("Invalid date format: %v. Expected YYYY-MM-DD", err)
		}
	}

	// Run the price scraping and keep the metrics server running in the background
	runPriceScraping(date, *enablePantry, *outputFile)
	
	// Wait for interrupt signal to gracefully shutdown the server
	log.Println("Press Ctrl+C to exit")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
	
	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Shutdown the metrics server
	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down metrics server: %v", err)
	}
}

func runPriceScraping(date time.Time, enablePantry bool, outputFile string) {
	// Create a new EMMSA scraper
	s, err := scraper.NewEMMSAScraper()
	if err != nil {
		log.Fatalf("Failed to create scraper: %v", err)
	}

	// Scrape prices with metrics
	startTime := time.Now()
	prices, err := s.ScrapePrices(date)
	duration := time.Since(startTime).Seconds()
	
	// Record metrics
	status := "success"
	if err != nil {
		status = "error"
	}
	metrics.RecordPriceRequest(status, duration, "scrape")
	if err != nil {
		s.Close()
		log.Fatalf("Failed to fetch prices: %v", err)
	}
	// Don't use defer with Fatalf as it won't run deferred functions
	s.Close()

	// Prepare data for storage
	data := map[string]interface{}{
		"date":    date.Format("2006-01-02"),
		"prices":  prices,
		"fetched": time.Now().Format(time.RFC3339),
	}

	// Save to Pantry if enabled
	if enablePantry {
		startTime := time.Now()
		err = saveToPantry(date, data)
		duration := time.Since(startTime).Seconds()
		
		// Record metrics
		status := "success"
		if err != nil {
			status = "error"
		}
		metrics.RecordPantryOperation("save", status, duration)
		
		if err != nil {
			log.Fatalf("Failed to save to Pantry: %v", err)
		}
	}

	// Marshal prices to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal prices to JSON: %v", err)
	}

	// Output results
	if outputFile != "" {
		// Write to file
		err = os.WriteFile(outputFile, jsonData, 0o600) // Use 0o600 for octal literal
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		log.Printf("Prices written to %s", outputFile)
	} else {
		// Print to stdout
		fmt.Println(string(jsonData))
	}
}

func saveToPantry(date time.Time, data interface{}) error {
	// Initialize Pantry config
	cfg, err := pantry.NewConfigFromEnv()
	if err != nil {
		return fmt.Errorf("error loading Pantry config: %w", err)
	}

	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize BasketManager
	manager := pantry.NewBasketManager(cfg)
	basketName := pantry.BasketName(date)

	// Check if basket exists
	exists, err := manager.BasketExists(ctx, basketName)
	if err != nil {
		return fmt.Errorf("error checking if basket exists: %w", err)
	}

	// Create basket if it doesn't exist
	if !exists {
		if err := manager.CreateBasket(ctx, basketName); err != nil {
			return fmt.Errorf("error creating basket: %w", err)
		}
		log.Printf("Created new Pantry basket: %s", basketName)
	}

	// Update basket with data
	if err := manager.UpdateBasket(ctx, basketName, data); err != nil {
		return fmt.Errorf("error updating basket: %w", err)
	}

	log.Printf("Successfully updated Pantry basket: %s", basketName)
	return nil
}
