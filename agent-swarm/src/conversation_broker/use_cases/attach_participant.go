package use_cases

import (
	"context"

	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
)

type AttachParticipantInterface interface {
	Attach(ctx context.Context, request *requests.ConversationParticipantRequest) error
}

type AttachParticipant struct {
	conversationRepository            repositories.ConversationRepositoryInterface
	participantRepository             repositories.ParticipantRepositoryInterface
	conversationParticipantRepository repositories.ConversationParticipantRepositoryInterface
}

func NewAttachParticipant(
	conversationRepository repositories.ConversationRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
	conversationParticipantRepository repositories.ConversationParticipantRepositoryInterface,
) *AttachParticipant {
	return &AttachParticipant{
		conversationRepository:            conversationRepository,
		participantRepository:             participantRepository,
		conversationParticipantRepository: conversationParticipantRepository,
	}
}

func (uc *AttachParticipant) Attach(ctx context.Context, request *requests.ConversationParticipantRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	conversation, err := uc.conversationRepository.GetById(ctx, request.ConversationID)
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}
	if conversation == nil {
		return cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
	}

	participant, err := uc.participantRepository.GetById(ctx, request.ParticipantID)
	if err != nil {
		return fmt.Errorf("failed to load participant: %w", err)
	}
	if participant == nil {
		return cb_errors.NewNotFoundError("participant", request.ParticipantID.String())
	}

	link := &models.ConversationParticipant{
		ConversationID: request.ConversationID,
		ParticipantID:  request.ParticipantID,
	}
	if _, err := uc.conversationParticipantRepository.Save(ctx, link); err != nil {
		return fmt.Errorf("failed to attach participant: %w", err)
	}
	return nil
}
