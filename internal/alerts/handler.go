package alerts

import (
	"net/http"

	"github.com/PabloJTechDev/loyalty-core-backoffice/internal/shared"
)

// Handler handles operational alert endpoints.
type Handler struct {
	items []OperationalAlert
}

// NewHandler creates a Handler with the given alert seed data.
func NewHandler(items []OperationalAlert) *Handler {
	return &Handler{items: items}
}

// List responds with all registered operational alerts.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	shared.WriteJSON(w, http.StatusOK, ListResponse{
		Total: len(h.items),
		Items: h.items,
	})
}

// DefaultAlerts is the seed data used when no external source is available.
var DefaultAlerts = []OperationalAlert{
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
