package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"

type TransferEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1Transfer) error
}
