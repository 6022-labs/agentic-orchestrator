package responses

type ListAgentContextKeysResponse struct {
	PaginateResponse
	ContextKeys []string `json:"contextKeys"`
}
