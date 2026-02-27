package responses

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/google/uuid"
)

type GetAgentContextResponse struct {
	Id  uuid.UUID `json:"id"`
	Key string    `json:"key"`
}

func GetAgentContextResponseFromAgentContext(agentContext *models.AgentContext) *GetAgentContextResponse {
	return &GetAgentContextResponse{
		Id:  agentContext.ID,
		Key: agentContext.Key,
	}
}

func GetAgentContextResponsesFromAgentContexts(agentContexts []*models.AgentContext) []*GetAgentContextResponse {
	response := make([]*GetAgentContextResponse, len(agentContexts))
	for i, agentContext := range agentContexts {
		response[i] = GetAgentContextResponseFromAgentContext(agentContext)
	}
	return response
}
