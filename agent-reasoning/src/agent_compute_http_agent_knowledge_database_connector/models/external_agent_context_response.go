package models

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/google/uuid"
)

type ExternalAgentContextResponse struct {
	Id  uuid.UUID `json:"id"`
	Key string    `json:"key"`
}

type ExternalAgentContextResponses []*ExternalAgentContextResponse

func (e *ExternalAgentContextResponse) ToAgentContext() *models.AgentContext {
	return &models.AgentContext{
		Id:  e.Id,
		Key: e.Key,
	}
}

func (e *ExternalAgentContextResponses) ToAgentContexts() []*models.AgentContext {
	agentContexts := make([]*models.AgentContext, len(*e))

	for i, response := range *e {
		agentContexts[i] = response.ToAgentContext()
	}

	return agentContexts
}
