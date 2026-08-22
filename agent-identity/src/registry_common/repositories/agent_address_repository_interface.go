package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/google/uuid"
)

// AgentAddressRepositoryInterface persists the agents' typed on-chain addresses,
// one row per (agent_id, type); values are stored canonical and matched exactly.
type AgentAddressRepositoryInterface interface {
	// GetByTypeAndAddress returns the row holding the typed address.
	GetByTypeAndAddress(
		ctx context.Context,
		addressType enums.AddressType,
		address string,
	) (*models.AgentAddress, error)
	// GetAllByAgentId returns the agent's address rows.
	GetAllByAgentId(ctx context.Context, agentId uuid.UUID) ([]*models.AgentAddress, error)
	// GetAllByAgentIds returns the address rows of all the given agents.
	GetAllByAgentIds(ctx context.Context, agentIds []uuid.UUID) ([]*models.AgentAddress, error)
	// ReplaceAllForAgent replaces the agent's rows with the given addresses.
	ReplaceAllForAgent(ctx context.Context, agentId uuid.UUID, addresses map[enums.AddressType]string) error
}
