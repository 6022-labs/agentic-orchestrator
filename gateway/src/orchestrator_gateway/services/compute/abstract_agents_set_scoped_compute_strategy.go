package compute

import (
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type AbstractAgentSetScopedComputeStrategy struct {
	*AbstractComputeStrategy
	agentContextRequester interfaces.AgentContextRequesterInterface
}

func NewAbstractAgentSetScopedComputeStrategy(
	logger *zap.Logger,
	agentContextRequester interfaces.AgentContextRequesterInterface,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
) *AbstractAgentSetScopedComputeStrategy {
	return &AbstractAgentSetScopedComputeStrategy{
		AbstractComputeStrategy: NewAbstractComputeStrategy(logger, agentContextMessageRequester),
		agentContextRequester:   agentContextRequester,
	}
}

func (cs *AbstractAgentSetScopedComputeStrategy) createAgentContext(
	agent *models.Agent,
	key string,
) (*models.AgentContext, error) {
	err := cs.agentContextRequester.CreateAgentContext(agent.ID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent context: %w", err)
	}

	agentContext, isFound, err := cs.agentContextRequester.GetAgentContextByKey(agent.ID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context by key: %w", err)
	}
	if !*isFound {
		return nil, fmt.Errorf("agent context not found after creation")
	}

	return agentContext, nil
}

func (cs *AbstractAgentSetScopedComputeStrategy) getContextKey(
	agents []*ConversationContextAgent,
) string {
	agentKeys := lo.Map(agents, func(a *ConversationContextAgent, _ int) string {
		return fmt.Sprintf("%s:%d", strings.ToLower(a.CollectionAddress), a.CollectionTokenId)
	})

	slices.SortFunc(agentKeys, func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	baseKey := strings.Join(agentKeys, "|")
	sum := sha256.Sum256([]byte(baseKey))

	return fmt.Sprintf("%x", sum)
}
