package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"go.uber.org/zap"
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
	if err != nil {
		return fmt.Errorf("failed to get agent by collection %s and id %s: %w", event.Raw.Address.Hex(), event.TokenId.String(), err)
	}
	if agent == nil {
		return fmt.Errorf("agent not found for collection %s and id %s", event.Raw.Address.Hex(), event.TokenId.String())
	}

	if agent.Images == nil {
		agent.Images = make(map[string]string)
	}

	agent.Images[event.Key.String()] = event.NewValue

	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	return nil
}
