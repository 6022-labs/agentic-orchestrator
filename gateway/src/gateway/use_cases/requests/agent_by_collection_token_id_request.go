package requests

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/ethereum/go-ethereum/common"
)

type AgentByCollectionTokenIdRequest struct {
	AuthenticatedWallet    *common.Address `swaggerignore:"true"`
	AgentCollectionAddress string          `swaggerignore:"true" param:"agentCollectionAddress"`
	AgentCollectionTokenId uint64          `swaggerignore:"true" param:"agentCollectionTokenId"`
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
