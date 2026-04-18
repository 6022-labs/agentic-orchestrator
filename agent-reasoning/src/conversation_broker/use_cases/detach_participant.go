package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
)

type DetachParticipantInterface interface {
	Detach(request *requests.ConversationParticipantRequest) error
}

type DetachParticipant struct {
	conversationParticipantRepository repositories.ConversationParticipantRepositoryInterface
}

func NewDetachParticipant(conversationParticipantRepository repositories.ConversationParticipantRepositoryInterface) *DetachParticipant {
	return &DetachParticipant{conversationParticipantRepository: conversationParticipantRepository}
}

func (uc *DetachParticipant) Detach(request *requests.ConversationParticipantRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	if err := uc.conversationParticipantRepository.DeleteByConversationIdAndParticipantId(request.ConversationID, request.ParticipantID); err != nil {
		return fmt.Errorf("failed to detach participant: %w", err)
	}
	return nil
}
