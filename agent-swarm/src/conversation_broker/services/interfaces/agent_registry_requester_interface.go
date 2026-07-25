package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"

// AgentRegistryRequesterInterface is the broker's outbound port to the
// agent-identity registry. Returns nil (no error) when the agent is not found.
type AgentRegistryRequesterInterface interface {
	GetAgent(collectionAddress string, collectionTokenId uint64) (*models.RegistryAgent, error)
	GetMintedAgentByEnsDomain(ensDomain string) (*models.RegistryAgent, error)
}
