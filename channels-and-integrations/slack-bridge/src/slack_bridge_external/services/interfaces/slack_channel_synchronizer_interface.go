package interfaces

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackChannelSynchronizerInterface interface {
	Sync(channelId string) (*models.SlackChannel, error)
}
