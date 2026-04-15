package use_cases

import (
	"errors"
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"gorm.io/gorm"
)

type AttachParticipantInterface interface {
	Attach(request *requests.ConversationParticipantRequest) error
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

func (uc *AttachParticipant) Attach(request *requests.ConversationParticipantRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	_, err := uc.conversationRepository.GetById(request.ConversationID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
	}
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}

	_, err = uc.participantRepository.GetById(request.ParticipantID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return cb_errors.NewNotFoundError("participant", request.ParticipantID.String())
	}
	if err != nil {
		return fmt.Errorf("failed to load participant: %w", err)
	}

	link := &models.ConversationParticipant{
		ConversationID: request.ConversationID,
		ParticipantID:  request.ParticipantID,
	}
	if _, err := uc.conversationParticipantRepository.Save(link); err != nil {
		return fmt.Errorf("failed to attach participant: %w", err)
	}
	return nil
}
