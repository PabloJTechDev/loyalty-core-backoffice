package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalizeRoute(t *testing.T) {
	if route := normalizeRoute("/v1/backoffice-capabilities/support-dashboard"); route != "/v1/backoffice-capabilities/:capabilityId" {
		t.Fatalf("unexpected route: %s", route)
	}

	if route := normalizeRoute("/v1/operational-alerts/alert_points_sync"); route != "/v1/operational-alerts/:alertId" {
		t.Fatalf("unexpected route: %s", route)
	}
}

func TestWriteJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeJSON(recorder, http.StatusCreated, map[string]string{"status": "accepted"})

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if payload["status"] != "accepted" {
		t.Fatalf("expected accepted, got %q", payload["status"])
	}
}

func TestGetPortFallback(t *testing.T) {
	t.Setenv("PORT", "invalid")
	if port := getPort(); port != 3005 {
		t.Fatalf("expected fallback port 3005, got %d", port)
	}
}

func TestCapabilitiesSeeded(t *testing.T) {
	if len(capabilities) == 0 {
		t.Fatal("expected seeded capabilities")
	}
	if capabilities[0].CapabilityID == "" {
		t.Fatal("expected capability id")
	}
}
