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

type IngestMessageInterface interface {
	Ingest(request *requests.IngestMessageRequest) (*models.ConversationMessage, error)
}

type IngestMessage struct {
	conversationRepository        repositories.ConversationRepositoryInterface
	participantRepository         repositories.ParticipantRepositoryInterface
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface
}

func NewIngestMessage(
	conversationRepository repositories.ConversationRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface,
) *IngestMessage {
	return &IngestMessage{
		conversationRepository:        conversationRepository,
		participantRepository:         participantRepository,
		conversationMessageRepository: conversationMessageRepository,
	}
}

func (uc *IngestMessage) Ingest(request *requests.IngestMessageRequest) (*models.ConversationMessage, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	_, err := uc.conversationRepository.GetById(request.ConversationID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}

	_, err = uc.participantRepository.GetById(request.ParticipantID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, cb_errors.NewNotFoundError("participant", request.ParticipantID.String())
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load participant: %w", err)
	}

	message := &models.ConversationMessage{
		ConversationID: request.ConversationID,
		ParticipantID:  request.ParticipantID,
		Content:        request.Content,
	}

	saved, err := uc.conversationMessageRepository.Save(message)
	if err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}
	return saved, nil
}
