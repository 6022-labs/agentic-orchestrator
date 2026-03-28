package services

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/clients"
	agent_knowledge_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/models"
	"github.com/google/uuid"
)

type AgentContextConversationMappingRequester struct {
	client clients.AgentKnowledgeClientInterface
}

func NewAgentContextConversationMappingRequester(
	client clients.AgentKnowledgeClientInterface,
) *AgentContextConversationMappingRequester {
	return &AgentContextConversationMappingRequester{
		client: client,
	}
}

func (a *AgentContextConversationMappingRequester) GetLinkedAgentContextKeyOfConversation(agentId uuid.UUID, conversationId string) (*models.AgentContext, *bool, error) {
	response, err := a.client.GetLinkedAgentContextOfConversation(&agent_knowledge_models.ExternalGetLinkedAgentContextOfConversationRequest{
		AgentId:        agentId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, nil, err
	}

	isFound := response.Status == agent_knowledge_models.GetLinkedAgentContextOfConversationResponseStatusFound

	return &models.AgentContext{
		Id:  response.Id,
		Key: response.Key,
	}, &isFound, nil
}

func (a *AgentContextConversationMappingRequester) LinkConversationToAgentContext(agentId uuid.UUID, conversationId string, key string) error {
	return a.client.LinkConversationToAgentContext(&agent_knowledge_models.ExternalLinkConversationToAgentContextRequest{
		AgentId:        agentId,
		ConversationId: conversationId,
		Key:            key,
	})
}
