package requests

type UpdateAgentAttributeRequest struct {
	CollectionTokenId uint64
	CollectionAddress string
	Key               string
	NewValue          string
}
