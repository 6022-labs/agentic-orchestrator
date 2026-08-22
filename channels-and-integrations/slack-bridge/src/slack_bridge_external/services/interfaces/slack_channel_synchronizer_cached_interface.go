package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type SlackChannelSynchronizerCachedInterface interface {
	Sync(ctx context.Context, channelId string) (*models.SlackChannel, error)
}
