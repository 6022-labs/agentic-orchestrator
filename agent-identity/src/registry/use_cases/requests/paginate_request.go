package requests

import "github.com/samber/lo"

type SortDirection string
type SortDirectionList []SortDirection

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

func GetAllSortDirection() SortDirectionList {
	return SortDirectionList{SortDirectionAsc, SortDirectionDesc}
}

func (s *SortDirection) ToDatabaseDirection() *string {
	if s == nil {
		return nil
	}

	switch *s {
	case SortDirectionAsc:
		return lo.ToPtr("asc")
	case SortDirectionDesc:
		return lo.ToPtr("desc")
	default:
		return nil
	}
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
	if r.Page == nil || *r.Page < 1 {
		defaultPage := uint(1)
		r.Page = &defaultPage
	}

	if r.ElementsPerPage == nil || *r.ElementsPerPage < 1 {
		defaultElementsPerPage := uint(10)
		r.ElementsPerPage = &defaultElementsPerPage
	}

	if *r.ElementsPerPage > uint(100) {
		maxElementsPerPage := uint(100)
		r.ElementsPerPage = &maxElementsPerPage
	}
}
