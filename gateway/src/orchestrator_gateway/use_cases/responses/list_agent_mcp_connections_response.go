package responses

type ListAgentMCPConnectionsResponse struct {
	Page            uint                     `json:"page"`
	ElementsPerPage uint                     `json:"elementsPerPage"`
	TotalElements   uint                     `json:"totalElements"`
	MCPConnections  []*MCPConnectionResponse `json:"mcpConnections"`
}
