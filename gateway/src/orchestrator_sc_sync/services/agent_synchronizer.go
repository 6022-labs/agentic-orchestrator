package services

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/settings"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type AgentSynchronizerInterface interface {
	Sync(agentCollectionAddress string, agentCollectionTokenId uint64) error
}

type AgentSynchronizer struct {
	orchestratorSettings           *settings.OrchestratorSettings
	agentRepository                repositories.AgentRepositoryInterface
	agentCollectionRequester       interfaces.AgentCollectionRequesterInterface
	agentCollectionRepository      repositories.AgentCollectionRepositoryInterface
	decentralizedDomainNameBuilder interfaces.DecentralizedDomainNameBuilderInterface
}

func NewAgentSynchronizer(
	orchestratorSettings *settings.OrchestratorSettings,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	decentralizedDomainNameBuilder interfaces.DecentralizedDomainNameBuilderInterface,
) *AgentSynchronizer {
	return &AgentSynchronizer{
		agentRepository:                agentRepository,
		orchestratorSettings:           orchestratorSettings,
		agentCollectionRequester:       agentCollectionRequester,
		agentCollectionRepository:      agentCollectionRepository,
		decentralizedDomainNameBuilder: decentralizedDomainNameBuilder,
	}
}

func (s *AgentSynchronizer) Sync(agentCollectionAddress string, agentCollectionTokenId uint64) error {
	typedAgentCollectionTokenId := new(big.Int).SetUint64(agentCollectionTokenId)
	typedAgentCollectionAddress := common.HexToAddress(agentCollectionAddress)

	agentInformation, err := s.agentCollectionRequester.GetAgentInformation(
		typedAgentCollectionAddress,
		*typedAgentCollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent information: %w", err)
	}

	var agentRole enums.AgentRole
	err = agentRole.Set(agentInformation.Role)
	if err != nil {
		return fmt.Errorf("failed to set agent role: %w", err)
	}

	agent, err := s.findAgentByWalletAddresses(agentInformation.Wallets)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to find agent by wallet addresses: %w", err)
	}

	if agent == nil {
		if s.orchestratorSettings.PublicKey.Hex() == agentInformation.Signer.Hex() {
			return fmt.Errorf("agent with address %s must be created by this orchestrator but it was not found", agentInformation.Wallets[0].Hex())
		}

		// Create a new external agent
		agent = &models.Agent{
			IsInternal: false,
		}
	}

	owner, err := s.agentCollectionRequester.GetAgentOwner(
		typedAgentCollectionAddress,
		*typedAgentCollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent owner: %w", err)
	}

	agentCollection, err := s.agentCollectionRepository.GetByAddress(agentCollectionAddress)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get agent collection by address: %w", err)
	}
	if agentCollection == nil {
		return fmt.Errorf("agent collection with address %s not found", agentCollectionAddress)
	}

	ownerHex := owner.Hex()

	agent.WalletAddress = agentInformation.Wallets[0].Hex()
	agent.Creator = agentInformation.Creator.Hex()
	agent.Role = agentRole
	agent.Owner = &ownerHex
	agent.Name = agentInformation.Name
	agent.EnsDomain = s.decentralizedDomainNameBuilder.BuildDomainName(agentInformation.Name, agentCollection.Name)
	agent.ProposalId = nil

	if agentInformation.CloneOf.IsSet {
		cloneOf := agentInformation.CloneOf.Value.Uint64()
		agent.CloneOf = &cloneOf
	}

	agent.CollectionAddress = agentCollectionAddress

	collectionTokenIdUint64 := agentCollectionTokenId
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

	_, err = s.agentRepository.Save(agent)
	if err != nil {
		return fmt.Errorf("failed to save agent: %w", err)
	}

	return nil
}

func (f *AgentSynchronizer) findAgentByWalletAddresses(walletAddress []common.Address) (*models.Agent, error) {
	for _, wallet := range walletAddress {
		agent, err := f.agentRepository.GetByWalletAddress(wallet.Hex())
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to get agent by wallet address: %w", err)
		}

		if agent != nil {
			return agent, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}
