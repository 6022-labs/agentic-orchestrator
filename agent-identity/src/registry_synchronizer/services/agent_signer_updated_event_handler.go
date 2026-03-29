package services

import (
	"errors"
	"log"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"gorm.io/gorm"
)

type AgentSignerUpdatedEventHandler struct {
	logger *log.Logger

	orchestratorSettings settings.OrchestratorSettings

	agentRepository          repositories.AgentRepositoryInterface
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface
}

func NewAgentSignerUpdatedEventHandler(
	logger *log.Logger,
	orchestratorSettings settings.OrchestratorSettings,
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
	handler.logger.Printf("Handling AgentSignerUpdated event for collection %s, agentId %d, new signer %s",
		event.Raw.Address.Hex(),
		event.TokenId.Uint64(),
		event.NewSigner.Hex(),
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
