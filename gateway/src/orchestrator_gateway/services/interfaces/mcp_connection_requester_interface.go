package interfaces

import "github.com/google/uuid"

type MCPConnection struct {
	ID      string
	AgentID string
	Name    string
	Url     string
}

type CreateMCPConnectionRequest struct {
	AgentId    uuid.UUID
	Name       string
	Url        string
	HeaderAuth map[string]string
}

type ListMCPConnectionsRequest struct {
	AgentId         *uuid.UUID
	Search          *string
	SortBy          *string
	Page            *uint
	SortDirection   *string
	ElementsPerPage *uint
}

type ListMCPConnectionsResponse struct {
	Page            uint
	ElementsPerPage uint
	TotalElements   uint
	MCPConnections  []*MCPConnection
}

type MCPConnectionRequesterInterface interface {
	Remove(id uuid.UUID) error
	Create(request *CreateMCPConnectionRequest) error
	List(request *ListMCPConnectionsRequest) (*ListMCPConnectionsResponse, error)
}
