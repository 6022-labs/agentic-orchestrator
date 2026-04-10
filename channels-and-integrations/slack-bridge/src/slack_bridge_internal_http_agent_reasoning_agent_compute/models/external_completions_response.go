package models

type ExternalAgentComputeResponses struct {
	IsFinal   bool      `json:"isFinal"`
	Responses []*string `json:"responses"`
}
