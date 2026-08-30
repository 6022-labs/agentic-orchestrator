package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// memorySessionNamespace keeps derived uuids from colliding with any other id.
const memorySessionNamespace = "6022-memory-session"

// MemorySessionResolver derives a session for the roles that do not name one.
type MemorySessionResolver struct {
	conversationMemoryBindingRepository repositories.ConversationMemoryBindingRepositoryInterface
	namespace                           uuid.UUID
}

func NewMemorySessionResolver(
	conversationMemoryBindingRepository repositories.ConversationMemoryBindingRepositoryInterface,
) *MemorySessionResolver {
	return &MemorySessionResolver{
		conversationMemoryBindingRepository: conversationMemoryBindingRepository,
		namespace:                           uuid.NewSHA1(uuid.NameSpaceOID, []byte(memorySessionNamespace)),
	}
}

func (r *MemorySessionResolver) Resolve(
	ctx context.Context,
	role *enums.ParticipantRole,
	conversation *models.Conversation,
	participant *models.Participant,
	participants []*models.Participant,
) (*uuid.UUID, error) {
	if role == nil {
		return nil, nil
	}

	binding, err := r.conversationMemoryBindingRepository.GetByConversationAndParticipant(ctx,
		conversation.ID,
		participant.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load memory binding: %w", err)
	}

	switch *role {
	case enums.ParticipantRoleExpert:
		if binding == nil {
			return nil, nil
		}
		return &binding.MemorySessionID, nil

	case enums.ParticipantRoleFacilitator:
		return r.pin(ctx, binding, conversation, participant, conversation.ID)

	case enums.ParticipantRoleClone:
		return r.pin(ctx, binding, conversation, participant, r.agentSetSessionId(participants))

	default:
		return nil, nil
	}
}

func (r *MemorySessionResolver) Bind(
	ctx context.Context,
	conversation *models.Conversation,
	participant *models.Participant,
	sessionId uuid.UUID,
) error {
	return r.conversationMemoryBindingRepository.Bind(ctx, &models.ConversationMemoryBinding{
		ConversationID:  conversation.ID,
		ParticipantID:   participant.ID,
		MemorySessionID: sessionId,
	})
}

// pin moves the binding when the derived key changed under it.
func (r *MemorySessionResolver) pin(
	ctx context.Context,
	binding *models.ConversationMemoryBinding,
	conversation *models.Conversation,
	participant *models.Participant,
	derived uuid.UUID,
) (*uuid.UUID, error) {
	if binding != nil && binding.MemorySessionID == derived {
		return &derived, nil
	}

	next := &models.ConversationMemoryBinding{
		ConversationID:  conversation.ID,
		ParticipantID:   participant.ID,
		MemorySessionID: derived,
	}

	if err := r.conversationMemoryBindingRepository.Bind(ctx, next); err != nil {
		return nil, fmt.Errorf("failed to bind memory session: %w", err)
	}
	return &derived, nil
}

// agentSetSessionId keys a clone on its peers, so one group shares one session.
func (r *MemorySessionResolver) agentSetSessionId(participants []*models.Participant) uuid.UUID {
	keys := lo.FilterMap(participants, func(participant *models.Participant, _ int) (string, bool) {
		if !participant.IsAgent() {
			return "", false
		}
		return fmt.Sprintf(
			"%s:%d",
			strings.ToLower(*participant.CollectionAddress),
			*participant.CollectionTokenId,
		), true
	})
	slices.Sort(keys)

	return uuid.NewSHA1(r.namespace, []byte(strings.Join(keys, "|")))
}
