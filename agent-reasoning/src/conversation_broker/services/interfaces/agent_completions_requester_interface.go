package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

// AgentCompletionsRequesterInterface is the broker's outbound port to the
// agent-reasoning compute service. The processor calls it once per swarm tick
// for the agent it just woke up.
type AgentCompletionsRequesterInterface interface {
	Completions(
		platform string,
		agent *models.Participant,
		conversation *models.Conversation,
		messages []*models.ConversationMessage,
		participants []*models.Participant,
		forceCompletion bool,
	) (responses []*string, isFinal bool, err error)
}
