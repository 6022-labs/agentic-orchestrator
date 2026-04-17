package interfaces

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type BrokerMessageForwarderInterface interface {
	Forward(
		slackChannel *models.SlackChannel,
		parentMessageTs, messageTs, content string,
		authorSlackUser *models.SlackUser,
		authorSlackAgent *models.SlackAgent,
	) error
}
