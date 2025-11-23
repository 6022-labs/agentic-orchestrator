package use_cases

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
)

type UpdateAgentOwnerInterface interface {
	Update(request *requests.UpdateAgentOwnerRequest) error
}

type UpdateAgentOwner struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewUpdateAgentOwner(agentRepository repositories.AgentRepositoryInterface) *UpdateAgentOwner {
	return &UpdateAgentOwner{
		agentRepository: agentRepository,
	}
}

func (u *UpdateAgentOwner) Update(request *requests.UpdateAgentOwnerRequest) error {
	agent, err := u.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.CollectionAddress,
		request.CollectionTokenId,
	)
	if err != nil {
		return err
	}

	agent.Owner = &request.NewOwner
	_, err = u.agentRepository.Save(agent)
	if err != nil {
		return err
	}

	return nil
}
