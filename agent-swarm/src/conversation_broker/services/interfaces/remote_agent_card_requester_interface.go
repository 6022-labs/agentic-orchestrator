package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"

// RemoteAgentCardRequesterInterface fetches a node's /.well-known/6022 card.
type RemoteAgentCardRequesterInterface interface {
	// GetCard also returns the body verbatim — signature verification runs over those bytes.
	GetCard(agentUrl string) (*models.RemoteAgentCard, []byte, error)
}
