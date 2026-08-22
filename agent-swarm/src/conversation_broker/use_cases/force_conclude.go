package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
)

type ForceConcludeInterface interface {
	Conclude(ctx context.Context, request *requests.ConversationIdRequest) error
}

type ForceConclude struct {
	conversationConclude interfaces.ConversationConcludeInterface
}

func NewForceConclude(conversationConclude interfaces.ConversationConcludeInterface) *ForceConclude {
	return &ForceConclude{conversationConclude: conversationConclude}
}

func (uc *ForceConclude) Conclude(ctx context.Context, request *requests.ConversationIdRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return uc.conversationConclude.ConcludeIfExists(ctx, request.ConversationID)
}
