package event_handlers

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"go.uber.org/zap"
)

type AgentSignerUpdatedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1AgentSignerUpdated) error
}

type AgentSignerUpdatedEventHandler struct {
	logger                 *zap.Logger
	orchestratorSettings   *settings.OrchestratorSettings
	finishAgentImportation use_cases.FinishAgentImportationInterface
}

func NewAgentSignerUpdatedEventHandler(
	logger *zap.Logger,
	orchestratorSettings *settings.OrchestratorSettings,
	finishAgentImportation use_cases.FinishAgentImportationInterface,
) *AgentSignerUpdatedEventHandler {
	return &AgentSignerUpdatedEventHandler{
		logger:                 logger,
		orchestratorSettings:   orchestratorSettings,
		finishAgentImportation: finishAgentImportation,
	}
}

func (h *AgentSignerUpdatedEventHandler) Handle(event *abi.AgentCollectionV1AgentSignerUpdated) error {
	if h.orchestratorSettings.PublicKey != event.NewSigner {
		h.logger.Info("Orchestrator is not the new signer", zap.String("newSigner", event.NewSigner.Hex()))
		return nil
	}

	request := &requests.FinishAgentImportationRequest{
		AgentCollectionAddress: event.Raw.Address.Hex(),
		AgentCollectionTokenId: event.TokenId.Uint64(),
	}
	return h.finishAgentImportation.Finish(request)
}
