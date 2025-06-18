package scraper

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestEMMSAScraper(t *testing.T) {
	// Create a new scraper
	s, err := NewEMMSAScraper()
	if err != nil {
		t.Fatalf("Failed to create scraper: %v", err)
	}
	defer s.Close()

	// Test with a recent date
	t.Run("ScrapeRecentPrices", func(t *testing.T) {
		t.Log("Testing price scraping for a recent date")

		// Use a date from the last 7 days
		date := time.Now().AddDate(0, 0, -1)
		prices, err := s.ScrapePrices(date)

		// Check for errors
		if err != nil {
			t.Fatalf("Failed to scrape prices: %v", err)
		}

		// Basic validation of the results
		if len(prices) == 0 {
			t.Error("No prices returned from the API")
			return
		}

		t.Logf("Successfully retrieved %d price entries\n", len(prices))

		// Convert to pretty-printed JSON
		jsonData, err := json.MarshalIndent(prices, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal prices to JSON: %v", err)
		}

		// Write to a file
		filename := "emmsa_prices.json"
		err = os.WriteFile(filename, jsonData, 0644)
		if err != nil {
			t.Fatalf("Failed to write JSON to file: %v", err)
		}

		t.Logf("First 5 entries (full data saved to %s):\n", filename)

		// Print first 5 entries to console
		for i, p := range prices {
			if i >= 5 {
				break
			}
			entry, _ := json.MarshalIndent(p, "", "  ")
			t.Logf("Entry %d:\n%s\n", i+1, string(entry))
		}

		t.Logf("\nFull data has been saved to %s", filename)
	})
}
