package models

type ExternalAnthropicMessagesStandardResponse struct {
	Model string `json:"model"`
	ID    string `json:"id"`
	ExternalAnthropicMessagesStandardMessage
}
