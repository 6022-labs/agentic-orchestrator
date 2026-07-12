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

type GetConversationInterface interface {
	Get(request *requests.ConversationIdRequest) (*responses.ConversationResponse, error)
}

type GetConversation struct {
	conversationRepository repositories.ConversationRepositoryInterface
}

func NewGetConversation(conversationRepository repositories.ConversationRepositoryInterface) *GetConversation {
	return &GetConversation{conversationRepository: conversationRepository}
}

func (uc *GetConversation) Get(request *requests.ConversationIdRequest) (*responses.ConversationResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	conv, err := uc.conversationRepository.GetById(request.ConversationID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}

	return responses.ConversationResponseFromModel(conv), nil
}
