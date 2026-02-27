package requests

import (
	"errors"
	"strconv"
	"strings"

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

func (r *CommandRequest) ToAddAgentToChannelRequest() (*requests.AddAgentToChannelRequest, error) {
	// Usage hint: [collection address] [token ID]
	parameters := strings.Split(r.Text, " ")
	collectionAddress := parameters[0]
	collectionTokenId, err := strconv.ParseUint(parameters[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return &requests.AddAgentToChannelRequest{
		ChannelID:         r.ChannelID,
		CollectionAddress: collectionAddress,
		CollectionTokenId: collectionTokenId,
	}, nil
}

func (r *CommandRequest) ToRemoveAgentFromChannelRequest() (*requests.RemoveAgentFromChannelRequest, error) {
	// Usage hint: [collection address] [token ID]
	parameters := strings.Split(r.Text, " ")
	collectionAddress := parameters[0]
	collectionTokenId, err := strconv.ParseUint(parameters[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return &requests.RemoveAgentFromChannelRequest{
		ChannelID:         r.ChannelID,
		CollectionAddress: collectionAddress,
		CollectionTokenId: collectionTokenId,
	}, nil
}

func (r *CommandRequest) ToVerifyWalletOwnershipRequest() *requests.VerifyWalletOwnershipRequest {
	// Usage hint: [no extra parameters for now]
	return &requests.VerifyWalletOwnershipRequest{
		UserId: r.UserID,
	}
}

func (r *CommandRequest) ToPairSlackUserToHumanAgentRequest() (*requests.PairSlackUserToHumanAgentRequest, error) {
	// Usage hint: [slack user id] [collection address] [token ID]
	parameters := strings.Split(r.Text, " ")
	if len(parameters) < 3 {
		return nil, errors.New("invalid command format")
	}

	slackUserId := parameters[0]
	if len(strings.TrimSpace(slackUserId)) == 0 {
		return nil, errors.New("invalid command format")
	}

	collectionAddress := parameters[1]
	if len(strings.TrimSpace(collectionAddress)) == 0 {
		return nil, errors.New("invalid command format")
	}

	collectionTokenId, err := strconv.ParseUint(parameters[2], 10, 64)
	if err != nil {
		return nil, errors.New("invalid command format")
	}

	return &requests.PairSlackUserToHumanAgentRequest{
		RequesterSlackUserId:        r.UserID,
		SlackUserId:                 slackUserId,
		HumanAgentCollectionAddress: collectionAddress,
		HumanAgentCollectionTokenId: collectionTokenId,
	}, nil
}
