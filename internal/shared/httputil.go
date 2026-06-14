package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ResponseWriter wrapper that captures status code.
type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

// WriteJSON writes a JSON-encoded payload with the given HTTP status.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// NormalizeRoute collapses dynamic path segments for metric labels.
func NormalizeRoute(path string) string {
	switch {
	case strings.HasPrefix(path, "/v1/backoffice-capabilities/"):
		return "/v1/backoffice-capabilities/:capabilityId"
	case strings.HasPrefix(path, "/v1/operational-alerts/"):
		return "/v1/operational-alerts/:alertId"
	default:
		return path
	}
}

// StatusClass returns e.g. "2xx" for a 200 status.
func StatusClass(statusCode int) string {
	return fmt.Sprintf("%dxx", statusCode/100)
}

// WithMetrics wraps a handler to record Prometheus request metrics.
func WithMetrics(route string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		wrapped := &ResponseWriter{ResponseWriter: w, Status: http.StatusOK}
		handler(wrapped, r)

		labels := prometheus.Labels{
			"method":       r.Method,
			"route":        NormalizeRoute(route),
			"status_class": StatusClass(wrapped.Status),
			"status_code":  strconv.Itoa(wrapped.Status),
		}
		CoreHTTPRequestsTotal.With(labels).Inc()
		CoreHTTPRequestDurationSeconds.With(labels).Observe(time.Since(startedAt).Seconds())
	}
}
