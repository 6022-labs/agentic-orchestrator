package interfaces

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type AgentCompletionsRequesterInterface interface {
	Completions(
		slackChannel *models.SlackChannel,
		slackThread *models.SlackThread,
		slackMessages []*models.SlackMessage,
		collectionAddress string,
		collectionTokenId uint64,
		forceCompletion bool,
	) ([]*string, *bool, error)
}
