package models

type ExternalOpenAICompletionsStandardResponse struct {
	ID          string                                            `json:"id"`
	Object      string                                            `json:"object"`
	Model       string                                            `json:"model"`
	Choices     []ExternalOpenAICompletionsStandardResponseChoice `json:"choices"`
	ServiceTier string                                            `json:"service_tier"`
}

type ExternalOpenAICompletionsStandardResponseChoice struct {
	Index        int                                                    `json:"index"`
	Message      ExternalOpenAICompletionsStandardResponseChoiceMessage `json:"message"`
	FinishReason string                                                 `json:"finish_reason"`
}

type ExternalOpenAICompletionsStandardResponseChoiceMessage struct {
	Role        string                                                           `json:"role"`
	Content     string                                                           `json:"content"`
	Refusal     interface{}                                                      `json:"refusal"`
	Annotations []interface{}                                                    `json:"annotations"`
	ToolCalls   []ExternalOpenAICompletionsStandardResponseChoiceMessageToolCall `json:"tool_calls,omitempty"`
}

type ExternalOpenAICompletionsStandardResponseChoiceMessageToolCall struct {
	ID       string                                            `json:"id"`
	Type     string                                            `json:"type"` // typically "function"
	Function ExternalOpenAICompletionsStandardToolCallFunction `json:"function"`
}
