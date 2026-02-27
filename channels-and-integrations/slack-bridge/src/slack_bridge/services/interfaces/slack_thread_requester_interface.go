package interfaces

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackThreadRequesterInterface interface {
	PushCompletionToThread(
		completion string,
		slackAgent *models.SlackAgent,
		slackChannelId string,
		slackThreadId string,
	) error
}
