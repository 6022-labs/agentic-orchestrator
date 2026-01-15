package compute

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type FacilitatorComputeStrategyInterface interface {
	Compute(agent *models.Agent, inputs []*models.AgentContextMessage, context *ComputeContext) ([]string, *bool, error)
}

type FacilitatorComputeStrategy struct {
	*AbstractComputeStrategy
	agentContextRequester                interfaces.AgentContextRequesterInterface
	agentPromptTemplateRequester         interfaces.AgentPromptTemplateRequesterInterface
	promptParametersBuilderFactory       interfaces.PromptParametersBuilderFactoryInterface
	promptFromPromptTemplateParser       interfaces.PromptFromPromptTemplateParserInterface
	llmConnectionChatCompletionRequester interfaces.LLMConnectionChatCompletionRequesterInterface
}

func NewFacilitatorComputeStrategy(
	logger *zap.Logger,
	agentContextRequester interfaces.AgentContextRequesterInterface,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface,
	promptParametersBuilderFactory interfaces.PromptParametersBuilderFactoryInterface,
	promptFromPromptTemplateParser interfaces.PromptFromPromptTemplateParserInterface,
	llmConnectionChatCompletionRequester interfaces.LLMConnectionChatCompletionRequesterInterface,
) *FacilitatorComputeStrategy {
	return &FacilitatorComputeStrategy{
		AbstractComputeStrategy:              NewAbstractComputeStrategy(logger, agentContextMessageRequester),
		agentContextRequester:                agentContextRequester,
		agentPromptTemplateRequester:         agentPromptTemplateRequester,
		promptFromPromptTemplateParser:       promptFromPromptTemplateParser,
		promptParametersBuilderFactory:       promptParametersBuilderFactory,
		llmConnectionChatCompletionRequester: llmConnectionChatCompletionRequester,
	}
}

func (cs *FacilitatorComputeStrategy) Compute(agent *models.Agent, inputs []*models.AgentContextMessage, context *ComputeContext) ([]string, *bool, error) {
	cs.logger.Debug("Computing facilitator strategy",
		zap.String("agentId", agent.ID.String()),
		zap.String("latestMessageAuthor", inputs[len(inputs)-1].Author),
		zap.String("conversationContext.Id", context.ConversationContext.Id),
	)

	key := cs.getContextKey(context.ConversationContext.Id)

	agentContext, isFound, err := cs.agentContextRequester.GetAgentContextByKey(agent.ID, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent context: %w", err)
	}

	if !*isFound {
		agentContext, err = cs.createAgentContext(agent, key)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create agent context: %w", err)
		}
	}

	storedMessages, err := cs.agentContextMessageRequester.GetAgentContextMessages(agent.ID, agentContext.Key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent context messages: %w", err)
	}

	messages := cs.unionMessagesByHash(agentContext.Id, inputs, storedMessages)
	messages = cs.groupMessages(messages, context.ConversationContext.Id)

	// In case of a new conversation, we need to create the system prompt
	if !*isFound {
		messages, err = cs.addSystemPrompt(agent, context, inputs[0], messages)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to add system prompt: %w", err)
		}
	}

	var defaultPromptCompletion *string
	defaultPromptCompletion, messages, err = cs.processDefaultPrompt(agent, context, messages, inputs[len(inputs)-1], inputs[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process default prompt: %w", err)
	}

	isAgentReadyToConclude := cs.isAgentReadyToConclude(messages)
	result := []string{*defaultPromptCompletion}

	isAgentConcluded := false

	if isAgentReadyToConclude || context.ForceConclusion {
		isAgentConcluded = true
		var concludePromptCompletion *string
		concludePromptCompletion, messages, err = cs.processConcludePrompt(agent, context, messages, inputs[len(inputs)-1], inputs[0])
		if err != nil {
			return nil, nil, fmt.Errorf("failed to process conclude prompt: %w", err)
		}

		result = append(result, *concludePromptCompletion)
	}

	nonStoredMessages := cs.diffMessagesByHash(agentContext.Id, messages, storedMessages)
	err = cs.saveMessages(agent, agentContext.Key, nonStoredMessages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save non-stored messages: %w", err)
	}

	return result, &isAgentConcluded, nil
}

func (cs *FacilitatorComputeStrategy) createAgentContext(
	agent *models.Agent,
	key string,
) (*models.AgentContext, error) {
	err := cs.agentContextRequester.CreateAgentContext(agent.ID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent context: %w", err)
	}

	agentContext, isFound, err := cs.agentContextRequester.GetAgentContextByKey(agent.ID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked agent context key of conversation: %w", err)
	}
	if !*isFound {
		return nil, fmt.Errorf("failed to find agent context after creating it")
	}

	return agentContext, nil
}

func (cs *FacilitatorComputeStrategy) addSystemPrompt(
	agent *models.Agent,
	computeContext *ComputeContext,
	initialRequest *models.AgentContextMessage,
	messages []*models.AgentContextMessage,
) ([]*models.AgentContextMessage, error) {
	cs.logger.Debug("Adding system prompt",
		zap.String("agentId", agent.ID.String()),
		zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
	)

	promptParameterBuilder := cs.promptParametersBuilderFactory.Create()
	promptParameters, err := promptParameterBuilder.
		WithBase(computeContext.PromptParameters).
		WithParam([]string{"InitialRequest", "Content"}, initialRequest.Content, nil).
		WithParam([]string{"ConversationId"}, computeContext.ConversationContext.Id, nil).
		WithParam([]string{"InitialRequest", "Author"}, initialRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt parameters: %w", err)
	}

	systemPrompt, err := cs.promptFromPromptTemplateParser.ParseFacilitatorSystem(agent.ID, promptParameters)
	if err != nil {
		return nil, fmt.Errorf("failed to compute prompt template: %w", err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Content:        *systemPrompt,
		Author:         "", // Author of the prompt is not the agent
		Role:           enums.AgentContextMessageRoleUser,
		ConversationId: &computeContext.ConversationContext.Id,
		Platform:       enums.AgentContextMessagePlatformNone,
	})

	return messages, nil
}

func (cs *FacilitatorComputeStrategy) processDefaultPrompt(
	agent *models.Agent,
	computeContext *ComputeContext,
	messages []*models.AgentContextMessage,
	latestRequest *models.AgentContextMessage,
	initialRequest *models.AgentContextMessage,
) (*string, []*models.AgentContextMessage, error) {
	cs.logger.Debug("Processing default prompt",
		zap.String("agentId", agent.ID.String()),
		zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
	)

	promptParameterBuilder := cs.promptParametersBuilderFactory.Create()
	promptParameters, err := promptParameterBuilder.
		WithParam([]string{"LatestRequest", "Content"}, latestRequest.Content, nil).
		WithParam([]string{"InitialRequest", "Content"}, initialRequest.Content, nil).
		WithParam([]string{"ConversationId"}, computeContext.ConversationContext.Id, nil).
		WithParam([]string{"LatestRequest", "Author"}, latestRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"InitialRequest", "Author"}, initialRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"ConversationMessageCount"}, len(messages), nil).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build prompt parameters: %w", err)
	}

	defaultPrompt, err := cs.promptFromPromptTemplateParser.ParseFacilitatorDefault(agent.ID, promptParameters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute prompt template: %w", err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Content:        *defaultPrompt,
		Author:         "", // Author of the prompt is not the agent
		ConversationId: &computeContext.ConversationContext.Id,
		Role:           enums.AgentContextMessageRoleUser,
		Platform:       enums.AgentContextMessagePlatformNone,
	})

	defaultPromptCompletion, err := cs.llmConnectionChatCompletionRequester.ChatCompletion(agent.ID, messages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get completion for prompt template key %v: %w", enums.AgentPromptTemplateKeyFacilitatorDefault, err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Author:         agent.EnsDomain,
		Content:        *defaultPromptCompletion,
		ConversationId: &computeContext.ConversationContext.Id,
		Platform:       enums.AgentContextMessagePlatformNone,
		Role:           enums.AgentContextMessageRoleAssistant,
	})

	return defaultPromptCompletion, messages, nil
}

func (cs *FacilitatorComputeStrategy) isAgentReadyToConclude(
	messages []*models.AgentContextMessage,
) bool {
	messagesToConsider := lo.Filter(messages, func(message *models.AgentContextMessage, _ int) bool {
		return message.Role == enums.AgentContextMessageRoleAssistant ||
			(message.Role == enums.AgentContextMessageRoleUser && message.Platform != enums.AgentContextMessagePlatformNone)
	})

	messagesToConsider = messagesToConsider[len(messagesToConsider)-2:]

	cs.logger.Debug("Checking if agent is ready to conclude",
		zap.Any("messagesToConsider", messagesToConsider),
	)

	for _, message := range messagesToConsider {
		if strings.Contains(strings.ToLower(message.Content), ":end:") {
			return true
		}
	}

	return false
}

func (cs *FacilitatorComputeStrategy) processConcludePrompt(
	agent *models.Agent,
	computeContext *ComputeContext,
	messages []*models.AgentContextMessage,
	latestRequest *models.AgentContextMessage,
	initialRequest *models.AgentContextMessage,
) (*string, []*models.AgentContextMessage, error) {
	cs.logger.Debug("Processing conclude prompt",
		zap.String("agentId", agent.ID.String()),
		zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
	)

	promptParameterBuilder := cs.promptParametersBuilderFactory.Create()
	promptParameters, err := promptParameterBuilder.
		WithParam([]string{"LatestRequest", "Content"}, latestRequest.Content, nil).
		WithParam([]string{"InitialRequest", "Content"}, initialRequest.Content, nil).
		WithParam([]string{"ConversationId"}, computeContext.ConversationContext.Id, nil).
		WithParam([]string{"LatestRequest", "Author"}, latestRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"InitialRequest", "Author"}, initialRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"ConversationMessageCount"}, len(messages), nil).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build prompt parameters: %w", err)
	}

	concludePrompt, err := cs.promptFromPromptTemplateParser.ParseFacilitatorConclude(agent.ID, promptParameters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute prompt template: %w", err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Content:        *concludePrompt,
		Author:         "", // Author of the prompt is not the agent
		ConversationId: &computeContext.ConversationContext.Id,
		Role:           enums.AgentContextMessageRoleUser,
		Platform:       enums.AgentContextMessagePlatformNone,
	})

	concludePromptCompletion, err := cs.llmConnectionChatCompletionRequester.ChatCompletion(agent.ID, messages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get completion for prompt template key %v: %w", enums.AgentPromptTemplateKeyFacilitatorConclude, err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Author:         agent.EnsDomain,
		Content:        *concludePromptCompletion,
		ConversationId: &computeContext.ConversationContext.Id,
		Platform:       enums.AgentContextMessagePlatformNone,
		Role:           enums.AgentContextMessageRoleAssistant,
	})

	return concludePromptCompletion, messages, nil
}

func (cs *FacilitatorComputeStrategy) getContextKey(
	conversationId string,
) string {
	sum := sha256.Sum256([]byte(conversationId))

	return fmt.Sprintf("%x", sum)
}
