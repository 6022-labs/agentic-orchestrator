package use_cases

import (
	"fmt"
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/ethereum/go-ethereum/common"
)

type PairProposalToAgentInterface interface {
	Pair(request *requests.PairProposalToAgentRequest) error
}

type PairProposalToAgent struct {
	agentRepository          repositories.AgentRepositoryInterface
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface
}

func NewPairProposalToAgent(
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionRequester interfaces.AgentCollectionRequesterInterface,
) *PairProposalToAgent {
	return &PairProposalToAgent{
		agentRepository:          agentRepository,
		agentCollectionRequester: agentCollectionRequester,
	}
}

func (p *PairProposalToAgent) Pair(request *requests.PairProposalToAgentRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	collectionAddress := common.HexToAddress(request.CollectionAddress)
	proposalId := new(big.Int).SetUint64(request.ProposalId)

	agentWallets, err := p.agentCollectionRequester.GetAgentWalletsOfProposal(collectionAddress, proposalId)
	if err != nil {
		return fmt.Errorf("failed to get agent wallets for collection address %s and proposal id %d: %w",
			request.CollectionAddress,
			request.ProposalId,
			err,
		)
	}
	if len(agentWallets) == 0 {
		return fmt.Errorf("no agent wallets found for collection address %s and proposal id %d",
			request.CollectionAddress,
			request.ProposalId,
		)
	}

	agentWallet := agentWallets[0]

	agent, err := p.agentRepository.GetNotYetMintedByCollectionAddressAndWalletAddress(
		request.CollectionAddress,
		agentWallet.Hex(),
	)
	if err != nil {
		return fmt.Errorf("failed to get agent by collection address %s and wallet address %s: %w",
			request.CollectionAddress,
			agentWallet.Hex(),
			err,
		)
	}
	if agent == nil {
		return fmt.Errorf("no agent found for collection address %s and wallet address %s",
			request.CollectionAddress,
			agentWallet.Hex(),
		)
	}

	agent.ProposalId = &request.ProposalId
	_, err = p.agentRepository.Save(agent)
	if err != nil {
		return err
	}

	return nil
}
