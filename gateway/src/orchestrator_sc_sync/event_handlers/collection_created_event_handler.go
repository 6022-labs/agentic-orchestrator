package event_handlers

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"go.uber.org/zap"
)

type CollectionCreatedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionsManagerCollectionCreated) error
}

type CollectionCreatedEventHandler struct {
	logger                *zap.Logger
	orchestratorSettings  *settings.OrchestratorSettings
	createAgentCollection use_cases.CreateAgentCollectionInterface
}

func NewCollectionCreatedEventHandler(
	logger *zap.Logger,
	orchestratorSettings *settings.OrchestratorSettings,
	createAgentCollection use_cases.CreateAgentCollectionInterface,
) *CollectionCreatedEventHandler {
	return &CollectionCreatedEventHandler{
		logger:                logger,
		orchestratorSettings:  orchestratorSettings,
		createAgentCollection: createAgentCollection,
	}
}

func (h *CollectionCreatedEventHandler) Handle(event *abi.AgentCollectionsManagerCollectionCreated) error {
	request := &requests.CreateAgentCollectionRequest{
		CollectionAddress: event.CollectionAddress.Hex(),
		CollectionName:    event.Name,
		CollectionSymbol:  event.Symbol,
	}
	return h.createAgentCollection.Create(request)
}
