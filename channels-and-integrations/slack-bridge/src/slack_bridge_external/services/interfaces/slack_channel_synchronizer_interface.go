package interfaces

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackChannelSynchronizerInterface interface {
	Sync(channelId string) (*models.SlackChannel, error)
}
