package capabilities

import (
	"net/http"

	"github.com/PabloJTechDev/loyalty-core-backoffice/internal/shared"
)

// Handler handles capability endpoints.
type Handler struct {
	items []Capability
}

// NewHandler creates a Handler with the given capability seed data.
func NewHandler(items []Capability) *Handler {
	return &Handler{items: items}
}

// List responds with all registered capabilities.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	shared.WriteJSON(w, http.StatusOK, ListResponse{
		Total: len(h.items),
		Items: h.items,
	})
}

// DefaultCapabilities is the seed data used when no external source is available.
var DefaultCapabilities = []Capability{
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
