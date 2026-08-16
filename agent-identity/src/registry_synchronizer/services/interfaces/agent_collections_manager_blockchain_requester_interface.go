package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
)

type AgentCollectionsManagerBlockchainRequesterInterface interface {
	GetAgentCollectionsManagerCollectionCreatedEvent(
		from, to uint64,
	) ([]*abi.AgentCollectionsManagerCollectionCreated, error)
}
