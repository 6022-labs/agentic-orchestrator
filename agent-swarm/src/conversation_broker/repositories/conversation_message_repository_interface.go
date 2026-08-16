package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

type ConversationMessageRepositoryInterface interface {
	Save(entity *models.ConversationMessage) (*models.ConversationMessage, error)
	ListByConversationId(conversationId uuid.UUID) ([]*models.ConversationMessage, error)

	// ListProcessableForParticipantAndConversation returns messages the given
	// agent participant should see when computing for a conversation.
	// Filtering follows the slack convention:
	//   - facilitator: every message
	//   - clone/expert: messages whose content mentions the participant's ENS domain
	//   - human: messages from this participant OR mentioning their ENS domain
	ListProcessableForParticipantAndConversation(
		participant *models.Participant,
		role *enums.ParticipantRole,
		conversationId uuid.UUID,
	) ([]*models.ConversationMessage, error)
}
