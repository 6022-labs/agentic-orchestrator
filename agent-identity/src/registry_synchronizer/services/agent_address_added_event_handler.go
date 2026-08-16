package services

import (
	"fmt"
	"maps"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"go.uber.org/zap"
)

type AgentAddressAddedEventHandler struct {
	*AbstractAgentAddressEventHandler
	logger                             *zap.Logger
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface
}

func NewAgentAddressAddedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentAddressRepository repositories.AgentAddressRepositoryInterface,
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface,
) *AgentAddressAddedEventHandler {
	return &AgentAddressAddedEventHandler{
		AbstractAgentAddressEventHandler: NewAbstractAgentAddressEventHandler(agentRepository, agentAddressRepository),

		logger:                             logger,
		agentCollectionBlockchainRequester: agentCollectionBlockchainRequester,
	}
}

func (handler *AgentAddressAddedEventHandler) Handle(event *abi.AgentCollectionV1AgentAddressAdded) error {
	handler.logger.Info("Handling AgentAddressAdded event",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("tokenId", event.TokenId.Uint64()),
		zap.String("addressType", event.AddressType),
		zap.String("value", event.Value),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		event.Raw.Address.Hex(),
		event.TokenId.Uint64(),
	)
	if err != nil {
		return fmt.Errorf("failed to get agent by collection address and token id: %w", err)
	}
	if agent == nil {
		return fmt.Errorf(
			"agent not found for collection %s and token id %d",
			event.Raw.Address.Hex(),
			event.TokenId.Uint64(),
		)
	}

	agentInformation, err := handler.agentCollectionBlockchainRequester.GetAgentInformation(
		event.Raw.Address,
		*event.TokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent information: %w", err)
	}

	addresses := handler.typedAddresses(agentInformation.Addresses)

	currentRows, err := handler.agentAddressRepository.GetAllByAgentId(agent.ID)
	if err != nil {
		return fmt.Errorf("failed to get agent addresses: %w", err)
	}
	current := map[enums.AddressType]string{}
	for _, row := range currentRows {
		current[row.Type] = row.Address
	}
	if maps.Equal(current, addresses) {
		return nil
	}

	if err := handler.agentAddressRepository.ReplaceAllForAgent(agent.ID, addresses); err != nil {
		return fmt.Errorf("failed to save agent addresses: %w", err)
	}

	return nil
}
