package requests

type ListAgentContextKeysSortBy string
type ListAgentContextKeysSortByList []ListAgentContextKeysSortBy

const (
	ListAgentContextKeysSortByKey ListAgentContextKeysSortBy = "key"
)

func GetAllListAgentContextKeysSortBy() ListAgentContextKeysSortByList {
	return ListAgentContextKeysSortByList{
		ListAgentContextKeysSortByKey,
	}
}

func (l ListAgentContextKeysSortByList) Contains(sortBy ListAgentContextKeysSortBy) bool {
	for _, s := range l {
		if s == sortBy {
			return true
		}
	}
	return false
}

type ListAgentContextKeysRequest struct {
	PaginateRequest
	AgentByCollectionTokenIdRequest
	Search *string                     `query:"search"`
	SortBy *ListAgentContextKeysSortBy `query:"sortBy"`
}

func (r *ListAgentContextKeysRequest) Sanitize() {
	r.PaginateRequest.Sanitize()

	if r.SortDirection == nil || !GetAllSortDirection().Contains(*r.SortDirection) {
		defaultSortDirection := SortDirectionAsc
		r.SortDirection = &defaultSortDirection
	}

	if r.SortBy == nil {
		defaultSortBy := ListAgentContextKeysSortByKey
		r.SortBy = &defaultSortBy
	}

	if !GetAllListAgentContextKeysSortBy().Contains(*r.SortBy) {
		defaultSortBy := ListAgentContextKeysSortByKey
		r.SortBy = &defaultSortBy
	}
}
