package models

type ExternalAgentByCollectionTokenIdRequest struct {
	AgentCollectionAddress string `json:"agentCollectionAddress"`
	AgentCollectionTokenId uint64 `json:"agentCollectionTokenId"`
}
