package use_cases

import (
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
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

	conversation, err := uc.conversationRepository.GetById(request.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}
	if conversation == nil {
		return nil, cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
	}

	return responses.ConversationResponseFromModel(conversation), nil
}
