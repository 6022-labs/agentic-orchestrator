package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/ethereum/go-ethereum/common"
)

// GetAgentPromptTemplateRequest identifies one agent-scoped template row.
type GetAgentPromptTemplateRequest struct {
	CollectionAddress string                       `json:"-" uri:"collectionAddress"`
	CollectionTokenId uint64                       `json:"-" uri:"collectionTokenId"`
	Key               enums.AgentPromptTemplateKey `json:"-" uri:"key"`
}

// ValidateAndSanitize lowercases the collection address so it matches the stored rows.
func (r *GetAgentPromptTemplateRequest) ValidateAndSanitize() error {
	r.CollectionAddress = strings.TrimSpace(r.CollectionAddress)
	if !strings.HasPrefix(r.CollectionAddress, "0x") || !common.IsHexAddress(r.CollectionAddress) {
		return errors.NewValidationError("collectionAddress", "must be a 0x-prefixed hex address")
	}
	// Token IDs start at 1 on-chain, so 0 can never identify an agent.
	if r.CollectionTokenId == 0 {
		return errors.NewValidationError("collectionTokenId", "must be a positive integer")
	}
	if err := r.Key.Set(string(r.Key)); err != nil {
		return errors.NewValidationError("key", "must be one of: "+enums.GetAllAgentPromptTemplateKeys().String())
	}

	r.CollectionAddress = strings.ToLower(r.CollectionAddress)
	return nil
}
