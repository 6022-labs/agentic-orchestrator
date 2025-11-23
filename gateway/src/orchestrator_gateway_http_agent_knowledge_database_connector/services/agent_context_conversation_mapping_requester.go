package services

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/clients"
	database_connector_models "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/models"
	"github.com/google/uuid"
)

type AgentContextConversationMappingRequester struct {
	client clients.AgentDatabaseConnectorClientInterface
}

func NewAgentContextConversationMappingRequester(
	client clients.AgentDatabaseConnectorClientInterface,
) *AgentContextConversationMappingRequester {
	return &AgentContextConversationMappingRequester{
		client: client,
	}
}

func (a *AgentContextConversationMappingRequester) GetLinkedAgentContextKeyOfConversation(agentId uuid.UUID, conversationId string) (*models.AgentContext, *bool, error) {
	response, err := a.client.GetLinkedAgentContextOfConversation(&database_connector_models.ExternalGetLinkedAgentContextOfConversationRequest{
		AgentId:        agentId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, nil, err
	}

	isFound := response.Status == database_connector_models.GetLinkedAgentContextOfConversationResponseStatusFound

	return &models.AgentContext{
		Id:  response.Id,
		Key: response.Key,
	}, &isFound, nil
}

func (a *AgentContextConversationMappingRequester) LinkConversationToAgentContext(agentId uuid.UUID, conversationId string, key string) error {
	return a.client.LinkConversationToAgentContext(&database_connector_models.ExternalLinkConversationToAgentContextRequest{
		AgentId:        agentId,
		ConversationId: conversationId,
		Key:            key,
	})
}
