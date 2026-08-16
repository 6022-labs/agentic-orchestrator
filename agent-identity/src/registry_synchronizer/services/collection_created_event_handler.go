package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type CollectionCreatedEventHandler struct {
	logger *zap.Logger

	collectionListeners []interfaces.CollectionEventSubscriberInterface

	agentCollectionRepository repositories.AgentCollectionRepositoryInterface
}

type NewCollectionCreatedEventHandlerParams struct {
	dig.In

	Logger *zap.Logger

	AgentCollectionRepository repositories.AgentCollectionRepositoryInterface

	CollectionListeners []interfaces.CollectionEventSubscriberInterface `group:"collection_event_subscribers"`
}

func NewCollectionCreatedEventHandler(
	params NewCollectionCreatedEventHandlerParams,
) *CollectionCreatedEventHandler {
	return &CollectionCreatedEventHandler{
		logger:                    params.Logger,
		collectionListeners:       params.CollectionListeners,
		agentCollectionRepository: params.AgentCollectionRepository,
	}
}

func (handler *CollectionCreatedEventHandler) Handle(event *abi.AgentCollectionsManagerCollectionCreated) error {
	handler.logger.Info("Handling CollectionCreated event",
		zap.String("collectionAddress", event.CollectionAddress.Hex()),
		zap.Uint64("collectionId", event.CollectionId.Uint64()),
		zap.String("name", event.Name),
		zap.String("symbol", event.Symbol),
	)

	collection := &models.AgentCollection{
		Address: event.CollectionAddress.Hex(),
		Name:    event.Name,
		Symbol:  event.Symbol,
	}

	_, err := handler.agentCollectionRepository.Save(collection)
	if err != nil {
		return fmt.Errorf("failed to save new collection: %w", err)
	}

	for _, collectionListener := range handler.collectionListeners {
		err = collectionListener.Subscribe(event.CollectionAddress)
		if err != nil {
			handler.logger.Error(
				"failed to notify collection listener about AgentCollectionsManager.CollectionCreated event",
				zap.Any("event", event),
				zap.Error(err),
			)
		}
	}

	return nil
}
