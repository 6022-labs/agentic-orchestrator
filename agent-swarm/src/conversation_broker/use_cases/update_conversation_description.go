package use_cases

import (
	"context"

	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
)

type UpdateConversationDescriptionInterface interface {
	Update(ctx context.Context, request *requests.UpdateConversationDescriptionRequest) error
}

type UpdateConversationDescription struct {
	conversationRepository repositories.ConversationRepositoryInterface
}

func NewUpdateConversationDescription(
	conversationRepository repositories.ConversationRepositoryInterface,
) *UpdateConversationDescription {
	return &UpdateConversationDescription{conversationRepository: conversationRepository}
}

func (uc *UpdateConversationDescription) Update(
	ctx context.Context,
	request *requests.UpdateConversationDescriptionRequest,
) error {
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

	conversation.Description = request.Description

	if _, err := uc.conversationRepository.Save(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation description: %w", err)
	}
	return nil
}
