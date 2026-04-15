package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
)

type RegisterNonMintedParticipantInterface interface {
	Register(request *requests.RegisterNonMintedParticipantRequest) (*models.Participant, error)
}

type RegisterNonMintedParticipant struct {
	participantRepository repositories.ParticipantRepositoryInterface
}

func NewRegisterNonMintedParticipant(participantRepository repositories.ParticipantRepositoryInterface) *RegisterNonMintedParticipant {
	return &RegisterNonMintedParticipant{participantRepository: participantRepository}
}

func (uc *RegisterNonMintedParticipant) Register(request *requests.RegisterNonMintedParticipantRequest) (*models.Participant, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	displayName := request.DisplayName
	participant := &models.Participant{
		DisplayName: &displayName,
	}

	saved, err := uc.participantRepository.Save(participant)
	if err != nil {
		return nil, fmt.Errorf("failed to save non-minted participant: %w", err)
	}

	return saved, nil
}
