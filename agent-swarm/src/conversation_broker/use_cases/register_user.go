package use_cases

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type RegisterUserInterface interface {
	Register(ctx context.Context, request *requests.RegisterUserRequest) (*responses.ParticipantResponse, error)
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

func (uc *RegisterUser) Register(
	ctx context.Context,
	request *requests.RegisterUserRequest,
) (*responses.ParticipantResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	displayName := request.DisplayName
	participant := &models.Participant{
		DisplayName: &displayName,
	}

	saved, err := uc.participantRepository.Save(ctx, participant)
	if err != nil {
		return nil, fmt.Errorf("failed to save user participant: %w", err)
	}

	return responses.ParticipantResponseFromModel(saved), nil
}
