package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
)

type AgentByCollectionTokenIdRequest struct {
	AgentCollectionAddress string `swaggerignore:"true" param:"agentCollectionAddress"`
	AgentCollectionTokenId uint64 `swaggerignore:"true" param:"agentCollectionTokenId"`
}

func (request *AgentByCollectionTokenIdRequest) ValidateAndSanitize() error {
	request.AgentCollectionAddress = strings.TrimSpace(request.AgentCollectionAddress)
	if !helpers.IsEthAddress(request.AgentCollectionAddress) {
		return errors.NewValidationError("agentCollectionAddress", "must be a 0x-prefixed 20-byte hex address")
	}
	// Token IDs start at 1 on-chain, so 0 can never identify an agent.
	if request.AgentCollectionTokenId == 0 {
		return errors.NewValidationError("agentCollectionTokenId", "must be a positive integer")
	}
	request.AgentCollectionAddress = strings.ToLower(request.AgentCollectionAddress)
	return nil
}
