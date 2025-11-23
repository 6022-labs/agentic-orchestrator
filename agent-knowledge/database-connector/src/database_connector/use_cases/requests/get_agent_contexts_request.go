package requests

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/google/uuid"
)

type GetAgentContextsSortBy string
type GetAgentContextsSortByList []GetAgentContextsSortBy
type SortDirection string
type SortDirectionList []SortDirection

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

func GetAllSortDirection() SortDirectionList {
	return SortDirectionList{
		SortDirectionAsc,
		SortDirectionDesc,
	}
}

func (s *SortDirection) ToDatabaseDirection() *string {
	value := "asc"
	if s == nil {
		return &value
	}

	switch *s {
	case SortDirectionAsc:
		value = "asc"
	case SortDirectionDesc:
		value = "desc"
	default:
		value = "asc"
	}

	return &value
}

func (l SortDirectionList) Contains(sortDirection SortDirection) bool {
	for _, s := range l {
		if s == sortDirection {
			return true
		}
	}
	return false
}

const (
	GetAgentContextsSortByKey GetAgentContextsSortBy = "key"
)

func GetAllGetAgentContextsSortBy() GetAgentContextsSortByList {
	return GetAgentContextsSortByList{
		GetAgentContextsSortByKey,
	}
}

func (l GetAgentContextsSortByList) Contains(sortBy GetAgentContextsSortBy) bool {
	for _, s := range l {
		if s == sortBy {
			return true
		}
	}
	return false
}

func (s *GetAgentContextsSortBy) ToDatabaseField() *string {
	value := "key"
	if s == nil {
		return &value
	}

	switch *s {
	case GetAgentContextsSortByKey:
		value = "key"
	default:
		value = "key"
	}

	return &value
}

type GetAgentContextsRequest struct {
	AgentId       uuid.UUID               `param:"agentId" swaggerignore:"true"`
	Search        *string                 `query:"search"`
	SortBy        *GetAgentContextsSortBy `query:"sortBy"`
	SortDirection *SortDirection          `query:"sortDirection"`
}

func (r *GetAgentContextsRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "must be provided")
	}

	return nil
}

func (r *GetAgentContextsRequest) Sanitize() {
	if r.SortBy == nil {
		defaultSortBy := GetAgentContextsSortByKey
		r.SortBy = &defaultSortBy
	}

	if !GetAllGetAgentContextsSortBy().Contains(*r.SortBy) {
		defaultSortBy := GetAgentContextsSortByKey
		r.SortBy = &defaultSortBy
	}

	if r.SortDirection == nil {
		defaultSortDirection := SortDirectionAsc
		r.SortDirection = &defaultSortDirection
	}

	if !GetAllSortDirection().Contains(*r.SortDirection) {
		defaultSortDirection := SortDirectionAsc
		r.SortDirection = &defaultSortDirection
	}
}
