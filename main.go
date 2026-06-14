package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PabloJTechDev/loyalty-core-backoffice/internal/alerts"
	"github.com/PabloJTechDev/loyalty-core-backoffice/internal/capabilities"
	"github.com/PabloJTechDev/loyalty-core-backoffice/internal/health"
	"github.com/PabloJTechDev/loyalty-core-backoffice/internal/shared"
)

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

func routes() http.Handler {
	capsHandler := capabilities.NewHandler(capabilities.DefaultCapabilities)
	alertsHandler := alerts.NewHandler(alerts.DefaultAlerts)
	healthHandler := health.NewHandler(len(capabilities.DefaultCapabilities))

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(shared.MetricsRegistry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/health", shared.WithMetrics("/health", healthHandler.Handle))
	mux.HandleFunc("/v1/backoffice-capabilities", shared.WithMetrics("/v1/backoffice-capabilities", capsHandler.List))
	mux.HandleFunc("/v1/operational-alerts", shared.WithMetrics("/v1/operational-alerts", alertsHandler.List))

	return mux
}

func main() {
	port := getPort()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shared.LogEvent("service.started", shared.LogFields{
		"port":    port,
		"storage": "mock-memory",
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
