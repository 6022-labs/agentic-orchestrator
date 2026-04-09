package requests

import (
	registry_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type AgentByIdRequest struct {
	AuthenticatedWalletRaw string         `json:"authenticatedWallet" form:"authenticatedWallet"`
	AuthenticatedWallet    common.Address `json:"-"`
	AgentId                uuid.UUID      `json:"agentId"`
}

func (r *AgentByIdRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return registry_errors.NewValidationError("agentId", "Agent ID is required")
	}
	if r.AuthenticatedWallet == (common.Address{}) {
		return registry_errors.NewValidationError("authenticatedWallet", "is required")
	}
	return nil
}
