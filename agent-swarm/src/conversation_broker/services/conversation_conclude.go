package services

import (
	"context"

	"fmt"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/google/uuid"
)

type ConversationConclude struct {
	conversationRepository repositories.ConversationRepositoryInterface
}

func NewConversationConclude(
	conversationRepository repositories.ConversationRepositoryInterface,
) *ConversationConclude {
	return &ConversationConclude{conversationRepository: conversationRepository}
}

func (s *ConversationConclude) Conclude(ctx context.Context, conversationId uuid.UUID) error {
	conversation, err := s.conversationRepository.GetById(ctx, conversationId)
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}
	if conversation == nil {
		return cb_errors.NewNotFoundError("conversation", conversationId.String())
	}
	return s.concludeConversation(ctx, conversation)
}

func (s *ConversationConclude) ConcludeIfExists(ctx context.Context, conversationId uuid.UUID) error {
	conversation, err := s.conversationRepository.GetById(ctx, conversationId)
	if err != nil {
		return fmt.Errorf("failed to load conversation: %w", err)
	}
	if conversation == nil {
		return nil
	}
	return s.concludeConversation(ctx, conversation)
}

func (s *ConversationConclude) concludeConversation(ctx context.Context, conversation *models.Conversation) error {
	if conversation.Status == enums.ConversationStatusConcluded {
		return nil
	}

	now := time.Now()
	conversation.Status = enums.ConversationStatusConcluded
	conversation.ConcludedAt = &now
	conversation.UnconcludedAt = nil

	if _, err := s.conversationRepository.Save(ctx, conversation); err != nil {
		return fmt.Errorf("failed to conclude conversation: %w", err)
	}
	return nil
}
