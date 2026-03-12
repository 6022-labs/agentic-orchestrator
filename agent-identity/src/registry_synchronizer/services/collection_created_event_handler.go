package services

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"go.uber.org/zap"
)

type CollectionCreatedEventHandler struct {
	logger                    *zap.Logger
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface
}

func NewCollectionCreatedEventHandler(
	logger *zap.Logger,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
) *CollectionCreatedEventHandler {
	return &CollectionCreatedEventHandler{
		logger:                    logger,
		agentCollectionRepository: agentCollectionRepository,
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
		return errors.New("failed to save new collection: " + err.Error())
	}

	return nil
}
