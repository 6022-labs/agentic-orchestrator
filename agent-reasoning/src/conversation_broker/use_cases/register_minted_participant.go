package use_cases

import (
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
)

type RegisterMintedParticipantInterface interface {
	Register(request *requests.RegisterMintedParticipantRequest) (*models.Participant, error)
}

type RegisterMintedParticipant struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
	participantRepository  repositories.ParticipantRepositoryInterface
}

func NewRegisterMintedParticipant(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
) *RegisterMintedParticipant {
	return &RegisterMintedParticipant{
		agentRegistryRequester: agentRegistryRequester,
		participantRepository:  participantRepository,
	}
}

func (uc *RegisterMintedParticipant) Register(request *requests.RegisterMintedParticipantRequest) (*models.Participant, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	participant, err := uc.agentRegistryRequester.GetMintedAgent(request.CollectionAddress, request.ParsedCollectionTokenId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent from registry: %w", err)
	}
	if participant == nil {
		return nil, cb_errors.NewNotFoundError("agent", request.CollectionAddress+"/"+request.CollectionTokenId)
	}

	saved, err := uc.participantRepository.Save(participant)
	if err != nil {
		return nil, fmt.Errorf("failed to save minted participant: %w", err)
	}
	return saved, nil
}
