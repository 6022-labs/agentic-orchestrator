package services

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/models"
	"github.com/google/uuid"
)

type AgentPromptTemplateRequester struct {
	client clients.AgentDatabaseConnectorClientInterface
}

func NewAgentPromptTemplateRequester(
	client clients.AgentDatabaseConnectorClientInterface,
) *AgentPromptTemplateRequester {
	return &AgentPromptTemplateRequester{
		client: client,
	}
}

func (a *AgentPromptTemplateRequester) SaveAgentPromptTemplate(agentId uuid.UUID, key enums.AgentPromptTemplateKey, content string) error {
	err := a.client.SaveAgentPromptTemplate(&models.ExternalSaveAgentPromptTemplateRequest{
		AgentId: agentId,
		Content: content,
		Key:     key,
	})

	if err != nil {
		return err
	}

	return nil
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
