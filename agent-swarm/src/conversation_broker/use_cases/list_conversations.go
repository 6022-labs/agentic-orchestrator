package use_cases

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type ListConversationsInterface interface {
	List(ctx context.Context, request *requests.ListConversationsRequest) (*responses.ListConversationsResponse, error)
}

type ListConversations struct {
	conversationRepository repositories.ConversationRepositoryInterface
}

func NewListConversations(conversationRepository repositories.ConversationRepositoryInterface) *ListConversations {
	return &ListConversations{conversationRepository: conversationRepository}
}

func (uc *ListConversations) List(
	ctx context.Context,
	request *requests.ListConversationsRequest,
) (*responses.ListConversationsResponse, error) {
	request.Sanitize()

	conversations, err := uc.conversationRepository.List(ctx, *request.Page, *request.ElementsPerPage)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}

	totalElements, err := uc.conversationRepository.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count conversations: %w", err)
	}

	conversationResponses := make([]*responses.ConversationResponse, len(conversations))
	for i, conversation := range conversations {
		conversationResponses[i] = responses.ConversationResponseFromModel(conversation)
	}

	return &responses.ListConversationsResponse{
		Page:            *request.Page,
		ElementsPerPage: *request.ElementsPerPage,
		TotalElements:   totalElements,
		Conversations:   conversationResponses,
	}, nil
}
