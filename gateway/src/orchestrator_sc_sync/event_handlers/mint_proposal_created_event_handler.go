package event_handlers

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"go.uber.org/zap"
)

type MintProposalCreatedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1MintProposalCreated) error
}

type MintProposalCreatedEventHandler struct {
	logger                     *zap.Logger
	pairProposalToAgentHandler use_cases.PairProposalToAgentInterface
}

func NewMintProposalCreatedEventHandler(
	logger *zap.Logger,
	pairProposalToAgentHandler use_cases.PairProposalToAgentInterface,
) *MintProposalCreatedEventHandler {
	return &MintProposalCreatedEventHandler{
		logger:                     logger,
		pairProposalToAgentHandler: pairProposalToAgentHandler,
	}
}

func (h *MintProposalCreatedEventHandler) Handle(event *abi.AgentCollectionV1MintProposalCreated) error {
	h.logger.Info("Handling MintProposalCreated event", zap.Any("event", event))

	err := h.pairProposalToAgentHandler.Pair(&requests.PairProposalToAgentRequest{
		ProposalId:        event.ProposalId.Uint64(),
		CollectionAddress: event.Raw.Address.Hex(),
	})
	if err != nil {
		h.logger.Error("Failed to pair proposal to agent", zap.Error(err))
		return err
	}

	h.logger.Info("Successfully paired proposal to agent",
		zap.Uint64("proposalID", event.ProposalId.Uint64()),
		zap.String("collectionAddress", event.Raw.Address.Hex()),
	)
	return nil
}
