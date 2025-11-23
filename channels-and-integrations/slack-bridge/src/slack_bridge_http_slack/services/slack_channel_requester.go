package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/clients"
	http_slack_models "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/models"
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

func (s *SlackChannelRequester) GetChannelInformation(channelId string) (*interfaces.ChannelInformation, error) {
	channelInfoResponse, err := s.client.GetConversationInfo(&http_slack_models.ExternalGetConversationInfoRequest{
		Channel: channelId,
	})
	if err != nil {
		return nil, err
	}

	channelMembers, err := s.client.GetConversationMembers(&http_slack_models.ExternalGetConversationMembersRequest{
		Channel: channelId,
	})
	if err != nil {
		return nil, err
	}
	if !channelMembers.Ok {
		jsonResp, _ := json.Marshal(channelMembers)
		return nil, fmt.Errorf("failed to get channel members: %s", string(jsonResp))
	}

	members := make([]string, 0)
	for _, member := range channelMembers.Members {
		userInfo, err := s.client.GetUsersInfo(&http_slack_models.ExternalGetUsersInfoRequest{
			User: member,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get user info: %w", err)
		}
		if !userInfo.Ok {
			return nil, fmt.Errorf("failed to get user info: %v", userInfo)
		}

		if len(strings.TrimSpace(userInfo.User.Profile.DisplayName)) != 0 {
			members = append(members, userInfo.User.Profile.DisplayName)
		}
	}

	return &interfaces.ChannelInformation{
		Id:          channelId,
		Members:     members,
		Name:        channelInfoResponse.Channel.Name,
		Description: channelInfoResponse.Channel.Purpose.Value,

		IsMember: channelInfoResponse.Channel.IsMember != nil && *channelInfoResponse.Channel.IsMember,
	}, nil
}

func (s *SlackChannelRequester) JoinChannel(channelId string) error {
	_, err := s.client.JoinConversation(&http_slack_models.ExternalJoinConversationRequest{
		Channel: channelId,
	})
	if err != nil {
		return err
	}

	return nil
}
