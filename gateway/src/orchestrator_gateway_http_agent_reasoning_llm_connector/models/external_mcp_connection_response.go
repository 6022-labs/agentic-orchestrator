package models

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"

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
