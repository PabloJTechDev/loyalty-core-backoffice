package alerts

// OperationalAlert represents a backoffice operational alert.
type OperationalAlert struct {
	AlertID     string `json:"alertId"`
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Status      string `json:"status"`
}

// ListResponse wraps a slice of OperationalAlert with a total count.
type ListResponse struct {
	Total int                `json:"total"`
	Items []OperationalAlert `json:"items"`
}
