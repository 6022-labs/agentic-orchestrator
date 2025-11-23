package responses

type PaginateResponse struct {
	Page            uint `json:"page"`
	ElementsPerPage uint `json:"elementsPerPage"`
	TotalElements   uint `json:"totalElements"`
}
