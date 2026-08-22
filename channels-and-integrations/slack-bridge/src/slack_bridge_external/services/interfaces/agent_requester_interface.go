package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type AgentRequesterInterface interface {
	GetAgent(ctx context.Context, collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error)
}
