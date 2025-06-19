package scraper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	emmsaAPIURL = "https://old.emmsa.com.pe/emmsa_spv/app/reportes/ajax/rpt07_gettable_new_web.php"
)

// EMMSAPrice represents the price data from EMMSA
type EMMSAPrice struct {
	Date       string  `json:"date"`
	Product    string  `json:"product"`
	Variedad   string  `json:"variedad"`
	PrecioMin  float64 `json:"precio_min"`
	PrecioMax  float64 `json:"precio_max"`
	PrecioProm float64 `json:"precio_prom"`
}

// EMMSAScraper handles fetching price data from the EMMSA API
type EMMSAScraper struct {
	httpClient *http.Client
}

// NewEMMSAScraper creates a new EMMSA scraper
func NewEMMSAScraper() (*EMMSAScraper, error) {
	return &EMMSAScraper{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// parsePriceTable parses the HTML table from the API response
func parsePriceTable(html []byte, date time.Time) ([]EMMSAPrice, error) {
	log.Printf("Parsing price table for date: %s", date.Format("2006-01-02"))
	log.Printf("Response length: %d bytes", len(html))

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var prices []EMMSAPrice

	// Find all table rows
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		// Skip header row
		if i == 0 {
			return
		}

		// Extract data from each cell
		cells := s.Find("td")
		if cells.Length() < 5 { // Ensure we have enough cells
			return
		}

		// Extract text from each cell
		product := strings.TrimSpace(cells.Eq(0).Text())
		variedad := strings.TrimSpace(cells.Eq(1).Text())
		precioMinStr := strings.TrimSpace(cells.Eq(2).Text())
		precioMaxStr := strings.TrimSpace(cells.Eq(3).Text())
		precioPromStr := strings.TrimSpace(cells.Eq(4).Text())

		// Convert price strings to float64 (handle potential errors)
		precioMin, err1 := strconv.ParseFloat(precioMinStr, 64)
		precioMax, err2 := strconv.ParseFloat(precioMaxStr, 64)
		precioProm, err3 := strconv.ParseFloat(precioPromStr, 64)

		// Skip rows with invalid price data
		if err1 != nil || err2 != nil || err3 != nil {
			log.Printf("Skipping row with invalid price data: %s, %s, %s",
				precioMinStr, precioMaxStr, precioPromStr)
			return
		}

		price := EMMSAPrice{
			Date:       date.Format("2006-01-02"),
			Product:    product,
			Variedad:   variedad,
			PrecioMin:  precioMin,
			PrecioMax:  precioMax,
			PrecioProm: precioProm,
		}

		prices = append(prices, price)
	})

	return prices, nil
}

// ScrapePrices fetches the daily prices from EMMSA API
func (s *EMMSAScraper) ScrapePrices(date time.Time) ([]EMMSAPrice, error) {
	// Format the date as dd/mm/yyyy for the API
	formattedDate := date.Format("02/01/2006")
	log.Printf("Fetching prices for date: %s", formattedDate)

	// Prepare form data
	formData := url.Values{
		"vid_tipo": {"1"}, // 1 = Precios Diarios
		"vprod":    {""},  // Empty for all products
		"vvari":    {""},  // Empty for all varieties
		"vfecha":   {formattedDate},
	}

	// Create a new request
	req, err := http.NewRequest("POST", emmsaAPIURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers to match the curl request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "text/html, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Origin", "https://old.emmsa.com.pe")
	req.Header.Set("Referer", "https://old.emmsa.com.pe/emmsa_spv/rpEstadistica/rpt_precios-diarios-web.php")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.5 Safari/605.1.15")

	// Send the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the HTML response
	return parsePriceTable(body, date)
}

// Close releases any resources used by the scraper
func (s *EMMSAScraper) Close() error {
	// No resources to close with HTTP client
	return nil
}
