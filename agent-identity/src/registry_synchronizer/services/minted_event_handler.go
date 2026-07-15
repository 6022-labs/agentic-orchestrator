package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	registry_common_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"go.uber.org/zap"
)

type MintedEventHandler struct {
	*AbstractAgentAddressEventHandler
	logger                             *zap.Logger
	agentCollectionRepository          repositories.AgentCollectionRepositoryInterface
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface
	decentralizedDomainNameBuilder     registry_common_interfaces.DecentralizedDomainNameBuilderInterface
}

func NewMintedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentAddressRepository repositories.AgentAddressRepositoryInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface,
	decentralizedDomainNameBuilder registry_common_interfaces.DecentralizedDomainNameBuilderInterface,
) *MintedEventHandler {
	return &MintedEventHandler{
		AbstractAgentAddressEventHandler: NewAbstractAgentAddressEventHandler(agentRepository, agentAddressRepository),

		logger:                             logger,
		agentCollectionRepository:          agentCollectionRepository,
		decentralizedDomainNameBuilder:     decentralizedDomainNameBuilder,
		agentCollectionBlockchainRequester: agentCollectionBlockchainRequester,
	}
}

func (handler *MintedEventHandler) Handle(event *abi.AgentCollectionV1Minted) error {
	handler.logger.Info("Handling Minted event for collection",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("tokenId", event.TokenId.Uint64()),
		zap.String("to", event.To.Hex()),
	)

	agentInformation, err := handler.agentCollectionBlockchainRequester.GetAgentInformation(
		event.Raw.Address,
		*event.TokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent information: %w", err)
	}

	var agentRole enums.AgentRole
	err = agentRole.Set(agentInformation.Role)
	if err != nil {
		return fmt.Errorf("failed to set agent role: %w", err)
	}

	addresses := handler.typedAddresses(agentInformation.Addresses)
	agent, err := handler.findAgentByAddresses(addresses)
	if err != nil {
		return fmt.Errorf("failed to find agent by addresses: %w", err)
	}

	if agent == nil {
		agent = &models.Agent{}
	}

	owner, err := handler.agentCollectionBlockchainRequester.GetAgentOwner(
		event.Raw.Address,
		*event.TokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent owner: %w", err)
	}

	agentCollection, err := handler.agentCollectionRepository.GetByAddress(event.Raw.Address.Hex())
	if err != nil {
		return fmt.Errorf("failed to get agent collection by address: %w", err)
	}
	if agentCollection == nil {
		return fmt.Errorf("agent collection with address %s not found", event.Raw.Address.Hex())
	}

	ownerHex := owner.Hex()

	agent.Creator = agentInformation.Creator.Hex()
	agent.Role = agentRole
	agent.Owner = &ownerHex
	agent.Name = agentInformation.Name
	agent.EnsDomain = handler.decentralizedDomainNameBuilder.BuildDomainName(agentInformation.Name, agentCollection.Name)
	agent.ProposalId = nil

	if agentInformation.CloneOf.IsSet {
		cloneOf := agentInformation.CloneOf.Value.Uint64()
		agent.CloneOf = &cloneOf
	}

	agent.CollectionAddress = event.Raw.Address.Hex()

	collectionTokenIdUint64 := event.TokenId.Uint64()
	agent.CollectionTokenId = &collectionTokenIdUint64

	if agent.Images == nil {
		agent.Images = map[string]string{}
	}
	for _, image := range agentInformation.Images {
		agent.Images[image.Key] = image.Value
	}

	if agent.Attributes == nil {
		agent.Attributes = map[string]string{}
	}
	for _, attribute := range agentInformation.Attributes {
		agent.Attributes[attribute.Key] = attribute.Value
	}

	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return fmt.Errorf("failed to save agent: %w", err)
	}

	if err := handler.agentAddressRepository.ReplaceAllForAgent(agent.ID, addresses); err != nil {
		return fmt.Errorf("failed to save agent addresses: %w", err)
	}

	return nil
}
