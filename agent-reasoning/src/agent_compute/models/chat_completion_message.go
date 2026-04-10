package models

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
