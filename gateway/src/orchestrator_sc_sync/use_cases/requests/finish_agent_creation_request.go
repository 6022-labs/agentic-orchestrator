package requests

type FinishAgentCreationRequest struct {
	CollectionTokenId uint64 `json:"collectionTokenId"`
	CollectionAddress string `json:"collectionAddress"`
	Signer            string `json:"signer"`
}
