package services

import (
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
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if agent == nil {
		return fmt.Errorf("no agent found for collection address %s and proposal id %d", event.Raw.Address.Hex(), event.ProposalId.Uint64())
	}

	err = handler.agentRepository.Remove(agent)
	if err != nil {
		return fmt.Errorf("failed to remove agent from proposal: %w", err)
	}

	return nil
}
