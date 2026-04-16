package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

// ConversationProcessorInterface drives one swarm tick for a (participant,
// conversation) pair: load context, call agent_compute, post completions back
// to the originating bridge, conclude when the compute response says so.
type ConversationProcessorInterface interface {
	ProcessConversation(participant *models.Participant, conversation *models.Conversation) error
}
