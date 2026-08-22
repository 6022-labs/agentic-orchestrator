package services

import (
	"context"

	"encoding/json"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/clients"
	http_slack_models "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/models"
)

type SlackChannelRequester struct {
	client clients.SlackClientInterface
}

func NewSlackChannelRequester(
	client clients.SlackClientInterface,
) *SlackChannelRequester {
	return &SlackChannelRequester{
		client: client,
	}
}

func (s *SlackChannelRequester) GetChannelInformation(
	ctx context.Context,
	channelId string,
) (*interfaces.ChannelInformation, error) {
	channelInfoResponse, err := s.client.GetConversationInfo(ctx, &http_slack_models.ExternalGetConversationInfoRequest{
		Channel: channelId,
	})
	if err != nil {
		return nil, err
	}

	channelMembers, err := s.client.GetConversationMembers(
		ctx,
		&http_slack_models.ExternalGetConversationMembersRequest{
			Channel: channelId,
		},
	)
	if err != nil {
		return nil, err
	}
	if !channelMembers.Ok {
		jsonResp, _ := json.Marshal(channelMembers)
		return nil, fmt.Errorf("failed to get channel members: %s", string(jsonResp))
	}

	return &interfaces.ChannelInformation{
		Id:          channelId,
		UserIds:     channelMembers.Members,
		Name:        channelInfoResponse.Channel.Name,
		Description: channelInfoResponse.Channel.Purpose.Value,

		IsMember: channelInfoResponse.Channel.IsMember != nil && *channelInfoResponse.Channel.IsMember,
	}, nil
}

func (s *SlackChannelRequester) JoinChannel(ctx context.Context, channelId string) error {
	_, err := s.client.JoinConversation(ctx, &http_slack_models.ExternalJoinConversationRequest{
		Channel: channelId,
	})
	if err != nil {
		return err
	}

	return nil
}
