package services

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
)

// ParticipantRoleResolver reads the role from the registry, which owns the attributes.
type ParticipantRoleResolver struct {
	registryAgentResolver interfaces.RegistryAgentResolverInterface
}

func NewParticipantRoleResolver(
	registryAgentResolver interfaces.RegistryAgentResolverInterface,
) *ParticipantRoleResolver {
	return &ParticipantRoleResolver{
		registryAgentResolver: registryAgentResolver,
	}
}

func (r *ParticipantRoleResolver) Resolve(
	ctx context.Context,
	participant *models.Participant,
) (*enums.ParticipantRole, error) {
	if participant == nil || !participant.IsAgent() {
		return nil, nil
	}

	agent, err := r.registryAgentResolver.ResolveByCollection(ctx,
		*participant.CollectionAddress,
		*participant.CollectionTokenId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve role for participant %s: %w", participant.ID, err)
	}

	// An unknown agent gets the same default an unset attribute would.
	var attributes map[string]string
	if agent != nil {
		attributes = agent.Attributes
	}
	role := enums.ParticipantRoleFromAttributes(attributes)
	return &role, nil
}
