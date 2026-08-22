package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type SlackUserSynchronizerCachedInterface interface {
	Sync(ctx context.Context, userId string) (*models.SlackUser, error)
}
