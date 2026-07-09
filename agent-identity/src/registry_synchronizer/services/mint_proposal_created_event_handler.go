package services

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	registry_common_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MintProposalCreatedEventHandler struct {
	logger                         *zap.Logger
	agentRepository                repositories.AgentRepositoryInterface
	agentCollectionRepository      repositories.AgentCollectionRepositoryInterface
	agentCollectionRequester       interfaces.AgentCollectionBlockchainRequesterInterface
	decentralizedDomainNameBuilder registry_common_interfaces.DecentralizedDomainNameBuilderInterface
}

func NewMintProposalCreatedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface,
	decentralizedDomainNameBuilder registry_common_interfaces.DecentralizedDomainNameBuilderInterface,
) *MintProposalCreatedEventHandler {
	return &MintProposalCreatedEventHandler{
		logger:                         logger,
		agentRepository:                agentRepository,
		agentCollectionRepository:      agentCollectionRepository,
		agentCollectionRequester:       agentCollectionRequester,
		decentralizedDomainNameBuilder: decentralizedDomainNameBuilder,
	}
}

func (handler *MintProposalCreatedEventHandler) Handle(event *abi.AgentCollectionV1MintProposalCreated) error {
	collectionAddress := common.HexToAddress(event.Raw.Address.Hex())
	proposalId := new(big.Int).SetUint64(event.ProposalId.Uint64())

	atBlock := event.Raw.BlockNumber
	mintProposal, err := handler.agentCollectionRequester.GetMintProposal(collectionAddress, proposalId, atBlock)
	if err != nil {
		return fmt.Errorf("failed to get mint proposal for collection address %s and proposal id %d: %w",
			collectionAddress.Hex(),
			proposalId.Uint64(),
			err,
		)
	}

	if len(mintProposal.Wallets) == 0 {
		return fmt.Errorf("no agent wallets found for collection address %s and proposal id %d",
			collectionAddress.Hex(),
			proposalId.Uint64(),
		)
	}

	agentWallet := mintProposal.Wallets[0]

	agent, err := handler.agentRepository.GetNotYetMintedByCollectionAddressAndWalletAddress(
		collectionAddress.Hex(),
		agentWallet.Hex(),
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get agent by collection address %s and wallet address %s: %w",
			collectionAddress.Hex(),
			agentWallet.Hex(),
			err,
		)
	}

	if agent == nil {
		existingAgent, err := handler.agentRepository.GetByWalletAddress(agentWallet.Hex())
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to get agent by wallet address %s: %w",
				agentWallet.Hex(),
				err,
			)
		}
		if existingAgent != nil {
			return nil
		}

		var agentRole enums.AgentRole
		if err := agentRole.Set(mintProposal.Role); err != nil {
			return fmt.Errorf("failed to set agent role: %w", err)
		}

		agentCollection, err := handler.agentCollectionRepository.GetByAddress(collectionAddress.Hex())
		if err != nil {
			return fmt.Errorf("failed to get agent collection by address %s: %w", collectionAddress.Hex(), err)
		}

		ownerHex := mintProposal.Owner.Hex()

		agent = &models.Agent{
			WalletAddress:     agentWallet.Hex(),
			CollectionAddress: collectionAddress.Hex(),
			Name:              mintProposal.Name,
			Role:              agentRole,
			Creator:           mintProposal.Creator.Hex(),
			Owner:             &ownerHex,
			EnsDomain:         handler.decentralizedDomainNameBuilder.BuildDomainName(mintProposal.Name, agentCollection.Name),
			Images:            map[string]string{},
			Attributes:        map[string]string{},
		}

		if mintProposal.CloneOf.IsSet {
			cloneOf := mintProposal.CloneOf.Value.Uint64()
			agent.CloneOf = &cloneOf
		}

		for _, image := range mintProposal.Images {
			agent.Images[image.Key] = image.Value
		}

		for _, attribute := range mintProposal.Attributes {
			agent.Attributes[attribute.Key] = attribute.Value
		}
	}

	agent.ProposalId = lo.ToPtr(proposalId.Uint64())
	_, err = handler.agentRepository.Save(agent)
	if err != nil {
		return fmt.Errorf("failed to save agent with collection address %s and wallet address %s: %w",
			collectionAddress.Hex(),
			agentWallet.Hex(),
			err,
		)
	}

	return nil
}
