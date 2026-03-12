package services

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AgentAttributeUpdatedEventHandler struct {
	logger          *zap.Logger
	agentRepository repositories.AgentRepositoryInterface
}

func NewAgentAttributeUpdatedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
) *AgentAttributeUpdatedEventHandler {
	return &AgentAttributeUpdatedEventHandler{
		logger:          logger,
		agentRepository: agentRepository,
	}
}

func (handler *AgentAttributeUpdatedEventHandler) Handle(event *abi.AgentCollectionV1AgentAttributeUpdated) error {
	handler.logger.Info("Handling AgentAttributeUpdated event",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("tokenId", event.TokenId.Uint64()),
		zap.String("key", event.Key.String()),
		zap.String("newValue", event.NewValue),
		zap.String("oldValue", event.OldValue),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(event.Raw.Address.Hex(), event.TokenId.Uint64())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to get agent by collection and id: " + err.Error())
	}
	if agent == nil {
		return errors.New("agent not found for collection " + event.Raw.Address.Hex() + " and id " + event.TokenId.String())
	}

	if agent.Attributes == nil {
		agent.Attributes = make(map[string]string)
	}

	agent.Attributes[event.Key.String()] = event.NewValue

	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return errors.New("failed to update agent: " + err.Error())
	}

	return nil
}
