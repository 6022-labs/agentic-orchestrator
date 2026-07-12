package responses

type AgentPromptTemplateResponse struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionTokenId uint64 `json:"collectionTokenId"`
	Key               string `json:"key"`
	Content           string `json:"content"`
}
