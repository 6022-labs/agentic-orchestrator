package services

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type MintProposalCreatedEventHandler struct {
	logger                   *zap.Logger
	agentRepository          repositories.AgentRepositoryInterface
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface
}

func NewMintProposalCreatedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionBlockchainRequesterInterface,
) *MintProposalCreatedEventHandler {
	return &MintProposalCreatedEventHandler{
		logger:                   logger,
		agentRepository:          agentRepository,
		agentCollectionRequester: agentCollectionRequester,
	}
}

func (handler *MintProposalCreatedEventHandler) Handle(event *abi.AgentCollectionV1MintProposalCreated) error {
	collectionAddress := common.HexToAddress(event.Raw.Address.Hex())
	proposalId := new(big.Int).SetUint64(event.ProposalId.Uint64())

	agentWallets, err := handler.agentCollectionRequester.GetAgentWalletsOfProposal(collectionAddress, proposalId)
	if err != nil {
		return fmt.Errorf("failed to get agent wallets for collection address %s and proposal id %d: %w",
			collectionAddress.Hex(),
			proposalId.Uint64(),
			err,
		)
	}
	if len(agentWallets) == 0 {
		return fmt.Errorf("no agent wallets found for collection address %s and proposal id %d",
			collectionAddress.Hex(),
			proposalId.Uint64(),
		)
	}

	agentWallet := agentWallets[0]

	agent, err := handler.agentRepository.GetNotYetMintedByCollectionAddressAndWalletAddress(
		collectionAddress.Hex(),
		agentWallet.Hex(),
	)
	if err != nil {
		return fmt.Errorf("failed to get agent by collection address %s and wallet address %s: %w",
			collectionAddress.Hex(),
			agentWallet.Hex(),
			err,
		)
	}
	if agent == nil {
		return fmt.Errorf("no agent found for collection address %s and wallet address %s",
			collectionAddress.Hex(),
			agentWallet.Hex(),
		)
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
