package models

import "github.com/google/uuid"

type ExternalAgentIdRequest struct {
	AgentId       uuid.UUID
	Search        *string `query:"search"`
	SortBy        *string `query:"sortBy"`
	SortDirection *string `query:"sortDirection"`
}
