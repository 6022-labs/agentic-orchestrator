package use_cases

import (
	"context"

	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type RegisterConversationInterface interface {
	Register(
		ctx context.Context,
		request *requests.RegisterConversationRequest,
	) (*responses.ConversationResponse, error)
}

type RegisterConversation struct {
	bridgeRepository                  repositories.BridgeRepositoryInterface
	conversationRepository            repositories.ConversationRepositoryInterface
	conversationParticipantRepository repositories.ConversationParticipantRepositoryInterface
	participantRepository             repositories.ParticipantRepositoryInterface
}

func NewRegisterConversation(
	bridgeRepository repositories.BridgeRepositoryInterface,
	conversationRepository repositories.ConversationRepositoryInterface,
	conversationParticipantRepository repositories.ConversationParticipantRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
) *RegisterConversation {
	return &RegisterConversation{
		bridgeRepository:                  bridgeRepository,
		conversationRepository:            conversationRepository,
		conversationParticipantRepository: conversationParticipantRepository,
		participantRepository:             participantRepository,
	}
}

func (uc *RegisterConversation) Register(
	ctx context.Context,
	request *requests.RegisterConversationRequest,
) (*responses.ConversationResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	bridge, err := uc.bridgeRepository.GetById(ctx, request.BridgeID)
	if err != nil {
		return nil, fmt.Errorf("failed to load bridge: %w", err)
	}
	if bridge == nil {
		return nil, cb_errors.NewNotFoundError("bridge", request.BridgeID.String())
	}

	for _, pid := range request.ParticipantIDs {
		participant, err := uc.participantRepository.GetById(ctx, pid)
		if err != nil {
			return nil, fmt.Errorf("failed to load participant %s: %w", pid, err)
		}
		if participant == nil {
			return nil, cb_errors.NewNotFoundError("participant", pid.String())
		}
	}

	conversation := &models.Conversation{
		BridgeID:    request.BridgeID,
		Description: request.Description,
	}
	saved, err := uc.conversationRepository.Save(ctx, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}

	for _, pid := range request.ParticipantIDs {
		link := &models.ConversationParticipant{
			ConversationID: saved.ID,
			ParticipantID:  pid,
		}
		if _, err := uc.conversationParticipantRepository.Save(ctx, link); err != nil {
			return nil, fmt.Errorf("failed to attach participant %s: %w", pid, err)
		}
	}

	return responses.ConversationResponseFromModel(saved), nil
}
