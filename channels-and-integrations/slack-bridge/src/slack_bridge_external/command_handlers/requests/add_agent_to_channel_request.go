package requests

type AddAgentToChannelRequest struct {
	ChannelID string `form:"channel_id"`
	AgentCommandArguments
}
