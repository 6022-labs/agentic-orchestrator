package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

// ConversationProcessorInterface drives one swarm tick for a (participant,
// conversation) pair: build the turn, invoke the remote agent node, post
// completions back to the originating bridge.
type ConversationProcessorInterface interface {
	ProcessConversation(participant *models.Participant, conversation *models.Conversation) error
}
