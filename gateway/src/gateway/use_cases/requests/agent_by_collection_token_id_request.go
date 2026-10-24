package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/ethereum/go-ethereum/common"
)

type AgentByCollectionTokenIdRequest struct {
	AuthenticatedWallet    *common.Address `swaggerignore:"true"`
	AgentCollectionAddress string          `swaggerignore:"true" param:"agentCollectionAddress"`
	AgentCollectionTokenId uint64          `swaggerignore:"true" param:"agentCollectionTokenId"`
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

func (request *AgentByCollectionTokenIdRequest) EnsureAuthenticatedWallet() error {
	if request.AuthenticatedWallet == nil {
		return errors.NewValidationError("auth", "is required")
	}

	if *request.AuthenticatedWallet == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return errors.NewValidationError("auth", "is invalid")
	}

	return nil
}
