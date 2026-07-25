package models

type ExternalChatCompletionResponse struct {
	Choices []*ExternalChatCompletionChoice `json:"choices"`
}

type ExternalChatCompletionChoice struct {
	Message *ExternalChatCompletionMessage `json:"message"`
}
