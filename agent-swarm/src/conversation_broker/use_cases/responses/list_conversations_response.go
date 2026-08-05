package responses

type ListConversationsResponse struct {
	Page            uint                    `json:"page"`
	ElementsPerPage uint                    `json:"elementsPerPage"`
	TotalElements   uint                    `json:"totalElements"`
	Conversations   []*ConversationResponse `json:"conversations"`
}
