package models

import "time"

type ExternalAgentComputeInputRequest struct {
	Author    *string   `json:"author"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ExternalAgentComputeContextConversationAgentRequest struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionTokenId uint64 `json:"collectionTokenId"`
}

type ExternalAgentComputeContextConversationContextRequest struct {
	Id     string                                                 `json:"id"`
	Agents []*ExternalAgentComputeContextConversationAgentRequest `json:"agents"`
}

type ExternalAgentComputeRequest struct {
	Platform            string                                                 `json:"platform"`
	ConversationContext *ExternalAgentComputeContextConversationContextRequest `json:"conversationContext"`
	ForceConclusion     bool                                                   `json:"forceConclusion"`
	Inputs              []*ExternalAgentComputeInputRequest                    `json:"inputs"`
	PromptParameters    map[string]any                                         `json:"promptParameters"`
}
