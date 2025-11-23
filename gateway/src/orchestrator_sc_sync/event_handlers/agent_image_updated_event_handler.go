package event_handlers

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"go.uber.org/zap"
)

type AgentImageUpdatedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1AgentImageUpdated) error
}

type AgentImageUpdatedEventHandler struct {
	logger                  *zap.Logger
	upsertAgentImageHandler use_cases.UpsertAgentImageInterface
}

func NewAgentImageUpdatedEventHandler(
	logger *zap.Logger,
	upsertAgentImageHandler use_cases.UpsertAgentImageInterface,
) *AgentImageUpdatedEventHandler {
	return &AgentImageUpdatedEventHandler{
		logger:                  logger,
		upsertAgentImageHandler: upsertAgentImageHandler,
	}
}

func (h *AgentImageUpdatedEventHandler) Handle(event *abi.AgentCollectionV1AgentImageUpdated) error {
	h.logger.Info("Handling AgentImageUpdated event", zap.Any("event", event))

	err := h.upsertAgentImageHandler.Upsert(&requests.UpsertAgentImageRequest{
		CollectionAddress: event.Raw.Address.Hex(),
		CollectionTokenId: event.TokenId.Uint64(),
		Key:               event.Key.String(),
		NewValue:          event.NewValue,
	})
	if err != nil {
		h.logger.Error("Failed to upsert agent image", zap.Error(err))
		return err
	}

	h.logger.Info("Successfully upserted agent image",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("collectionTokenId", event.TokenId.Uint64()),
		zap.String("key", event.Key.String()),
		zap.String("newValue", event.NewValue),
	)
	return nil
}
