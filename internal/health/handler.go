package health

import (
	"net/http"

	"github.com/PabloJTechDev/loyalty-core-backoffice/internal/shared"
)

// Handler handles the /health endpoint.
type Handler struct {
	capabilitiesCount int
}

// NewHandler creates a HealthHandler that reports the given capability count.
func NewHandler(capabilitiesCount int) *Handler {
	return &Handler{capabilitiesCount: capabilitiesCount}
}

// Handle responds with service health status.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	shared.WriteJSON(w, http.StatusOK, Response{
		Status:       "ok",
		Service:      "loyalty-core-backoffice",
		Storage:      "mock-memory",
		Capabilities: h.capabilitiesCount,
	})
}
