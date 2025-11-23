package requests

type UpdateAgentOwnerRequest struct {
	CollectionTokenId uint64
	CollectionAddress string
	NewOwner          string
}
