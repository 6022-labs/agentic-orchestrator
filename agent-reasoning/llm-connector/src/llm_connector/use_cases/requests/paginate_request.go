package requests

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

type PaginateRequest struct {
	Page            *uint          `query:"page"`
	SortDirection   *SortDirection `query:"sortDirection"`
	ElementsPerPage *uint          `query:"elementsPerPage"`
}

func (r *PaginateRequest) Sanitize() {
	if r.Page == nil {
		defaultPage := uint(1)
		r.Page = &defaultPage
	}

	if r.ElementsPerPage == nil {
		defaultElementsPerPage := uint(10)
		r.ElementsPerPage = &defaultElementsPerPage
	}
}
