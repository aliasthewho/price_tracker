package pantry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type (
	// ErrorResponse represents an error response from the Pantry API
	ErrorResponse struct {
		Message string `json:"message"`
	}

	// Basket represents a Pantry basket
	Basket map[string]interface{}
)

// BasketManager handles operations on Pantry baskets
type BasketManager struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Config holds the configuration for the Pantry client
type Config struct {
	APIKey string
}

// NewConfigFromEnv creates a new Config with values from environment variables
func NewConfigFromEnv() (Config, error) {
	apiKey := os.Getenv("PANTRY_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("PANTRY_API_KEY environment variable not set")
	}
	return Config{APIKey: apiKey}, nil
}

// NewBasketManager creates a new Pantry basket manager
func NewBasketManager(cfg Config) *BasketManager {
	return &BasketManager{
		baseURL:    "https://getpantry.cloud/apiv1/pantry",
		apiKey:     cfg.APIKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// BasketName generates a consistent name for a basket based on date
func BasketName(date time.Time) string {
	return fmt.Sprintf("prices_%s", date.Format("2006_01_02"))
}

// CreateBasket creates a new basket in Pantry
func (m *BasketManager) CreateBasket(ctx context.Context, basketName string) error {
	url := fmt.Sprintf("%s/%s/basket/%s", m.baseURL, m.apiKey, basketName)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil {
			return fmt.Errorf("failed to create basket: %s", errResp.Message)
		}
		return fmt.Errorf("failed to create basket: %s", resp.Status)
	}
	
	return nil
}

// BasketExists checks if a basket exists in Pantry
func (m *BasketManager) BasketExists(ctx context.Context, basketName string) (bool, error) {
	url := fmt.Sprintf("%s/%s/basket/%s", m.baseURL, m.apiKey, basketName)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK, nil
}

// ListBaskets lists all baskets in the pantry
func (m *BasketManager) ListBaskets(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/%s/baskets", m.baseURL, m.apiKey)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil {
			return nil, fmt.Errorf("failed to list baskets: %s", errResp.Message)
		}
		return nil, fmt.Errorf("failed to list baskets: %s", resp.Status)
	}
	
	var baskets []string
	if err := json.NewDecoder(resp.Body).Decode(&baskets); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return baskets, nil
}

// UpdateBasket updates the contents of a basket
func (m *BasketManager) UpdateBasket(ctx context.Context, basketName string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := fmt.Sprintf("%s/%s/basket/%s", m.baseURL, m.apiKey, basketName)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil {
			return fmt.Errorf("failed to update basket: %s", errResp.Message)
		}
		return fmt.Errorf("failed to update basket: %s", resp.Status)
	}
	
	return nil
}

// GetBasket retrieves the contents of a basket
func (m *BasketManager) GetBasket(ctx context.Context, basketName string, target interface{}) error {
	url := fmt.Sprintf("%s/%s/basket/%s", m.baseURL, m.apiKey, basketName)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil {
			return fmt.Errorf("failed to get basket: %s", errResp.Message)
		}
		return fmt.Errorf("failed to get basket: %s", resp.Status)
	}
	
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	return nil
}
