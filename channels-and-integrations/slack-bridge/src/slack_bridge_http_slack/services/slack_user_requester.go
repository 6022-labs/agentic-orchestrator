package services

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/clients"
	http_slack_models "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/models"
)

type SlackUserRequester struct {
	client clients.SlackClientInterface
}

func NewSlackUserRequester(
	client clients.SlackClientInterface,
) *SlackUserRequester {
	return &SlackUserRequester{
		client: client,
	}
}

func (s *SlackUserRequester) GetUserInformation(userId string) (*models.SlackUser, error) {
	userInfoResponse, err := s.client.GetUsersInfo(&http_slack_models.ExternalGetUsersInfoRequest{
		User: userId,
	})
	if err != nil {
		return nil, err
	}

	return &models.SlackUser{
		Id:       userId,
		Name:     userInfoResponse.User.Name,
		RealName: userInfoResponse.User.RealName,
	}, nil
}
