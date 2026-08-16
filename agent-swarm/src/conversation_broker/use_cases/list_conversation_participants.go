package use_cases

import (
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type ListConversationParticipantsInterface interface {
	List(request *requests.ConversationIdRequest) ([]*responses.ParticipantResponse, error)
}

type ListConversationParticipants struct {
	conversationRepository repositories.ConversationRepositoryInterface
	participantRepository  repositories.ParticipantRepositoryInterface
}

func NewListConversationParticipants(
	conversationRepository repositories.ConversationRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
) *ListConversationParticipants {
	return &ListConversationParticipants{
		conversationRepository: conversationRepository,
		participantRepository:  participantRepository,
	}
}

func (uc *ListConversationParticipants) List(
	request *requests.ConversationIdRequest,
) ([]*responses.ParticipantResponse, error) {
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

	participants, err := uc.participantRepository.ListByConversationId(request.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversation participants: %w", err)
	}

	participantResponses := make([]*responses.ParticipantResponse, len(participants))
	for i, participant := range participants {
		participantResponses[i] = responses.ParticipantResponseFromModel(participant)
	}

	return participantResponses, nil
}
