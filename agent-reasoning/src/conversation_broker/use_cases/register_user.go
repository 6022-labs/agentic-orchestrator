package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
)

type RegisterUserInterface interface {
	Register(request *requests.RegisterUserRequest) (*models.Participant, error)
}

// RegisterUser always creates a new participant. There is no natural key for
// users (display names collide), so dedup must happen on the bridge side via
// a local cache (e.g. slack_users.broker_participant_id).
type RegisterUser struct {
	participantRepository repositories.ParticipantRepositoryInterface
}

func NewRegisterUser(participantRepository repositories.ParticipantRepositoryInterface) *RegisterUser {
	return &RegisterUser{participantRepository: participantRepository}
}

func (uc *RegisterUser) Register(request *requests.RegisterUserRequest) (*models.Participant, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	displayName := request.DisplayName
	participant := &models.Participant{
		DisplayName: &displayName,
	}

	saved, err := uc.participantRepository.Save(participant)
	if err != nil {
		return nil, fmt.Errorf("failed to save user participant: %w", err)
	}

	return saved, nil
}
