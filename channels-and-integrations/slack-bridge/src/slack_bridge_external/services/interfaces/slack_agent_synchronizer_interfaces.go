package interfaces

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackAgentSynchronizerInterface interface {
	Sync(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error)
}
