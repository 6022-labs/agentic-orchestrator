package interfaces

import "context"

// ConversationProcessorDaemonInterface drives the swarm tick at process level:
// it wakes participants one by one and dispatches their processable
// conversations to the ConversationProcessor. Start is idempotent; Stop
// cancels the internal context and waits for in-flight processing to drain.
type ConversationProcessorDaemonInterface interface {
	Start(ctx context.Context)
	Stop()
}
