package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"

type MintProposalCreatedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1MintProposalCreated) error
}
