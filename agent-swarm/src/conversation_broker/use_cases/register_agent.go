package use_cases

import (
	"context"

	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type RegisterAgentInterface interface {
	Register(ctx context.Context, request *requests.RegisterAgentRequest) (*responses.ParticipantResponse, bool, error)
}

// RegisterAgent upserts an agent participant by its on-chain identity bundle
// (collectionAddress, collectionTokenId). Returns (existing, true, nil) if
// already known; otherwise resolves identity from the agent-identity registry,
// saves, and returns (created, false, nil).
type RegisterAgent struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
	participantRepository  repositories.ParticipantRepositoryInterface
}

func NewRegisterAgent(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
) *RegisterAgent {
	return &RegisterAgent{
		agentRegistryRequester: agentRegistryRequester,
		participantRepository:  participantRepository,
	}
}

func (uc *RegisterAgent) Register(
	ctx context.Context,
	request *requests.RegisterAgentRequest,
) (*responses.ParticipantResponse, bool, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, false, err
	}

	existing, err := uc.participantRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		request.CollectionAddress,
		request.CollectionTokenId,
	)
	if err != nil {
		return nil, false, fmt.Errorf("failed to lookup agent participant: %w", err)
	}
	if existing != nil {
		return responses.ParticipantResponseFromModel(existing), true, nil
	}

	registryAgent, err := uc.agentRegistryRequester.GetAgent(ctx, request.CollectionAddress, request.CollectionTokenId)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch agent from registry: %w", err)
	}
	if registryAgent == nil {
		agentKey := fmt.Sprintf("%s/%d", request.CollectionAddress, request.CollectionTokenId)
		return nil, false, cb_errors.NewNotFoundError("agent_not_found", "agent not found: "+agentKey)
	}

	participant := registryAgent.ToParticipant()
	participant.CollectionAddress = &request.CollectionAddress

	saved, err := uc.participantRepository.Save(ctx, participant)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save agent participant: %w", err)
	}
	return responses.ParticipantResponseFromModel(saved), false, nil
}
