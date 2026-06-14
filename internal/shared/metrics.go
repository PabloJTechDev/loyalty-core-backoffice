package shared

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	// MetricsRegistry is the application-wide Prometheus registry.
	MetricsRegistry = prometheus.NewRegistry()

	// CoreHTTPRequestsTotal counts HTTP requests by method, route, status class and code.
	CoreHTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_core_backoffice_http_requests_total",
			Help: "Total HTTP requests handled by the backoffice core service",
		},
		[]string{"method", "route", "status_class", "status_code"},
	)

	// CoreHTTPRequestDurationSeconds tracks request latency.
	CoreHTTPRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loyalty_core_backoffice_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds for the backoffice core service",
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.2, 0.35, 0.5, 0.75, 1, 1.5, 2, 3, 5},
		},
		[]string{"method", "route", "status_class", "status_code"},
	)
)

func init() {
	MetricsRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		CoreHTTPRequestsTotal,
		CoreHTTPRequestDurationSeconds,
	)
}
