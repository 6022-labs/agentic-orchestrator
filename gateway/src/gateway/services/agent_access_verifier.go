package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"go.uber.org/zap"
)

const AgentAccessForbiddenCode = "agent_access_forbidden"

// AgentAccessVerifier mirrors the on-chain AgentCollectionV1 access rule against the Transfer-synced registry.
type AgentAccessVerifier struct {
	logger                 *zap.Logger
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewAgentAccessVerifier(
	logger *zap.Logger,
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
) *AgentAccessVerifier {
	return &AgentAccessVerifier{logger: logger, agentRegistryRequester: agentRegistryRequester}
}

func (v *AgentAccessVerifier) EnsureAccess(
	ctx context.Context,
	callerId string,
	agentCollectionAddress string,
	agentCollectionTokenId uint64,
) error {
	agentCollectionAddress = strings.ToLower(agentCollectionAddress)
	agent, err := v.agentRegistryRequester.GetAgent(ctx, agentCollectionAddress, agentCollectionTokenId)
	if err != nil {
		v.logger.Error("Failed to load agent from registry",
			zap.String("agentCollectionAddress", agentCollectionAddress),
			zap.Uint64("agentCollectionTokenId", agentCollectionTokenId),
			zap.Error(err),
		)
		return errors.NewUnavailableError("agent_registry_unavailable", "agent registry is unavailable")
	}
	if agent == nil {
		return errors.NewNotFoundError(
			"agent_not_found",
			fmt.Sprintf("agent not found: %s/%d", agentCollectionAddress, agentCollectionTokenId),
		)
	}

	if !v.hasAccess(agent, callerId) {
		return errors.NewForbiddenError(
			AgentAccessForbiddenCode,
			"wallet is neither the owner nor an address of this agent",
		)
	}

	return nil
}

// hasAccess ignores blank addresses, so an unminted agent without an evm wallet admits nobody.
func (v *AgentAccessVerifier) hasAccess(agent *responses.AgentResponse, callerId string) bool {
	if agent.Owner != nil && *agent.Owner != "" && helpers.NormalizeCallerId(*agent.Owner) == callerId {
		return true
	}
	agentEvmAddress := agent.Addresses[enums.AddressTypeEvm.String()]
	return agentEvmAddress != "" && helpers.NormalizeCallerId(agentEvmAddress) == callerId
}
