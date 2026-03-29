package services

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AgentSignerUpdatedEventHandler struct {
	logger *zap.Logger

	orchestratorSettings *settings.OrchestratorSettings

	agentRepository          repositories.AgentRepositoryInterface
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface
}

func NewAgentSignerUpdatedEventHandler(
	logger *zap.Logger,
	orchestratorSettings *settings.OrchestratorSettings,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface,
) *AgentSignerUpdatedEventHandler {
	return &AgentSignerUpdatedEventHandler{
		logger: logger,

		orchestratorSettings: orchestratorSettings,

		agentRepository:          agentRepository,
		agentCollectionRequester: agentCollectionRequester,
	}
}

func (handler *AgentSignerUpdatedEventHandler) Handle(event *abi.AgentCollectionV1AgentSignerUpdated) error {
	handler.logger.Info("Handling AgentSignerUpdated event",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("agentId", event.TokenId.Uint64()),
		zap.String("newSigner", event.NewSigner.Hex()),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(event.Raw.Address.Hex(), event.TokenId.Uint64())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to get agent by collection and id: " + err.Error())
	}
	if agent == nil {
		return errors.New("agent not found for collection " + event.Raw.Address.Hex() + " and id " + event.TokenId.String())
	}

	isNewSignerInternal := event.NewSigner.Cmp(handler.orchestratorSettings.PublicKey) == 0
	agent.IsInternal = isNewSignerInternal

	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return errors.New("failed to update agent: " + err.Error())
	}

	return nil
}
