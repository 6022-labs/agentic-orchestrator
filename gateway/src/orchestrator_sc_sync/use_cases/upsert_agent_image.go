package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
)

type UpsertAgentImageInterface interface {
	Upsert(request *requests.UpsertAgentImageRequest) error
}

type UpsertAgentImage struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewUpsertAgentImage(agentRepository repositories.AgentRepositoryInterface) *UpsertAgentImage {
	return &UpsertAgentImage{
		agentRepository: agentRepository,
	}
}

func (u *UpsertAgentImage) Upsert(request *requests.UpsertAgentImageRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agent, err := u.agentRepository.GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, request.CollectionTokenId)
	if err != nil {
		return err
	}
	if agent == nil {
		return fmt.Errorf("no agent found for collection address %s and collection token id %d",
			request.CollectionAddress,
			request.CollectionTokenId,
		)
	}

	if agent.Images == nil {
		agent.Images = make(map[string]string)
	}

	agent.Images[request.Key] = request.NewValue

	_, err = u.agentRepository.Save(agent)
	if err != nil {
		return fmt.Errorf("failed to upsert agent image: %w", err)
	}

	return nil
}
