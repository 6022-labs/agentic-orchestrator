package interfaces

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type AgentRequesterInterface interface {
	GetAgent(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error)
}
