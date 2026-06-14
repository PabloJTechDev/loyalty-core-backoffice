package health

// Response is the payload returned by the health endpoint.
type Response struct {
	Status       string `json:"status"`
	Service      string `json:"service"`
	Storage      string `json:"storage"`
	Capabilities int    `json:"capabilities"`
}
