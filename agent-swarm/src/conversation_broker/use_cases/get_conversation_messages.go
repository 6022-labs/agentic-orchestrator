package use_cases

import (
	"context"

	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type GetConversationMessagesInterface interface {
	Get(ctx context.Context, request *requests.ConversationIdRequest) ([]*responses.ConversationMessageResponse, error)
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

func (uc *GetConversationMessages) Get(
	ctx context.Context,
	request *requests.ConversationIdRequest,
) ([]*responses.ConversationMessageResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	conversation, err := uc.conversationRepository.GetById(ctx, request.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}
	if conversation == nil {
		return nil, cb_errors.NewNotFoundError(
			"conversation_not_found",
			"conversation not found: "+request.ConversationID.String(),
		)
	}

	messages, err := uc.conversationMessageRepository.ListByConversationId(ctx, request.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversation messages: %w", err)
	}

	messageResponses := make([]*responses.ConversationMessageResponse, len(messages))
	for i, message := range messages {
		messageResponses[i] = responses.ConversationMessageResponseFromModel(message)
	}

	return messageResponses, nil
}
