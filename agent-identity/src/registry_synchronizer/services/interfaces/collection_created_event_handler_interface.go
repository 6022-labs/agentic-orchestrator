package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"

type CollectionCreatedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionsManagerCollectionCreated) error
}
