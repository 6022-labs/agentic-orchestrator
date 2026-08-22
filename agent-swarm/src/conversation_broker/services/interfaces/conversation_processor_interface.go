package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
)

// ConversationProcessorInterface drives one swarm tick for a (participant,
// conversation) pair: build the turn, invoke the remote agent node, post
// completions back to the originating bridge.
type ConversationProcessorInterface interface {
	ProcessConversation(ctx context.Context, participant *models.Participant, conversation *models.Conversation) error
}
