package services

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"go.uber.org/zap"
)

type AgentImageRemovedEventHandler struct {
	logger          *zap.Logger
	agentRepository repositories.AgentRepositoryInterface
}

func NewAgentImageRemovedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
) *AgentImageRemovedEventHandler {
	return &AgentImageRemovedEventHandler{
		logger:          logger,
		agentRepository: agentRepository,
	}
}

func (handler *AgentImageRemovedEventHandler) Handle(
	ctx context.Context,
	event *abi.AgentCollectionV1AgentImageRemoved,
) error {
	handler.logger.Info("Handling AgentImageRemoved event",
		zap.String("agentCollection", event.Raw.Address.Hex()),
		zap.Uint64("agentId", event.TokenId.Uint64()),
		zap.String("key", event.Key.Hex()),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		event.Raw.Address.Hex(),
		event.TokenId.Uint64(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get agent by collection %s and id %s: %w",
			event.Raw.Address.Hex(),
			event.TokenId.String(),
			err,
		)
	}
	if agent == nil {
		return fmt.Errorf(
			"agent not found for collection %s and id %s",
			event.Raw.Address.Hex(),
			event.TokenId.String(),
		)
	}

	if agent.Images == nil {
		return nil
	}

	delete(agent.Images, event.Key.String())

	_, err = handler.agentRepository.Save(ctx, agent)
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	return nil
}
