package services

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"go.uber.org/zap"
)

type TransferEventHandler struct {
	logger          *zap.Logger
	agentRepository repositories.AgentRepositoryInterface
}

func NewTransferEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
) *TransferEventHandler {
	return &TransferEventHandler{
		logger:          logger,
		agentRepository: agentRepository,
	}
}

func (handler *TransferEventHandler) Handle(event *abi.AgentCollectionV1Transfer) error {
	handler.logger.Info("Handling Transfer event for collection",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("tokenId", event.TokenId.Uint64()),
		zap.String("from", event.From.Hex()),
		zap.String("to", event.To.Hex()),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		event.Raw.Address.Hex(),
		event.TokenId.Uint64(),
	)
	if err != nil {
		return fmt.Errorf("failed to get agent by collection and token id: %w", err)
	}
	if agent == nil {
		return errors.New("agent not found for transfer event")
	}

	ownerHex := event.To.Hex()
	agent.Owner = &ownerHex

	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return err
	}

	return nil
}
