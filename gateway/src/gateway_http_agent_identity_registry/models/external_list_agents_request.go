package models

type ExternalListAgentsRequest struct {
	Page            *uint
	ElementsPerPage *uint
	SortBy          *string
	SortDirection   *string
}
