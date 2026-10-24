package requests

type RemoveAgentFromChannelRequest struct {
	ChannelID string `form:"channel_id"`
	AgentCommandArguments
}
