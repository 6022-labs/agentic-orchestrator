package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type SlackThreadRequesterInterface interface {
	PushCompletionToThread(
		ctx context.Context,
		completion string,
		slackAgent *models.SlackAgent,
		slackChannelId string,
		slackThreadId string,
	) error
}
