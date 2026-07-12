package responses

type ListBrokerConversationsResponse struct {
	Page            uint                          `json:"page"`
	ElementsPerPage uint                          `json:"elementsPerPage"`
	TotalElements   uint                          `json:"totalElements"`
	Conversations   []*BrokerConversationResponse `json:"conversations"`
}
