package use_cases

import (
	"fmt"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type IngestMessageInterface interface {
	Ingest(request *requests.IngestMessageRequest) (*responses.ConversationMessageResponse, error)
}

type IngestMessage struct {
	conversationRepository        repositories.ConversationRepositoryInterface
	participantRepository         repositories.ParticipantRepositoryInterface
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface
	participantRoleResolver       interfaces.ParticipantRoleResolverInterface
}

func NewIngestMessage(
	conversationRepository repositories.ConversationRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface,
	participantRoleResolver interfaces.ParticipantRoleResolverInterface,
) *IngestMessage {
	return &IngestMessage{
		conversationRepository:        conversationRepository,
		participantRepository:         participantRepository,
		conversationMessageRepository: conversationMessageRepository,
		participantRoleResolver:       participantRoleResolver,
	}
}

func (uc *IngestMessage) Ingest(request *requests.IngestMessageRequest) (*responses.ConversationMessageResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	conversation, err := uc.conversationRepository.GetById(request.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}
	if conversation == nil {
		return nil, cb_errors.NewNotFoundError("conversation", request.ConversationID.String())
	}

	participant, err := uc.participantRepository.GetById(request.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load participant: %w", err)
	}
	if participant == nil {
		return nil, cb_errors.NewNotFoundError("participant", request.ParticipantID.String())
	}

	humanAuthored, err := uc.isHumanAuthored(participant)
	if err != nil {
		return nil, err
	}

	// Reopen the conversation when a human posts
	if conversation.Status == enums.ConversationStatusConcluded && humanAuthored {
		now := time.Now()
		conversation.Status = enums.ConversationStatusActive
		conversation.UnconcludedAt = &now
		if _, err := uc.conversationRepository.Save(conversation); err != nil {
			return nil, fmt.Errorf("failed to reopen concluded conversation: %w", err)
		}
	}

	message := &models.ConversationMessage{
		ConversationID: request.ConversationID,
		ParticipantID:  request.ParticipantID,
		Content:        request.Content,
	}

	saved, err := uc.conversationMessageRepository.Save(message)
	if err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}
	return responses.ConversationMessageResponseFromModel(saved), nil
}

// isHumanAuthored covers plain users and the agent NFTs a person holds themselves.
func (uc *IngestMessage) isHumanAuthored(participant *models.Participant) (bool, error) {
	role, err := uc.participantRoleResolver.Resolve(participant)
	if err != nil {
		return false, err
	}
	return role == nil || *role == enums.ParticipantRoleHuman, nil
}
