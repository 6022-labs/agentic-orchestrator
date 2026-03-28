package services

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/models"
	"github.com/google/uuid"
)

type AgentPromptTemplateRequester struct {
	client clients.AgentKnowledgeClientInterface
}

func NewAgentPromptTemplateRequester(
	client clients.AgentKnowledgeClientInterface,
) *AgentPromptTemplateRequester {
	return &AgentPromptTemplateRequester{
		client: client,
	}
}

func (a *AgentPromptTemplateRequester) GetAgentPromptTemplateByAgentIdAndKey(agentId uuid.UUID, key enums.AgentPromptTemplateKey) (*string, error) {
	promptTemplate, err := a.client.GetAgentPromptTemplate(&models.ExternalGetAgentPromptTemplateRequest{
		AgentId:     agentId,
		TemplateKey: key,
	})

	if err != nil {
		return nil, err
	}

	return &promptTemplate.PromptTemplate, nil
}
