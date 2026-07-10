package models

// ExternalCreateResponseRequest is the OpenAI-Responses-compatible request a
// 6022 node accepts on its responses endpoint. Only the fields the broker
// sends are modeled.
type ExternalCreateResponseRequest struct {
	Instructions string                          `json:"instructions,omitempty"`
	Input        []*ExternalResponseInputMessage `json:"input"`
	Metadata     map[string]string               `json:"metadata,omitempty"`
}

type ExternalResponseInputMessage struct {
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content string `json:"content"`
}
