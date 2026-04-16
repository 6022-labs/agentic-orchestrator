package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
)

type ForceConcludeInterface interface {
	Conclude(request *requests.ConversationIdRequest) error
}

type ForceConclude struct {
	conversationConclude interfaces.ConversationConcludeInterface
}

func NewForceConclude(conversationConclude interfaces.ConversationConcludeInterface) *ForceConclude {
	return &ForceConclude{conversationConclude: conversationConclude}
}

func (uc *ForceConclude) Conclude(request *requests.ConversationIdRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return uc.conversationConclude.ConcludeIfExists(request.ConversationID)
}
