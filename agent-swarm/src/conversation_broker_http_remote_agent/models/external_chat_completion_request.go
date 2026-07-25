package models

// Minimal OpenAI chat-completions request for nodes that only advertise the
// chatCompletions dialect.
type ExternalChatCompletionRequest struct {
	Messages []*ExternalChatCompletionMessage `json:"messages"`
	Metadata map[string]string                `json:"metadata,omitempty"`
}

type ExternalChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
