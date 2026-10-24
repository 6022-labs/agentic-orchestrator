package requests

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
)

type CommandRequest struct {
	ChannelID string `form:"channel_id"`
	Text      string `form:"text"`
	UserID    string `form:"user_id"`
}

func (r *CommandRequest) ToListAgentInChannelRequest() *requests.ListAgentInChannelRequest {
	return &requests.ListAgentInChannelRequest{
		ChannelID: r.ChannelID,
	}
}

func (r *CommandRequest) ToAddAgentToChannelRequest() *requests.AddAgentToChannelRequest {
	request := &requests.AddAgentToChannelRequest{ChannelID: r.ChannelID}
	request.TryParse(r.Text)
	return request
}

func (r *CommandRequest) ToRemoveAgentFromChannelRequest() *requests.RemoveAgentFromChannelRequest {
	request := &requests.RemoveAgentFromChannelRequest{ChannelID: r.ChannelID}
	request.TryParse(r.Text)
	return request
}

func (r *CommandRequest) ToVerifyWalletOwnershipRequest() *requests.VerifyWalletOwnershipRequest {
	// Usage hint: [no extra parameters for now]
	return &requests.VerifyWalletOwnershipRequest{
		UserId: r.UserID,
	}
}

func (r *CommandRequest) ToPairSlackUserToHumanAgentRequest() *requests.PairSlackUserToHumanAgentRequest {
	request := &requests.PairSlackUserToHumanAgentRequest{RequesterSlackUserId: r.UserID}
	request.TryParse(r.Text)
	return request
}
