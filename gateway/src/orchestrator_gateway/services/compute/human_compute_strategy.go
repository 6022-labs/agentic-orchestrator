package compute

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"go.uber.org/zap"
)

type HumanComputeStrategyInterface interface {
	Compute(agent *models.Agent, inputs []*models.AgentContextMessage, context *ComputeContext) ([]string, error)
}

type HumanComputeStrategy struct {
	*AbstractAgentSetScopedComputeStrategy
}

func NewHumanComputeStrategy(
	logger *zap.Logger,
	agentContextRequester interfaces.AgentContextRequesterInterface,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
) *HumanComputeStrategy {
	return &HumanComputeStrategy{
		AbstractAgentSetScopedComputeStrategy: NewAbstractAgentSetScopedComputeStrategy(
			logger,
			agentContextRequester,
			agentContextMessageRequester,
		),
	}
}

func (cs *HumanComputeStrategy) Compute(
	agent *models.Agent,
	inputs []*models.AgentContextMessage,
	context *ComputeContext,
) ([]string, error) {
	cs.logger.Debug("Computing human strategy",
		zap.String("agentId", agent.ID.String()),
		zap.String("latestMessageAuthor", inputs[len(inputs)-1].Author),
		zap.String("conversationContext.Id", context.ConversationContext.Id),
	)

	key := cs.getContextKey(context.ConversationContext.Agents)

	agentContext, isFound, err := cs.agentContextRequester.GetAgentContextByKey(agent.ID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context: %w", err)
	}

	if !*isFound {
		agentContext, err = cs.createAgentContext(agent, key)
		if err != nil {
			return nil, fmt.Errorf("failed to create agent context: %w", err)
		}
	}

	storedMessages, err := cs.agentContextMessageRequester.GetAgentContextMessages(agent.ID, agentContext.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context messages: %w", err)
	}

	messages := cs.unionMessagesByHash(agentContext.Id, inputs, storedMessages)
	err = cs.saveMessages(agent, agentContext.Key, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to save messages: %w", err)
	}

	return nil, nil
}
