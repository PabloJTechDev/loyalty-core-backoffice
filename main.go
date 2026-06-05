package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type logFields map[string]any

func logEvent(event string, fields logFields) {
	payload := map[string]any{
		"ts":      time.Now().UTC().Format(time.RFC3339),
		"service": "loyalty-core-backoffice",
		"event":   event,
	}

	for key, value := range fields {
		if value != nil {
			payload[key] = value
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("{\"ts\":\"%s\",\"service\":\"loyalty-core-backoffice\",\"event\":\"logger.error\",\"message\":\"json_marshal_failed\"}\n", time.Now().UTC().Format(time.RFC3339))
		return
	}

	fmt.Println(string(bytes))
}

var (
	metricsRegistry = prometheus.NewRegistry()
	coreHTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_core_backoffice_http_requests_total",
			Help: "Total HTTP requests handled by the backoffice core service",
		},
		[]string{"method", "route", "status_class", "status_code"},
	)
	coreHTTPRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loyalty_core_backoffice_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds for the backoffice core service",
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.2, 0.35, 0.5, 0.75, 1, 1.5, 2, 3, 5},
		},
		[]string{"method", "route", "status_class", "status_code"},
	)
)

func init() {
	metricsRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		coreHTTPRequestsTotal,
		coreHTTPRequestDurationSeconds,
	)
}

type healthResponse struct {
	Status       string `json:"status"`
	Service      string `json:"service"`
	Storage      string `json:"storage"`
	Capabilities int    `json:"capabilities"`
}

type capability struct {
	CapabilityID string   `json:"capabilityId"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	DependsOn    []string `json:"dependsOn"`
	Status       string   `json:"status"`
}

type capabilityListResponse struct {
	Total int          `json:"total"`
	Items []capability `json:"items"`
}

type operationalAlert struct {
	AlertID     string `json:"alertId"`
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Status      string `json:"status"`
}

type alertListResponse struct {
	Total int                `json:"total"`
	Items []operationalAlert `json:"items"`
}

type writer struct {
	http.ResponseWriter
	status int
}

func (w *writer) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func getPort() int {
	portValue := strings.TrimSpace(os.Getenv("PORT"))
	if portValue == "" {
		return 3005
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || port <= 0 {
		log.Printf("invalid PORT %q, falling back to 3005", portValue)
		return 3005
	}

	return port
}

func normalizeRoute(path string) string {
	switch {
	case strings.HasPrefix(path, "/v1/backoffice-capabilities/"):
		return "/v1/backoffice-capabilities/:capabilityId"
	case strings.HasPrefix(path, "/v1/operational-alerts/"):
		return "/v1/operational-alerts/:alertId"
	default:
		return path
	}
}

func statusClass(statusCode int) string {
	return fmt.Sprintf("%dxx", statusCode/100)
}

func withMetrics(route string, handler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		wrapped := &writer{ResponseWriter: response, status: http.StatusOK}
		handler(wrapped, request)

		labels := prometheus.Labels{
			"method":       request.Method,
			"route":        normalizeRoute(route),
			"status_class": statusClass(wrapped.status),
			"status_code":  strconv.Itoa(wrapped.status),
		}
		coreHTTPRequestsTotal.With(labels).Inc()
		coreHTTPRequestDurationSeconds.With(labels).Observe(time.Since(startedAt).Seconds())
	}
}

func writeJSON(response http.ResponseWriter, status int, payload any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(payload)
}

var capabilities = []capability{
	{
		CapabilityID: "support-dashboard",
		Title:        "Support dashboard",
		Description:  "Operational summary for support, ops and audit workflows.",
		DependsOn:    []string{"loyalty-core-points"},
		Status:       "planned",
	},
	{
		CapabilityID: "customer-review",
		Title:        "Customer review",
		Description:  "Technical capability to aggregate customer snapshots for backoffice operations.",
		DependsOn:    []string{"loyalty-core-points"},
		Status:       "planned",
	},
	{
		CapabilityID: "order-monitoring",
		Title:        "Order monitoring",
		Description:  "Operational follow-up for ecommerce orders and flagged cases.",
		DependsOn:    []string{"loyalty-core-backoffice", "loyalty-core-points"},
		Status:       "planned",
	},
}

var alerts = []operationalAlert{
	{
		AlertID:     "alert_points_sync",
		Severity:    "medium",
		Title:       "Points sync pending",
		Description: "Backoffice still relies on mock data while loyalty-core-points integration is pending.",
		Domain:      "points",
		Status:      "open",
	},
	{
		AlertID:     "alert_order_review",
		Severity:    "low",
		Title:       "Order review backlog",
		Description: "Mock backlog kept visible until ecommerce operational review is connected to real sources.",
		Domain:      "ecommerce",
		Status:      "open",
	},
}

func main() {
	port := getPort()
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/health", withMetrics("/health", func(response http.ResponseWriter, request *http.Request) {
		writeJSON(response, http.StatusOK, healthResponse{
			Status:       "ok",
			Service:      "loyalty-core-backoffice",
			Storage:      "mock-memory",
			Capabilities: len(capabilities),
		})
	}))

	mux.HandleFunc("/v1/backoffice-capabilities", withMetrics("/v1/backoffice-capabilities", func(response http.ResponseWriter, request *http.Request) {
		writeJSON(response, http.StatusOK, capabilityListResponse{Total: len(capabilities), Items: capabilities})
	}))

	mux.HandleFunc("/v1/operational-alerts", withMetrics("/v1/operational-alerts", func(response http.ResponseWriter, request *http.Request) {
		writeJSON(response, http.StatusOK, alertListResponse{Total: len(alerts), Items: alerts})
	}))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logEvent("service.started", logFields{
		"port":    port,
		"storage": "mock-memory",
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
