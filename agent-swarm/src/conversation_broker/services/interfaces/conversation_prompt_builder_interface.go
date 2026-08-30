package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

// ConversationPromptBuilderInterface prepares the inference turn for an agent:
// role-templated instructions plus the conversation window as input messages.
type ConversationPromptBuilderInterface interface {
	// BuildInput renders the turn; a nil memorySessionId means the expert has not bound a session yet.
	BuildInput(
		ctx context.Context,
		platform string,
		agent *models.Participant,
		role *enums.ParticipantRole,
		conversation *models.Conversation,
		messages []*models.ConversationMessage,
		participants []*models.Participant,
		memorySessionId *uuid.UUID,
	) (*models.RemoteAgentInput, error)
	// BuildConcludeInput extends a turn with the agent's completions and the facilitator conclude prompt.
	BuildConcludeInput(
		ctx context.Context,
		turn *models.RemoteAgentInput,
		completions []string,
		agent *models.Participant,
		role *enums.ParticipantRole,
		conversation *models.Conversation,
		messages []*models.ConversationMessage,
		participants []*models.Participant,
		memorySessionId *uuid.UUID,
	) (*models.RemoteAgentInput, error)
}
