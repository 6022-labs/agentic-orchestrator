package requests

type ListAgentSortBy string
type ListAgentSortByList []ListAgentSortBy

const (
	ListAgentSortByName              ListAgentSortBy = "name"
	ListAgentSortByOwner             ListAgentSortBy = "owner"
	ListAgentSortByCreatedAt         ListAgentSortBy = "createdAt"
	ListAgentSortByWalletAddress     ListAgentSortBy = "walletAddress"
	ListAgentSortByCollectionAddress ListAgentSortBy = "collectionAddress"
)

func GetAllListAgentSortBy() ListAgentSortByList {
	return ListAgentSortByList{
		ListAgentSortByName,
		ListAgentSortByOwner,
		ListAgentSortByCreatedAt,
		ListAgentSortByWalletAddress,
		ListAgentSortByCollectionAddress,
	}
}

func (s *ListAgentSortBy) ToDatabaseField() *string {
	value := "created_at"
	if s == nil {
		return &value
	}

	switch *s {
	case ListAgentSortByName:
		value = "name"
	case ListAgentSortByOwner:
		value = "owner"
	case ListAgentSortByCreatedAt:
		value = "created_at"
	case ListAgentSortByWalletAddress:
		value = "wallet_address"
	case ListAgentSortByCollectionAddress:
		value = "collection_address"
	default:
		value = "created_at"
	}

	return &value
}

func (l ListAgentSortByList) Contains(sortBy ListAgentSortBy) bool {
	for _, s := range l {
		if s == sortBy {
			return true
		}
	}
	return false
}

type ListAgentsRequest struct {
	PaginateRequest
	SortBy     *ListAgentSortBy `query:"sortBy"`
}

func (r *ListAgentsRequest) Sanitize() {
	r.PaginateRequest.Sanitize()

	if r.SortDirection == nil || !GetAllSortDirection().Contains(*r.SortDirection) {
		defaultSortDirection := SortDirectionAsc
		r.SortDirection = &defaultSortDirection
	}

	if r.SortBy == nil || !GetAllListAgentSortBy().Contains(*r.SortBy) {
		defaultSortBy := ListAgentSortByCreatedAt
		r.SortBy = &defaultSortBy
	}
}
