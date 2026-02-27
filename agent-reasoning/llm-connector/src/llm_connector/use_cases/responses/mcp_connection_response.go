package responses

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/google/uuid"
)

type MCPConnectionResponse struct {
	ID      uuid.UUID `json:"id"`
	AgentID uuid.UUID `json:"agentId"`
	Name    string    `json:"name"`
	Url     string    `json:"url"`
}

func (response *MCPConnectionResponse) FromModel(model *models.MCPConnection) {
	response.ID = model.ID
	response.AgentID = model.AgentID
	response.Name = model.Name
	response.Url = model.Url
}
