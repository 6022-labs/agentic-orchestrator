package models

import "github.com/google/uuid"

type ExternalListMCPConnectionsRequest struct {
	AgentId         *uuid.UUID `query:"agentId"`
	SortBy          *string    `query:"sortBy"`
	SortDirection   *string    `query:"sortDirection"`
	Page            *uint      `query:"page"`
	ElementsPerPage *uint      `query:"elementsPerPage"`
	Search          *string    `query:"search"`
}
