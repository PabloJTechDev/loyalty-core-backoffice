package capabilities

// Capability represents a backoffice operational capability.
type Capability struct {
	CapabilityID string   `json:"capabilityId"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	DependsOn    []string `json:"dependsOn"`
	Status       string   `json:"status"`
}

// ListResponse wraps a slice of Capability with a total count.
type ListResponse struct {
	Total int          `json:"total"`
	Items []Capability `json:"items"`
}
