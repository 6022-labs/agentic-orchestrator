package services

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	registry_common_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type MintProposalCreatedEventHandler struct {
	*AbstractAgentAddressEventHandler
	logger                         *zap.Logger
	agentCollectionRepository      repositories.AgentCollectionRepositoryInterface
	agentCollectionRequester       interfaces.AgentCollectionBlockchainRequesterInterface
	decentralizedDomainNameBuilder registry_common_interfaces.DecentralizedDomainNameBuilderInterface
}

func NewMintProposalCreatedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentAddressRepository repositories.AgentAddressRepositoryInterface,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface,
	decentralizedDomainNameBuilder registry_common_interfaces.DecentralizedDomainNameBuilderInterface,
) *MintProposalCreatedEventHandler {
	return &MintProposalCreatedEventHandler{
		AbstractAgentAddressEventHandler: NewAbstractAgentAddressEventHandler(agentRepository, agentAddressRepository),

		logger:                         logger,
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

	addresses := handler.typedAddresses(mintProposal.Addresses)
	if len(addresses) == 0 {
		return fmt.Errorf("no known agent address found for collection address %s and proposal id %d",
			collectionAddress.Hex(),
			proposalId.Uint64(),
		)
	}

	agent, err := handler.findAgentByAddresses(addresses)
	if err != nil {
		return fmt.Errorf("failed to get agent by address: %w", err)
	}

	// An agent already minted (or drafted in another collection) holds one of the addresses; nothing to sync.
	if agent != nil && (agent.CollectionAddress != collectionAddress.Hex() || agent.CollectionTokenId != nil) {
		return nil
	}

	if agent == nil {
		agentCollection, err := handler.agentCollectionRepository.GetByAddress(collectionAddress.Hex())
		if err != nil {
			return fmt.Errorf("failed to get agent collection by address %s: %w", collectionAddress.Hex(), err)
		}

		ownerHex := mintProposal.Owner.Hex()

		agent = &models.Agent{
			CollectionAddress: collectionAddress.Hex(),
			Name:              mintProposal.Name,
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
		return fmt.Errorf("failed to save agent with collection address %s and proposal id %d: %w",
			collectionAddress.Hex(),
			proposalId.Uint64(),
			err,
		)
	}

	if err := handler.agentAddressRepository.ReplaceAllForAgent(agent.ID, addresses); err != nil {
		return fmt.Errorf("failed to save agent addresses: %w", err)
	}

	return nil
}
