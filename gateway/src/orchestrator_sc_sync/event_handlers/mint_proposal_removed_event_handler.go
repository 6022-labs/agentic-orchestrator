package event_handlers

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"go.uber.org/zap"
)

type MintProposalRemovedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1MintProposalRemoved) error
}

type MintProposalRemovedEventHandler struct {
	logger                  *zap.Logger
	removeAgentFromProposal use_cases.RemoveAgentFromProposalInterface
}

func NewMintProposalRemovedEventHandler(
	logger *zap.Logger,
	removeAgentFromProposal use_cases.RemoveAgentFromProposalInterface,
) *MintProposalRemovedEventHandler {
	return &MintProposalRemovedEventHandler{
		logger:                  logger,
		removeAgentFromProposal: removeAgentFromProposal,
	}
}

func (h *MintProposalRemovedEventHandler) Handle(event *abi.AgentCollectionV1MintProposalRemoved) error {
	h.logger.Info("Handling MintProposalRemoved event", zap.Any("event", event))

	err := h.removeAgentFromProposal.Remove(&requests.RemoveAgentFromProposalRequest{
		ProposalId:        event.ProposalId.Uint64(),
		CollectionAddress: event.Raw.Address.Hex(),
	})
	if err != nil {
		return err
	}

	h.logger.Info("Successfully removed agent from proposal",
		zap.Uint64("proposalID", event.ProposalId.Uint64()),
		zap.String("collectionAddress", event.Raw.Address.Hex()),
	)
	return nil
}
