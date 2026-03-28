package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/clients"
	agent_knowledge_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/models"
	"github.com/google/uuid"
)

type AgentContextRequester struct {
	client clients.AgentKnowledgeClientInterface
}

func NewAgentContextRequester(
	client clients.AgentKnowledgeClientInterface,
) *AgentContextRequester {
	return &AgentContextRequester{
		client: client,
	}
}

func (a *AgentContextRequester) CreateAgentContext(agentId uuid.UUID, key string) error {
	return a.client.CreateAgentContext(&agent_knowledge_models.ExternalCreateAgentContextRequest{
		AgentId: agentId,
		Key:     key,
	})
}

func (a *AgentContextRequester) GetAgentContextByKey(agentId uuid.UUID, key string) (*models.AgentContext, *bool, error) {
	agentContext, err := a.client.GetAgentContextByKey(&agent_knowledge_models.ExternalGetAgentContextByKeyRequest{
		AgentId: agentId,
		Key:     key,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent context by key: %w", err)
	}

	isFound := agentContext != nil
	if !isFound {
		return nil, &isFound, nil
	}

	return agentContext.ToAgentContext(), &isFound, nil
}

func (a *AgentContextRequester) GetAgentContexts(agentId uuid.UUID, search, sortBy, sortDirection *string) ([]*models.AgentContext, error) {
	agentContexts, err := a.client.GetAgentContexts(&agent_knowledge_models.ExternalAgentIdRequest{
		AgentId:       agentId,
		Search:        search,
		SortBy:        sortBy,
		SortDirection: sortDirection,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context keys: %w", err)
	}

	return agentContexts.ToAgentContexts(), nil
}
