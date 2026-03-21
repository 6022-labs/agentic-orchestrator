package models

type ExternalListAgentsResponse struct {
	Page            uint                     `json:"page"`
	ElementsPerPage uint                     `json:"elementsPerPage"`
	TotalElements   uint                     `json:"totalElements"`
	Agents          []*ExternalAgentResponse `json:"agents"`
}
