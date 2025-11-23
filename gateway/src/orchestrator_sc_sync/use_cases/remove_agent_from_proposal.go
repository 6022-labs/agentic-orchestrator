package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"gorm.io/gorm"
)

type RemoveAgentFromProposalInterface interface {
	Remove(request *requests.RemoveAgentFromProposalRequest) error
}

type RemoveAgentFromProposal struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewRemoveAgentFromProposal(agentRepository repositories.AgentRepositoryInterface) *RemoveAgentFromProposal {
	return &RemoveAgentFromProposal{
		agentRepository: agentRepository,
	}
}

func (r *RemoveAgentFromProposal) Remove(request *requests.RemoveAgentFromProposalRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agent, err := r.agentRepository.GetNotYetMintedByCollectionAddressAndProposalId(request.CollectionAddress, request.ProposalId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if agent == nil {
		return fmt.Errorf("no agent found for collection address %s and proposal id %d", request.CollectionAddress, request.ProposalId)
	}

	err = r.agentRepository.Remove(agent)
	if err != nil {
		return fmt.Errorf("failed to remove agent from proposal: %w", err)
	}

	return nil
}
