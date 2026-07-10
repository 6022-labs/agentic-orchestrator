package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

// RemoteAgentCardResolverInterface resolves an agent's ENS domain to its
// published 6022 card, caching results with a configurable TTL.
type RemoteAgentCardResolverInterface interface {
	ResolveCard(ensDomain string) (*models.RemoteAgentCard, error)
}
