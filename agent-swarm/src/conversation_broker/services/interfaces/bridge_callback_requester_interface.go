package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"

// BridgeCallbackRequesterInterface is the broker's outbound port for delivering
// agent completions back to the originating bridge. The bridge's callback URL
// is read from the bridge row passed in.
type BridgeCallbackRequesterInterface interface {
	PostCompletion(
		bridge *models.Bridge,
		conversation *models.Conversation,
		agent *models.Participant,
		content string,
	) error
}
