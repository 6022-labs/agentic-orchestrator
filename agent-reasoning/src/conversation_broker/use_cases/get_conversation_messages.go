package use_cases

import (
	"errors"
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/responses"
	"gorm.io/gorm"
)

type GetConversationMessagesInterface interface {
	Get(request *requests.ConversationIdRequest) ([]*responses.ConversationMessageResponse, error)
}

type GetConversationMessages struct {
	conversationRepository        repositories.ConversationRepositoryInterface
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface
}

func NewGetConversationMessages(
	conversationRepository repositories.ConversationRepositoryInterface,
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface,
) *GetConversationMessages {
	return &GetConversationMessages{
		conversationRepository:        conversationRepository,
		conversationMessageRepository: conversationMessageRepository,
	}
}

func (uc *GetConversationMessages) Get(request *requests.ConversationIdRequest) ([]*responses.ConversationMessageResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	if _, err := uc.conversationRepository.GetById(request.ConversationID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
		}
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}

	messages, err := uc.conversationMessageRepository.ListByConversationId(request.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversation messages: %w", err)
	}

	messageResponses := make([]*responses.ConversationMessageResponse, len(messages))
	for i, message := range messages {
		messageResponses[i] = responses.ConversationMessageResponseFromModel(message)
	}

	return messageResponses, nil
}
