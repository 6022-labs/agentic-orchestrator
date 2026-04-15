package services

import (
	"fmt"
	"strconv"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_identity_registry/clients"
)

type AgentRegistryRequester struct {
	client clients.AgentIdentityRegistryClientInterface
}

func NewAgentRegistryRequester(client clients.AgentIdentityRegistryClientInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{client: client}
}

func (r *AgentRegistryRequester) GetMintedAgent(collectionAddress string, collectionTokenId uint64) (*models.Participant, error) {
	external, err := r.client.GetAgent(collectionAddress, collectionTokenId)
	if err != nil {
		return nil, err
	}
	if external == nil {
		return nil, nil
	}

	var role enums.ParticipantRole
	if err := role.Set(external.Role); err != nil {
		return nil, fmt.Errorf("registry returned invalid role: %w", err)
	}

	collectionAddressCopy := external.CollectionAddress
	ensDomainCopy := external.EnsDomain

	var collectionTokenIdStr *string
	if external.CollectionTokenId != nil {
		s := strconv.FormatUint(*external.CollectionTokenId, 10)
		collectionTokenIdStr = &s
	}

	return &models.Participant{
		ID:                external.ID,
		CollectionAddress: &collectionAddressCopy,
		CollectionTokenId: collectionTokenIdStr,
		EnsDomain:         &ensDomainCopy,
		Role:              &role,
	}, nil
}
