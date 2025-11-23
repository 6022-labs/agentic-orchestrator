package requests

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type AgentByIdRequest struct {
	AuthenticatedWallet common.Address `swaggerignore:"true"`
	AgentId             uuid.UUID      `json:"agentId"`
}

func (r *AgentByIdRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "Agent ID is required")
	}

	return nil
}
