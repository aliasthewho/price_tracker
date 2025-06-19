package pantry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeJSON is a helper function to write JSON responses with proper error handling
func writeJSON(t *testing.T, w http.ResponseWriter, v interface{}) error {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func TestBasketNaming(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "specific date",
			date:     time.Date(2025, 6, 17, 0, 0, 0, 0, time.UTC),
			expected: "prices_2025_06_17",
		},
		{
			name:     "another date",
			date:     time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			expected: "prices_2023_12_01",
		},
		{
			name:     "another_date",
			date:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: "prices_2023_12_31",
		},
		{
			name:     "leap_year",
			date:     time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: "prices_2024_02_29",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BasketName(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewConfigFromEnv(t *testing.T) {
	t.Parallel()
	// Save and restore environment variables
	oldKey := os.Getenv("PANTRY_API_KEY")
	defer os.Setenv("PANTRY_API_KEY", oldKey)

	tests := []struct {
		name    string
		setup   func()
		err     string
		wantKey string
	}{
		{
			name: "valid key",
			setup: func() {
				os.Setenv("PANTRY_API_KEY", "test-key-123")
			},
			wantKey: "test-key-123",
		},
		{
			name: "missing key",
			setup: func() {
				os.Unsetenv("PANTRY_API_KEY")
			},
			err: "PANTRY_API_KEY environment variable not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			os.Clearenv()
			if tt.setup != nil {
				tt.setup()
			}

			// Test
			got, err := NewConfigFromEnv()

			// Verify
			if tt.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantKey, got.APIKey)
		})
	}
}

func TestBasketManager(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		apiKey      string
		setupServer func(t *testing.T) *httptest.Server
		testFunc    func(t *testing.T, manager *BasketManager, server *httptest.Server)
	}{
		{
			name:   "NewBasketManager",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				assert.NotNil(t, manager)
				assert.Equal(t, "https://getpantry.cloud/apiv1/pantry", manager.baseURL)
				assert.Equal(t, "test-key", manager.apiKey)
				assert.NotNil(t, manager.httpClient)
			},
		},
		{
			name:   "CreateBasket success",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/apiv1/pantry/test-key/basket/test-basket", r.URL.Path)
					assert.Equal(t, http.MethodPost, r.Method)
					w.WriteHeader(http.StatusOK)
					if err := writeJSON(t, w, map[string]string{"message": "Basket created"}); err != nil {
						t.Errorf("Failed to write JSON response: %v", err)
					}
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				err := manager.CreateBasket(context.Background(), "test-basket")
				assert.NoError(t, err)
			},
		},
		{
			name:   "BasketExists true",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/apiv1/pantry/test-key/basket/existing-basket", r.URL.Path)
					w.WriteHeader(http.StatusOK)
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				exists, err := manager.BasketExists(context.Background(), "existing-basket")
				assert.NoError(t, err)
				assert.True(t, exists)
			},
		},
		{
			name:   "BasketExists false",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				exists, err := manager.BasketExists(context.Background(), "non-existent-basket")
				assert.NoError(t, err)
				assert.False(t, exists)
			},
		},
		{
			name:   "ListBaskets success",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/apiv1/pantry/test-key/baskets", r.URL.Path)
					w.WriteHeader(http.StatusOK)
					if err := writeJSON(t, w, []string{"basket1", "basket2"}); err != nil {
						t.Errorf("Failed to write JSON response: %v", err)
					}
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				baskets, err := manager.ListBaskets(context.Background())
				require.NoError(t, err)
				assert.Equal(t, []string{"basket1", "basket2"}, baskets)
			},
		},
		{
			name:   "UpdateBasket success",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/apiv1/pantry/test-key/basket/test-basket", r.URL.Path)
					assert.Equal(t, http.MethodPut, r.Method)

					var data map[string]interface{}
					err := json.NewDecoder(r.Body).Decode(&data)
					require.NoError(t, err)
					assert.Equal(t, "test", data["key"])

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					if err := writeJSON(t, w, map[string]string{"key": "value"}); err != nil {
						t.Errorf("Failed to write JSON response: %v", err)
					}
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				err := manager.UpdateBasket(context.Background(), "test-basket", map[string]string{"key": "test"})
				assert.NoError(t, err)
			},
		},
		{
			name:   "GetBasket success",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/apiv1/pantry/test-key/basket/test-basket", r.URL.Path)
					w.WriteHeader(http.StatusOK)
					if err := writeJSON(t, w, map[string]string{"key": "value"}); err != nil {
						t.Errorf("Failed to write JSON response: %v", err)
					}
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				var result map[string]string
				err := manager.GetBasket(context.Background(), "test-basket", &result)
				require.NoError(t, err)
				assert.Equal(t, "value", result["key"])
			},
		},
		{
			name:   "CreateBasket error response",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					if err := writeJSON(t, w, ErrorResponse{Message: "invalid request"}); err != nil {
						t.Errorf("Failed to write JSON response: %v", err)
					}
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				err := manager.CreateBasket(context.Background(), "test-basket")
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid request")
			},
		},
		{
			name:   "UpdateBasket error response",
			apiKey: "test-key",
			setupServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					if err := writeJSON(t, w, ErrorResponse{Message: "invalid data"}); err != nil {
						t.Errorf("Failed to write JSON response: %v", err)
					}
				}))
			},
			testFunc: func(t *testing.T, manager *BasketManager, server *httptest.Server) {
				manager.baseURL = server.URL + "/apiv1/pantry"
				err := manager.UpdateBasket(context.Background(), "test-basket", map[string]string{"key": "value"})
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid data")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer(t)
			defer server.Close()

			manager := NewBasketManager(Config{APIKey: tt.apiKey})
			tt.testFunc(t, manager, server)
		})
	}
}
