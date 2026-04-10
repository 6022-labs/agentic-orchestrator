package models

type ExternalOllamaChatMessageStandardResponse struct {
	Model   string                                   `json:"model"`
	Message ExternalOllamaChatMessageStandardMessage `json:"message"`
}
