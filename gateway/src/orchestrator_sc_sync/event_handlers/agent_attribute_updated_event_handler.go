package event_handlers

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"go.uber.org/zap"
)

type AgentAttributeUpdatedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1AgentAttributeUpdated) error
}

type AgentAttributeUpdatedEventHandler struct {
	logger                   *zap.Logger
	orchestratorSettings     *settings.OrchestratorSettings
	updateAgentAttribute     use_cases.UpdateAgentAttributeInterface
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface
}

func NewAgentAttributeUpdatedEventHandler(
	logger *zap.Logger,
	orchestratorSettings *settings.OrchestratorSettings,
	updateAgentAttribute use_cases.UpdateAgentAttributeInterface,
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
) *AgentAttributeUpdatedEventHandler {
	return &AgentAttributeUpdatedEventHandler{
		logger:                   logger,
		orchestratorSettings:     orchestratorSettings,
		updateAgentAttribute:     updateAgentAttribute,
		agentCollectionRequester: agentCollectionRequester,
	}
}

func (h *AgentAttributeUpdatedEventHandler) Handle(event *abi.AgentCollectionV1AgentAttributeUpdated) error {
	_, key, _, err := h.agentCollectionRequester.GetAddOrUpdateAgentAttributeInputsFromTransaction(event.Raw.Address, event.Raw.TxHash)
	if err != nil {
		return fmt.Errorf("failed to get addOrUpdateAgentAttribute inputs from transaction: %w", err)
	}

	request := &requests.UpdateAgentAttributeRequest{
		CollectionAddress: event.Raw.Address.Hex(),
		CollectionTokenId: event.TokenId.Uint64(),
		Key:               *key,
		NewValue:          event.NewValue,
	}

	return h.updateAgentAttribute.Update(request)
}
