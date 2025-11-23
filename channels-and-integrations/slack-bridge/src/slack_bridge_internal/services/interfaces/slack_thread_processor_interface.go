package interfaces

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackThreadProcessorInterface interface {
	ProcessThread(slackAgent *models.SlackAgent, slackThread *models.SlackThread) error
}
