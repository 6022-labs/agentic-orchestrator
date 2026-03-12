package services

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MintedEventHandler struct {
	logger                             *zap.Logger
	orchestratorSettings               *settings.OrchestratorSettings
	agentRepository                    repositories.AgentRepositoryInterface
	agentCollectionRepository          repositories.AgentCollectionRepositoryInterface
	decentralizedDomainNameBuilder     interfaces.DecentralizedDomainNameBuilderInterface
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface
}

func NewMintedEventHandler(
	logger *zap.Logger,
	orchestratorSettings *settings.OrchestratorSettings,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	decentralizedDomainNameBuilder interfaces.DecentralizedDomainNameBuilderInterface,
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface,
) *MintedEventHandler {
	return &MintedEventHandler{
		logger:               logger,
		orchestratorSettings: orchestratorSettings,

		agentRepository:                    agentRepository,
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

	agent, err := handler.findAgentByWalletAddresses(agentInformation.Wallets)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to find agent by wallet addresses: %w", err)
	}

	if agent == nil {
		if handler.orchestratorSettings.PublicKey.Hex() == agentInformation.Signer.Hex() {
			return fmt.Errorf("agent with address %s must be created by this orchestrator but it was not found", agentInformation.Wallets[0].Hex())
		}

		// Create a new external agent
		agent = &models.Agent{
			IsInternal: false,
		}
	}

	owner, err := handler.agentCollectionBlockchainRequester.GetAgentOwner(
		event.Raw.Address,
		*event.TokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent owner: %w", err)
	}

	agentCollection, err := handler.agentCollectionRepository.GetByAddress(event.Raw.Address.Hex())
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get agent collection by address: %w", err)
	}
	if agentCollection == nil {
		return fmt.Errorf("agent collection with address %s not found", event.Raw.Address.Hex())
	}

	ownerHex := owner.Hex()

	agent.WalletAddress = agentInformation.Wallets[0].Hex()
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

	return nil
}

func (handler *MintedEventHandler) findAgentByWalletAddresses(walletAddress []common.Address) (*models.Agent, error) {
	for _, wallet := range walletAddress {
		agent, err := handler.agentRepository.GetByWalletAddress(wallet.Hex())
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to get agent by wallet address: %w", err)
		}

		if agent != nil {
			return agent, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}
