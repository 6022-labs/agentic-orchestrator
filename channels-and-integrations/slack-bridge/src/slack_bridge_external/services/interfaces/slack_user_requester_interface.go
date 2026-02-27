package interfaces

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackUserRequesterInterface interface {
	GetUserInformation(userId string) (*models.SlackUser, error)
}
