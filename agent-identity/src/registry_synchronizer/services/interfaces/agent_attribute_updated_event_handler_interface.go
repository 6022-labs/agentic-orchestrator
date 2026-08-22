package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
)

type AgentAttributeUpdatedEventHandlerInterface interface {
	Handle(ctx context.Context, event *abi.AgentCollectionV1AgentAttributeUpdated) error
}
