package use_cases

import (
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type RegisterAgentInterface interface {
	Register(request *requests.RegisterAgentRequest) (*responses.ParticipantResponse, bool, error)
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

func (uc *RegisterAgent) Register(request *requests.RegisterAgentRequest) (*responses.ParticipantResponse, bool, error) {
	if err := request.Validate(); err != nil {
		return nil, false, err
	}

	tokenId := request.ParsedCollectionTokenId()
	existing, err := uc.participantRepository.GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, tokenId)
	if err != nil {
		return nil, false, fmt.Errorf("failed to lookup agent participant: %w", err)
	}
	if existing != nil {
		return responses.ParticipantResponseFromModel(existing), true, nil
	}

	registryAgent, err := uc.agentRegistryRequester.GetAgent(request.CollectionAddress, tokenId)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch agent from registry: %w", err)
	}
	if registryAgent == nil {
		return nil, false, cb_errors.NewNotFoundError("agent", request.CollectionAddress+"/"+request.CollectionTokenId)
	}

	participant, err := registryAgent.ToParticipant()
	if err != nil {
		return nil, false, err
	}

	saved, err := uc.participantRepository.Save(participant)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save agent participant: %w", err)
	}
	return responses.ParticipantResponseFromModel(saved), false, nil
}
