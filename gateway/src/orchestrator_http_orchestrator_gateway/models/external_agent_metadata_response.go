package models

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
)

type ExternalAgentMetadataResponse struct {
	LLMModel    string                      `json:"llmModel"`
	LLMPlatform enums.LLMConnectionPlatform `json:"llmPlatform"`
}

func (m *ExternalAgentMetadataResponse) ToInternalAgentMetadataInformation() *interfaces.AgentMetadataInformation {
	return &interfaces.AgentMetadataInformation{
		LLMModel:    m.LLMModel,
		LLMPlatform: m.LLMPlatform,
	}
}
