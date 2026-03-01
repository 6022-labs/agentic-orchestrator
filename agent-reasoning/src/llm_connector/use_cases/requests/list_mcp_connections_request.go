package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/errors"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ListAgentMCPConnectionsSortBy string
type ListAgentMCPConnectionsSortByList []ListAgentMCPConnectionsSortBy

const (
	ListAgentMCPConnectionsSortByUrl  ListAgentMCPConnectionsSortBy = "url"
	ListAgentMCPConnectionsSortByName ListAgentMCPConnectionsSortBy = "name"
)

func (s *ListAgentMCPConnectionsSortBy) ToDatabaseField() *string {
	if s == nil {
		return nil
	}

	var value *string

	switch *s {
	case ListAgentMCPConnectionsSortByUrl:
		value = lo.ToPtr("url")
	case ListAgentMCPConnectionsSortByName:
		value = lo.ToPtr("name")
	default:
		value = nil
	}

	return value
}

type ListMCPConnectionsRequest struct {
	PaginateRequest
	AgentId *uuid.UUID                     `query:"agentId"`
	Search  *string                        `query:"search"`
	SortBy  *ListAgentMCPConnectionsSortBy `query:"sortBy"`
}

func (r *ListMCPConnectionsRequest) SanitizeAndValidate() error {
	r.PaginateRequest.Sanitize()

	if r.Page != nil && *r.Page < 1 {
		return errors.NewValidationError("page", "must be greater than 0")
	}

	if r.ElementsPerPage != nil && *r.ElementsPerPage < 1 {
		return errors.NewValidationError("elementsPerPage", "must be greater than 0")
	}
	if r.ElementsPerPage != nil && *r.ElementsPerPage > uint(100) {
		return errors.NewValidationError("elementsPerPage", "must be less than or equal to 100")
	}

	if r.SortDirection != nil && !GetAllSortDirection().Contains(*r.SortDirection) {
		return errors.NewValidationError("sortDirection", "is not a valid sort direction")
	}

	return nil
}
