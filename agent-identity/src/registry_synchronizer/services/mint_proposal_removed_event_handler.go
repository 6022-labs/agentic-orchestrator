package services

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MintProposalRemovedEventHandler struct {
	logger                    *zap.Logger
	agentRepository           repositories.AgentRepositoryInterface
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface
}

func NewMintProposalRemovedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
) *MintProposalRemovedEventHandler {
	return &MintProposalRemovedEventHandler{
		logger:                    logger,
		agentRepository:           agentRepository,
		agentCollectionRepository: agentCollectionRepository,
	}
}

func (handler *MintProposalRemovedEventHandler) Handle(event *abi.AgentCollectionV1MintProposalRemoved) error {
	handler.logger.Info("Handling MintProposalRemoved event",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("proposalId", event.ProposalId.Uint64()),
	)

	agent, err := handler.agentRepository.GetNotYetMintedByCollectionAddressAndProposalId(event.Raw.Address.Hex(), event.ProposalId.Uint64())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if agent == nil {
		handler.logger.Warn(
			"MintProposalRemoved event has no matching agent record; skipping",
			zap.String("collectionAddress", event.Raw.Address.Hex()),
			zap.Uint64("proposalId", event.ProposalId.Uint64()),
		)
		return nil
	}

	err = handler.agentRepository.Remove(agent)
	if err != nil {
		return fmt.Errorf("failed to remove agent from proposal: %w", err)
	}

	return nil
}
