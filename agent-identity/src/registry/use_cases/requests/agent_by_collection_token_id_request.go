package requests

type AgentByCollectionTokenIdRequest struct {
	AgentCollectionAddress string `param:"agentCollectionAddress"`
	AgentCollectionTokenId uint64 `param:"agentCollectionTokenId"`
}
