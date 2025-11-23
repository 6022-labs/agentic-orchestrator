package responses

type ListMCPConnectionsResponse struct {
	Page            uint                     `json:"page"`
	TotalElements   uint64                   `json:"totalElements"`
	ElementsPerPage uint                     `json:"elementsPerPage"`
	MCPConnections  []*MCPConnectionResponse `json:"mcpConnections"`
}
