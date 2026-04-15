package use_cases

import (
	"errors"
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"gorm.io/gorm"
)

type UpdateConversationDescriptionInterface interface {
	Update(request *requests.UpdateConversationDescriptionRequest) error
}

type UpdateConversationDescription struct {
	conversationRepository repositories.ConversationRepositoryInterface
}

func NewUpdateConversationDescription(conversationRepository repositories.ConversationRepositoryInterface) *UpdateConversationDescription {
	return &UpdateConversationDescription{conversationRepository: conversationRepository}
}

func (uc *UpdateConversationDescription) Update(request *requests.UpdateConversationDescriptionRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	conversation, err := uc.conversationRepository.GetById(request.ConversationID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
	}
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}

	conversation.Description = request.Description

	if _, err := uc.conversationRepository.Save(conversation); err != nil {
		return fmt.Errorf("failed to update conversation description: %w", err)
	}
	return nil
}
