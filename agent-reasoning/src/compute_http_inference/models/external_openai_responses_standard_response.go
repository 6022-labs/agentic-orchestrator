package models

import "encoding/json"

type ExternalOpenAIResponsesStandardResponse struct {
	ID                 string                                           `json:"id"`
	Object             string                                           `json:"object"`
	CreatedAt          int                                              `json:"created_at"`
	Status             string                                           `json:"status"`
	Error              interface{}                                      `json:"error"`
	IncompleteDetails  interface{}                                      `json:"incomplete_details"`
	Instructions       interface{}                                      `json:"instructions"`
	MaxOutputTokens    interface{}                                      `json:"max_output_tokens"`
	Model              string                                           `json:"model"`
	Output             []ExternalOpenAIResponsesStandardResponseOutput  `json:"output"`
	ParallelToolCalls  bool                                             `json:"parallel_tool_calls"`
	PreviousResponseID interface{}                                      `json:"previous_response_id"`
	Reasoning          ExternalOpenAIResponsesStandardResponseReasoning `json:"reasoning"`
	Store              bool                                             `json:"store"`
	Temperature        float64                                          `json:"temperature"`
	Text               ExternalOpenAIResponsesStandardResponseText      `json:"text"`
	ToolChoice         string                                           `json:"tool_choice"`
	Tools              []interface{}                                    `json:"tools"`
	TopP               float64                                          `json:"top_p"`
	Truncation         string                                           `json:"truncation"`
	Usage              ExternalOpenAIResponsesStandardResponseUsage     `json:"usage"`
	User               interface{}                                      `json:"user"`
	Metadata           ExternalOpenAIResponsesStandardResponseMetadata  `json:"metadata"`
}

// Output item type for the response
type ExternalOpenAIResponsesStandardResponseOutput struct {
	Type      string                                                 `json:"type"`
	ID        string                                                 `json:"id"`
	Status    string                                                 `json:"status"`
	Role      *string                                                `json:"role,omitempty"`
	Name      *string                                                `json:"name,omitempty"`
	CallId    *string                                                `json:"call_id,omitempty"`
	Arguments *string                                                `json:"arguments,omitempty"`
	Content   []ExternalOpenAIResponsesStandardResponseOutputContent `json:"content,omitempty"`
}

func (o *ExternalOpenAIResponsesStandardResponseOutput) ParsedArguments() (map[string]any, error) {
	if o.Arguments == nil {
		return make(map[string]any), nil
	}

	var args map[string]any
	err := json.Unmarshal([]byte(*o.Arguments), &args)
	if err != nil {
		return nil, err
	}
	return args, nil
}

// Content item contained within an output entry
type ExternalOpenAIResponsesStandardResponseOutputContent struct {
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	Annotations []interface{} `json:"annotations"`
}

// Reasoning details for the response
type ExternalOpenAIResponsesStandardResponseReasoning struct {
	Effort  interface{} `json:"effort"`
	Summary interface{} `json:"summary"`
}

// Text formatting section
type ExternalOpenAIResponsesStandardResponseText struct {
	Format ExternalOpenAIResponsesStandardResponseTextFormat `json:"format"`
}

type ExternalOpenAIResponsesStandardResponseTextFormat struct {
	Type string `json:"type"`
}

// Usage metrics for the response
type ExternalOpenAIResponsesStandardResponseUsage struct {
	InputTokens         int                                                             `json:"input_tokens"`
	InputTokensDetails  ExternalOpenAIResponsesStandardResponseUsageInputTokensDetails  `json:"input_tokens_details"`
	OutputTokens        int                                                             `json:"output_tokens"`
	OutputTokensDetails ExternalOpenAIResponsesStandardResponseUsageOutputTokensDetails `json:"output_tokens_details"`
	TotalTokens         int                                                             `json:"total_tokens"`
}

type ExternalOpenAIResponsesStandardResponseUsageInputTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type ExternalOpenAIResponsesStandardResponseUsageOutputTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

// Metadata associated with the response
type ExternalOpenAIResponsesStandardResponseMetadata struct{}
