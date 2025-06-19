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

// BasketManager provides methods to interact with Pantry baskets.
// It handles all communication with the Pantry API, including
// creating, updating, and querying baskets.
//
// The zero value is not usable, use NewBasketManager instead.
type BasketManager struct {
	// baseURL is the base URL of the Pantry API
	baseURL string
	// apiKey is the API key for authentication
	apiKey string
	// httpClient is the HTTP client for making requests
	httpClient *http.Client
}

// Config holds the configuration required to initialize a Pantry client.
type Config struct {
	// APIKey is the authentication token for the Pantry API.
	// It can be obtained from the Pantry dashboard at https://getpantry.cloud/
	APIKey string
}

// NewConfigFromEnv creates a new Config by reading the PANTRY_API_KEY environment variable.
//
// Returns an error if the environment variable is not set or is empty.
//
// Example:
//
//	cfg, err := NewConfigFromEnv()
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewConfigFromEnv() (Config, error) {
	apiKey := os.Getenv("PANTRY_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("PANTRY_API_KEY environment variable not set")
	}
	return Config{APIKey: apiKey}, nil
}

// NewBasketManager creates a new BasketManager with the provided configuration.
//
// The returned BasketManager is ready to interact with the Pantry API.
// The default HTTP client has a 10-second timeout.
//
// Example:
//
//	cfg := Config{APIKey: "your-api-key"}
//	manager := NewBasketManager(cfg)
func NewBasketManager(cfg Config) *BasketManager {
	return &BasketManager{
		baseURL:    "https://getpantry.cloud/apiv1/pantry",
		apiKey:     cfg.APIKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// BasketName generates a consistent name for a basket based on the provided date.
// The format is "prices_YYYY_MM_DD".
//
// Example:
//
//	name := BasketName(time.Date(2023, 6, 18, 0, 0, 0, 0, time.UTC))
//	// name is "prices_2023_06_18"
func BasketName(date time.Time) string {
	return fmt.Sprintf("prices_%s", date.Format("2006_01_02"))
}

// CreateBasket creates a new basket in Pantry with the given name.
//
// The basket name must be unique within your Pantry. If a basket with the same name
// already exists, this function will return an error.
//
// Example:
//
//	err := manager.CreateBasket(ctx, "my-basket")
//	if err != nil {
//	    return fmt.Errorf("failed to create basket: %w", err)
//	}
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

// BasketExists checks if a basket with the given name exists in Pantry.
//
// Returns true if the basket exists, false otherwise. Any error during the check
// (e.g., network issues) is returned as the second parameter.
//
// Example:
//
//	exists, err := manager.BasketExists(ctx, "my-basket")
//	if err != nil {
//	    return false, fmt.Errorf("failed to check basket existence: %w", err)
//	}
func (m *BasketManager) BasketExists(ctx context.Context, basketName string) (bool, error) {
	url := fmt.Sprintf("%s/%s/basket/%s", m.baseURL, m.apiKey, basketName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
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

// ListBaskets retrieves the names of all baskets in your Pantry.
//
// Returns a slice of basket names. The list includes all baskets regardless
// of their content or creation time.
//
// Example:
//
//	baskets, err := manager.ListBaskets(ctx)
//	if err != nil {
//	    return nil, fmt.Errorf("failed to list baskets: %w", err)
//	}
func (m *BasketManager) ListBaskets(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/%s/baskets", m.baseURL, m.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
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

// UpdateBasket updates the contents of an existing basket with the provided data.
//
// The data parameter can be any value that can be marshaled to JSON.
// If the basket doesn't exist, it will be created.
//
// Example:
//
//	data := map[string]interface{}{"key": "value"}
//	if err := manager.UpdateBasket(ctx, "my-basket", data); err != nil {
//	    return fmt.Errorf("failed to update basket: %w", err)
//	}
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

// GetBasket retrieves and unmarshals the contents of a basket into the target.
//
// The target parameter must be a pointer to a value that can hold the unmarshaled
// JSON data. Returns an error if the basket doesn't exist or if the data cannot
// be unmarshaled into the target type.
//
// Example:
//
//	var data map[string]interface{}
//	if err := manager.GetBasket(ctx, "my-basket", &data); err != nil {
//	    return fmt.Errorf("failed to get basket: %w", err)
//	}
func (m *BasketManager) GetBasket(ctx context.Context, basketName string, target interface{}) error {
	url := fmt.Sprintf("%s/%s/basket/%s", m.baseURL, m.apiKey, basketName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
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
