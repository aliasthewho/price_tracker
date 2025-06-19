package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// PriceRequestsTotal counts the total number of price requests
	PriceRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "price_requests_total",
		Help: "Total number of price requests",
	}, []string{"status"}) // "success" or "error"

	// PriceRequestDuration tracks the duration of price requests
	PriceRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "price_request_duration_seconds",
		Help:    "Duration of price requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"endpoint"})

	// PantryOperationsTotal counts the total number of Pantry operations
	PantryOperationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pantry_operations_total",
		Help: "Total number of Pantry operations",
	}, []string{"operation", "status"}) // operation: "get", "set", "delete"; status: "success", "error"

	// PantryOperationDuration tracks the duration of Pantry operations
	PantryOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "pantry_operation_duration_seconds",
		Help:    "Duration of Pantry operations in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
)

// RecordPriceRequest records metrics for a price request
func RecordPriceRequest(status string, duration float64, endpoint string) {
	PriceRequestsTotal.WithLabelValues(status).Inc()
	PriceRequestDuration.WithLabelValues(endpoint).Observe(duration)
}

// RecordPantryOperation records metrics for a Pantry operation
func RecordPantryOperation(operation, status string, duration float64) {
	PantryOperationsTotal.WithLabelValues(operation, status).Inc()
	PantryOperationDuration.WithLabelValues(operation).Observe(duration)
}
