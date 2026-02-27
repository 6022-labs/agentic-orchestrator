package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
)

type UpdateAgentAttributeInterface interface {
	Update(request *requests.UpdateAgentAttributeRequest) error
}

type UpdateAgentAttribute struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewUpdateAgentAttribute(agentRepository repositories.AgentRepositoryInterface) *UpdateAgentAttribute {
	return &UpdateAgentAttribute{
		agentRepository: agentRepository,
	}
}

func (u *UpdateAgentAttribute) Update(request *requests.UpdateAgentAttributeRequest) error {
	agent, err := u.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.CollectionAddress,
		request.CollectionTokenId,
	)
	if err != nil {
		return err
	}

	if agent.Attributes == nil {
		agent.Attributes = make(map[string]string)
	}

	agent.Attributes[request.Key] = request.NewValue

	_, err = u.agentRepository.Save(agent)
	if err != nil {
		return err
	}

	return nil
}
