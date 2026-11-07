package use_cases

import (
	"context"

	"fmt"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type IngestMessageInterface interface {
	Ingest(ctx context.Context, request *requests.IngestMessageRequest) (*responses.ConversationMessageResponse, error)
}

type IngestMessage struct {
	brokerMetrics                 metrics_interfaces.BrokerMetricsInterface
	conversationRepository        repositories.ConversationRepositoryInterface
	participantRepository         repositories.ParticipantRepositoryInterface
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface
	participantRoleResolver       interfaces.ParticipantRoleResolverInterface
}

func NewIngestMessage(
	brokerMetrics metrics_interfaces.BrokerMetricsInterface,
	conversationRepository repositories.ConversationRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface,
	participantRoleResolver interfaces.ParticipantRoleResolverInterface,
) *IngestMessage {
	return &IngestMessage{
		brokerMetrics:                 brokerMetrics,
		conversationRepository:        conversationRepository,
		participantRepository:         participantRepository,
		conversationMessageRepository: conversationMessageRepository,
		participantRoleResolver:       participantRoleResolver,
	}
}

func (uc *IngestMessage) Ingest(
	ctx context.Context,
	request *requests.IngestMessageRequest,
) (*responses.ConversationMessageResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	conversation, err := uc.conversationRepository.GetById(ctx, request.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %w", err)
	}
	if conversation == nil {
		return nil, cb_errors.NewNotFoundError(
			"conversation_not_found",
			"conversation not found: "+request.ConversationID.String(),
		)
	}

	participant, err := uc.participantRepository.GetById(ctx, request.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load participant: %w", err)
	}
	if participant == nil {
		return nil, cb_errors.NewNotFoundError(
			"participant_not_found",
			"participant not found: "+request.ParticipantID.String(),
		)
	}

	role, err := uc.participantRoleResolver.Resolve(ctx, participant)
	if err != nil {
		return nil, err
	}
	humanAuthored := uc.isHumanAuthored(role)

	// Reopen the conversation when a human posts
	if conversation.Status == enums.ConversationStatusConcluded && humanAuthored {
		now := time.Now()
		conversation.Status = enums.ConversationStatusActive
		conversation.UnconcludedAt = &now
		if _, err := uc.conversationRepository.Save(ctx, conversation); err != nil {
			return nil, fmt.Errorf("failed to reopen concluded conversation: %w", err)
		}
	}

	message := &models.ConversationMessage{
		ConversationID: request.ConversationID,
		ParticipantID:  request.ParticipantID,
		Content:        request.Content,
	}

	saved, err := uc.conversationMessageRepository.Save(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	uc.brokerMetrics.RecordMessageIngested(ctx, uc.participantRoleLabel(role))

	return responses.ConversationMessageResponseFromModel(saved), nil
}

// isHumanAuthored covers plain users and the agent NFTs a person holds themselves.
func (uc *IngestMessage) isHumanAuthored(role *enums.ParticipantRole) bool {
	return role == nil || *role == enums.ParticipantRoleHuman
}

func (uc *IngestMessage) participantRoleLabel(role *enums.ParticipantRole) string {
	if role == nil {
		return "unknown"
	}
	return string(*role)
}
