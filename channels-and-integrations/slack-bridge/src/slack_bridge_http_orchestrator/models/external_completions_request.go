package models

import (
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

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

func (r *ExternalAgentComputeContextConversationContextRequest) SetAgents(agents []models.SlackAgent) *ExternalAgentComputeContextConversationContextRequest {
	r.Agents = make([]*ExternalAgentComputeContextConversationAgentRequest, len(agents))
	for i, agent := range agents {
		r.Agents[i] = &ExternalAgentComputeContextConversationAgentRequest{
			CollectionAddress: agent.CollectionAddress,
			CollectionTokenId: agent.CollectionTokenId,
		}
	}

	return r
}

type ExternalAgentComputeRequest struct {
	Platform            string                                                 `json:"platform"`
	ConversationContext *ExternalAgentComputeContextConversationContextRequest `json:"conversationContext"`
	ForceConclusion     bool                                                   `json:"forceConclusion"`
	Inputs              []*ExternalAgentComputeInputRequest                    `json:"inputs"`
	PromptParameters    map[string]any                                         `json:"promptParameters"`
}
