package requests

type CreateAgentCollectionRequest struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionName    string `json:"collectionName"`
	CollectionSymbol  string `json:"collectionSymbol"`
}
