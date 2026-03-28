package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/clients"
	agent_knowledge_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/models"
	"github.com/google/uuid"
)

type AgentContextMessageRequester struct {
	client clients.AgentKnowledgeClientInterface
}

func NewAgentContextMessageRequester(
	client clients.AgentKnowledgeClientInterface,
) *AgentContextMessageRequester {
	return &AgentContextMessageRequester{
		client: client,
	}
}

func (a *AgentContextMessageRequester) AddMessagesToAgentContext(request *interfaces.AddMessageToAgentContextRequest) error {
	externalRequest := &agent_knowledge_models.ExternalAddMessagesToAgentContextRequest{
		AgentId:        request.AgentId,
		Messages:       []agent_knowledge_models.ExternalMessage{},
		ContextKey:     request.ContextKey,
		ConversationId: request.ConversationId,
	}

	for _, message := range request.Messages {
		externalRequest.Messages = append(externalRequest.Messages, agent_knowledge_models.ExternalMessage{
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
	messages, err := a.client.GetAgentContextMessages(&agent_knowledge_models.ExternalGetAgentContextMessagesRequest{
		AgentId: agentId,
		Key:     key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context messages: %w", err)
	}

	return messages.ToAgentContextMessages(), nil
}
