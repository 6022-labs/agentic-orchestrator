package models

import "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"

type ExternalMCPConnectionResponse struct {
	ID      string `json:"id"`
	AgentID string `json:"agentId"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

func (response ExternalMCPConnectionResponse) ToMCPConnection() *interfaces.MCPConnection {
	return &interfaces.MCPConnection{
		ID:      response.ID,
		Url:     response.Url,
		Name:    response.Name,
		AgentID: response.AgentID,
	}
}
