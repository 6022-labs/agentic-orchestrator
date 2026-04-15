package use_cases

import (
	"errors"
	"fmt"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"gorm.io/gorm"
)

type ForceConcludeInterface interface {
	Conclude(request *requests.ConversationIdRequest) error
}

type ForceConclude struct {
	conversationRepository repositories.ConversationRepositoryInterface
}

func NewForceConclude(conversationRepository repositories.ConversationRepositoryInterface) *ForceConclude {
	return &ForceConclude{conversationRepository: conversationRepository}
}

func (uc *ForceConclude) Conclude(request *requests.ConversationIdRequest) error {
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

	now := time.Now()
	conversation.Status = enums.ConversationStatusConcluded
	conversation.ConcludedAt = &now
	conversation.UnconcludedAt = nil

	if _, err := uc.conversationRepository.Save(conversation); err != nil {
		return fmt.Errorf("failed to conclude conversation: %w", err)
	}
	return nil
}
