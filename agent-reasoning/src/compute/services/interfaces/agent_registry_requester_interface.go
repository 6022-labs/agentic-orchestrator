package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases/responses"
)

type AgentRegistryRequesterInterface interface {
	GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*responses.AgentResponse, error)
	GetMintedAgentByEnsDomain(ensDomain string) (*responses.AgentResponse, error)
	GetAgentsByEnsDomains(ensDomains []string) ([]*responses.AgentResponse, error)
}
