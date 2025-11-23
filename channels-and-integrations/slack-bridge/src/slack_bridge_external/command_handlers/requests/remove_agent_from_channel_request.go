package requests

type RemoveAgentFromChannelRequest struct {
	ChannelID         string `form:"channel_id"`
	CollectionAddress string `form:"collection_address"`
	CollectionTokenId uint64 `form:"collection_token_id"`
}
