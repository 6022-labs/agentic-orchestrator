package services

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AgentImageUpdatedEventHandler struct {
	logger          *zap.Logger
	agentRepository repositories.AgentRepositoryInterface
}

func NewAgentImageUpdatedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
) *AgentImageUpdatedEventHandler {
	return &AgentImageUpdatedEventHandler{
		logger:          logger,
		agentRepository: agentRepository,
	}
}

func (handler *AgentImageUpdatedEventHandler) Handle(event *abi.AgentCollectionV1AgentImageUpdated) error {
	handler.logger.Info("Handling AgentImageUpdated event",
		zap.String("agentCollection", event.Raw.Address.Hex()),
		zap.Uint64("agentId", event.TokenId.Uint64()),
		zap.String("key", event.Key.Hex()),
		zap.String("oldValue", event.OldValue),
		zap.String("newValue", event.NewValue),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(event.Raw.Address.Hex(), event.TokenId.Uint64())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to get agent by collection and id: " + err.Error())
	}
	if agent == nil {
		return errors.New("agent not found for collection " + event.Raw.Address.Hex() + " and id " + event.TokenId.String())
	}

	if agent.Images == nil {
		agent.Images = make(map[string]string)
	}

	agent.Images[event.Key.String()] = event.NewValue

	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return errors.New("failed to update agent: " + err.Error())
	}

	return nil
}
