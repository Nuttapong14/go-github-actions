package handlers

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// DeploymentsTotal tracks the total number of deployments by environment and status
	DeploymentsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "deployments_total",
			Help: "Total number of deployments",
		},
		[]string{"environment", "status"},
	)

	// HealthCheckDuration tracks health check response times
	HealthCheckDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "health_check_duration_seconds",
			Help:    "Health check response time in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"check_type"},
	)

	// HTTPRequestsInFlight tracks currently processing HTTP requests
	HTTPRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// DatabaseConnectionsActive tracks active database connections
	DatabaseConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)
)

func init() {
	// Register custom metrics
	prometheus.MustRegister(
		DeploymentsTotal,
		HealthCheckDuration,
		HTTPRequestsInFlight,
		DatabaseConnectionsActive,
	)
}

// SetupMetrics configures Prometheus metrics collection for the Fiber app
// Returns the prometheus middleware for use in the app
func SetupMetrics(app *fiber.App, serviceName string) *fiberprometheus.FiberPrometheus {
	// Create fiberprometheus instance
	prometheus := fiberprometheus.New(serviceName)

	// Register the /metrics endpoint
	prometheus.RegisterAt(app, "/metrics")

	// Return the middleware for the app to use
	return prometheus
}

// RecordDeployment records a deployment metric
func RecordDeployment(environment, status string) {
	DeploymentsTotal.WithLabelValues(environment, status).Inc()
}

// RecordHealthCheckDuration records the duration of a health check
func RecordHealthCheckDuration(checkType string, durationSeconds float64) {
	HealthCheckDuration.WithLabelValues(checkType).Observe(durationSeconds)
}

// SetHTTPRequestsInFlight sets the number of HTTP requests currently in flight
func SetHTTPRequestsInFlight(count float64) {
	HTTPRequestsInFlight.Set(count)
}

// SetDatabaseConnectionsActive sets the number of active database connections
func SetDatabaseConnectionsActive(count float64) {
	DatabaseConnectionsActive.Set(count)
}
