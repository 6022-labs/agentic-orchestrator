package requests

type ListConversationsRequest struct {
	Page            *uint `query:"page"`
	ElementsPerPage *uint `query:"elementsPerPage"`
}

func (r *ListConversationsRequest) Sanitize() {
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
