package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

// ConversationPromptBuilderInterface prepares the inference turn for an agent:
// role-templated instructions plus the conversation window as input messages.
// BuildConcludeInput extends a turn with the agent's completions and the
// facilitator conclude prompt.
type ConversationPromptBuilderInterface interface {
	BuildInput(
		platform string,
		agent *models.Participant,
		conversation *models.Conversation,
		messages []*models.ConversationMessage,
		participants []*models.Participant,
	) (*models.RemoteAgentInput, error)
	BuildConcludeInput(
		turn *models.RemoteAgentInput,
		completions []string,
		agent *models.Participant,
		conversation *models.Conversation,
		messages []*models.ConversationMessage,
		participants []*models.Participant,
	) (*models.RemoteAgentInput, error)
}
