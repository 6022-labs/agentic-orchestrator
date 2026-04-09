package models

import "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"

type ExternalListMCPConnectionsResponse struct {
	MCPConnections  []ExternalMCPConnectionResponse `json:"mcpConnections"`
	Page            uint                            `json:"page"`
	ElementsPerPage uint                            `json:"elementsPerPage"`
	TotalElements   uint                            `json:"totalElements"`
}

func (response ExternalListMCPConnectionsResponse) ToInternal() *interfaces.ListMCPConnectionsResponse {
	mcpConnections := make([]*interfaces.MCPConnection, len(response.MCPConnections))
	for i, conn := range response.MCPConnections {
		mcpConnections[i] = conn.ToMCPConnection()
	}

	return &interfaces.ListMCPConnectionsResponse{
		Page:            response.Page,
		ElementsPerPage: response.ElementsPerPage,
		TotalElements:   response.TotalElements,
		MCPConnections:  mcpConnections,
	}
}
