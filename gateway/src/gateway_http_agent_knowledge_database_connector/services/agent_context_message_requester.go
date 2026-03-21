package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_knowledge_database_connector/clients"
	database_connector_models "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_knowledge_database_connector/models"
	"github.com/google/uuid"
)

type AgentContextMessageRequester struct {
	client clients.AgentDatabaseConnectorClientInterface
}

func NewAgentContextMessageRequester(
	client clients.AgentDatabaseConnectorClientInterface,
) *AgentContextMessageRequester {
	return &AgentContextMessageRequester{
		client: client,
	}
}

func (a *AgentContextMessageRequester) AddMessagesToAgentContext(request *interfaces.AddMessageToAgentContextRequest) error {
	externalRequest := &database_connector_models.ExternalAddMessagesToAgentContextRequest{
		AgentId:        request.AgentId,
		Messages:       []database_connector_models.ExternalMessage{},
		ContextKey:     request.ContextKey,
		ConversationId: request.ConversationId,
	}

	for _, message := range request.Messages {
		externalRequest.Messages = append(externalRequest.Messages, database_connector_models.ExternalMessage{
			Role:           message.Role,
			Author:         message.Author,
			Content:        message.Content,
			Platform:       message.Platform,
			ConversationId: message.ConversationId,
		})
	}

	return a.client.AddMessagesToAgentContext(externalRequest)
}

func (a *AgentContextMessageRequester) GetAgentContextMessages(agentId uuid.UUID, key string) ([]*models.AgentContextMessage, error) {
	messages, err := a.client.GetAgentContextMessages(&database_connector_models.ExternalGetAgentContextMessagesRequest{
		AgentId: agentId,
		Key:     key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context messages: %w", err)
	}

	return messages.ToAgentContextMessages(), nil
}
