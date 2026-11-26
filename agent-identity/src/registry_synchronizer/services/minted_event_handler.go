package services

import (
	"context"

	"fmt"
	"strings"

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

func (handler *MintedEventHandler) Handle(ctx context.Context, event *abi.AgentCollectionV1Minted) error {
	handler.logger.Info("Handling Minted event for collection",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("tokenId", event.TokenId.Uint64()),
		zap.String("to", event.To.Hex()),
	)

	agentInformation, err := handler.agentCollectionBlockchainRequester.GetAgentInformation(
		ctx,
		event.Raw.Address,
		*event.TokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent information: %w", err)
	}

	collectionAddress := strings.ToLower(event.Raw.Address.Hex())
	addresses := handler.typedAddresses(agentInformation.Addresses)
	agent, err := handler.findAgentToMint(
		ctx,
		collectionAddress,
		event.TokenId.Uint64(),
		agentInformation.Name,
		addresses,
	)
	if err != nil {
		return err
	}

	if agent == nil {
		agent = &models.Agent{}
	}

	owner, err := handler.agentCollectionBlockchainRequester.GetAgentOwner(
		ctx,
		event.Raw.Address,
		*event.TokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent owner: %w", err)
	}

	agentCollection, err := handler.agentCollectionRepository.GetByAddress(ctx, collectionAddress)
	if err != nil {
		return fmt.Errorf("failed to get agent collection by address: %w", err)
	}
	if agentCollection == nil {
		return fmt.Errorf("agent collection with address %s not found", event.Raw.Address.Hex())
	}

	ownerHex := strings.ToLower(owner.Hex())

	agent.Creator = strings.ToLower(agentInformation.Creator.Hex())
	agent.Owner = &ownerHex
	agent.Name = agentInformation.Name
	agent.EnsDomain = handler.decentralizedDomainNameBuilder.BuildDomainName(
		agentInformation.Name,
		agentCollection.Name,
	)
	agent.ProposalId = nil

	if agentInformation.CloneOf.IsSet {
		cloneOf := agentInformation.CloneOf.Value.Uint64()
		agent.CloneOf = &cloneOf
	}

	agent.CollectionAddress = collectionAddress

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

	_, err = handler.agentRepository.Save(ctx, agent)
	if err != nil {
		return fmt.Errorf("failed to save agent: %w", err)
	}

	if err := handler.agentAddressRepository.ReplaceAllForAgent(ctx, agent.ID, addresses); err != nil {
		return fmt.Errorf("failed to save agent addresses: %w", err)
	}

	return nil
}

// findAgentToMint returns the row this mint fills: the token's own row on replay, else the proposal it executes.
func (handler *MintedEventHandler) findAgentToMint(
	ctx context.Context,
	collectionAddress string,
	tokenId uint64,
	name string,
	addresses map[enums.AddressType]string,
) (*models.Agent, error) {
	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(ctx, collectionAddress, tokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent by collection and token id: %w", err)
	}
	if agent != nil {
		return agent, nil
	}

	agent, err = handler.agentRepository.GetNotYetMintedByCollectionAddressNameAndAddresses(
		ctx,
		collectionAddress,
		name,
		addresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal by name and addresses: %w", err)
	}

	return agent, nil
}
