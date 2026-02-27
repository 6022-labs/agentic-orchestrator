package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/clients"
	database_connector_models "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/models"
	"github.com/google/uuid"
)

type AgentContextRequester struct {
	client clients.AgentDatabaseConnectorClientInterface
}

func NewAgentContextRequester(
	client clients.AgentDatabaseConnectorClientInterface,
) *AgentContextRequester {
	return &AgentContextRequester{
		client: client,
	}
}

func (a *AgentContextRequester) CreateAgentContext(agentId uuid.UUID, key string) error {
	return a.client.CreateAgentContext(&database_connector_models.ExternalCreateAgentContextRequest{
		AgentId: agentId,
		Key:     key,
	})
}

func (a *AgentContextRequester) GetAgentContextByKey(agentId uuid.UUID, key string) (*models.AgentContext, *bool, error) {
	agentContext, err := a.client.GetAgentContextByKey(&database_connector_models.ExternalGetAgentContextByKeyRequest{
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
	agentContexts, err := a.client.GetAgentContexts(&database_connector_models.ExternalAgentIdRequest{
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
