package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationRepositoryInterface interface {
	Save(entity *models.Conversation) (*models.Conversation, error)
	GetById(id uuid.UUID) (*models.Conversation, error)
	List(page, elementsPerPage uint) ([]*models.Conversation, error)
	Count() (uint, error)

	// ListProcessableForParticipant returns active conversations the given
	// participant should be considered for in the current swarm tick.
	// Facilitators are additionally throttled: a conversation is returned
	// only when the facilitator's last message in it is older than
	// facilitatorTickIntervalSeconds (or no message at all).
	ListProcessableForParticipant(participant *models.Participant, facilitatorTickIntervalSeconds uint64) ([]*models.Conversation, error)
}
