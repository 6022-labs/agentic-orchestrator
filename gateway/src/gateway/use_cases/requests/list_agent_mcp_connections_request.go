package requests

import "github.com/samber/lo"

type ListAgentMCPConnectionsSortBy string
type ListAgentMCPConnectionsSortByList []ListAgentMCPConnectionsSortBy

const (
	ListAgentMCPConnectionsSortByUrl  ListAgentMCPConnectionsSortBy = "url"
	ListAgentMCPConnectionsSortByName ListAgentMCPConnectionsSortBy = "name"
)

func GetAllListAgentMCPConnectionsSortBy() ListAgentMCPConnectionsSortByList {
	return ListAgentMCPConnectionsSortByList{
		ListAgentMCPConnectionsSortByUrl,
		ListAgentMCPConnectionsSortByName,
	}
}

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

type ListAgentMCPConnectionsRequest struct {
	AgentByCollectionTokenIdRequest
	PaginateRequest
	Search *string                        `query:"search"`
	SortBy *ListAgentMCPConnectionsSortBy `query:"sortBy"`
}

func (r *ListAgentMCPConnectionsRequest) SanitizeAndValidate() error {
	r.Sanitize()

	err := r.AgentByCollectionTokenIdRequest.EnsureAuthenticatedWallet()
	if err != nil {
		return err
	}

	return nil
}
