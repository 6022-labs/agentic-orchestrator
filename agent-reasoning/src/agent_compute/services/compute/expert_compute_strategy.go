package compute

import (
	"fmt"
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	LegacyContextDiscoveryKey               = "threading"
	ConversationContextDiscoveryKeyPrefix   = "conversation-context-discovery"
	ConversationContextDiscoveryKeyTemplate = "conversation-context-discovery/%s"
)

type ExpertComputeStrategyInterface interface {
	Compute(agent *models.Agent, inputs []*models.AgentContextMessage, context *ComputeContext) ([]string, error)
}

type ExpertComputeStrategy struct {
	*AbstractComputeStrategy
	agentContextRequester                    interfaces.AgentContextRequesterInterface
	agentPromptTemplateRequester             interfaces.AgentPromptTemplateRequesterInterface
	promptFromPromptTemplateParser           interfaces.PromptFromPromptTemplateParserInterface
	promptParametersBuilderFactory           interfaces.PromptParametersBuilderFactoryInterface
	llmConnectionChatCompletionRequester     interfaces.LLMConnectionChatCompletionRequesterInterface
	agentContextConversationMappingRequester interfaces.AgentContextConversationMappingRequesterInterface
}

func NewExpertComputeStrategy(
	logger *zap.Logger,
	agentContextRequester interfaces.AgentContextRequesterInterface,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface,
	promptParametersBuilderFactory interfaces.PromptParametersBuilderFactoryInterface,
	promptFromPromptTemplateParser interfaces.PromptFromPromptTemplateParserInterface,
	llmConnectionChatCompletionRequester interfaces.LLMConnectionChatCompletionRequesterInterface,
	agentContextConversationMappingRequester interfaces.AgentContextConversationMappingRequesterInterface,
) *ExpertComputeStrategy {
	return &ExpertComputeStrategy{
		AbstractComputeStrategy: NewAbstractComputeStrategy(
			logger,
			agentContextMessageRequester,
		),
		agentContextRequester:                    agentContextRequester,
		agentPromptTemplateRequester:             agentPromptTemplateRequester,
		promptFromPromptTemplateParser:           promptFromPromptTemplateParser,
		promptParametersBuilderFactory:           promptParametersBuilderFactory,
		llmConnectionChatCompletionRequester:     llmConnectionChatCompletionRequester,
		agentContextConversationMappingRequester: agentContextConversationMappingRequester,
	}
}

func (cs *ExpertComputeStrategy) Compute(agent *models.Agent, inputs []*models.AgentContextMessage, computeContext *ComputeContext) ([]string, error) {
	cs.logger.Debug(
		"Start to compute expert agent context completion",
		zap.Any("inputs", inputs),
		zap.Any("context", computeContext),
		zap.String("agentId", agent.ID.String()),
		zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
	)

	agentContext, isFound, err := cs.agentContextConversationMappingRequester.GetLinkedAgentContextKeyOfConversation(agent.ID, computeContext.ConversationContext.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked agent context key of conversation: %w", err)
	}

	latestMessage := inputs[len(inputs)-1]
	initialRequest := inputs[0]

	var responses []string
	if !*isFound {
		cs.logger.Debug(
			"Agent context not found for this conversation, trying to find it",
			zap.String("agentId", agent.ID.String()),
			zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
		)

		var threadNotFoundCompletion *string
		threadNotFoundCompletion, agentContext, err = cs.tryFindConversationContext(agent, computeContext, inputs, latestMessage, initialRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to try find agent context key: %w", err)
		}
		responses = append(responses, *threadNotFoundCompletion)
		if agentContext == nil {
			cs.logger.Debug(
				"Agent context not found for this conversation, no agent context key found",
				zap.String("agentId", agent.ID.String()),
				zap.String("conversationId", computeContext.ConversationContext.Id),
			)
			return responses, nil
		}

		err = cs.agentContextConversationMappingRequester.LinkConversationToAgentContext(agent.ID, computeContext.ConversationContext.Id, agentContext.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to link conversation to agent context: %w", err)
		}
	}

	cs.logger.Debug(
		"Agent context found for this conversation",
		zap.String("agentId", agent.ID.String()),
		zap.String("agentContextKey", agentContext.Key),
		zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
	)

	storedMessages, err := cs.agentContextMessageRequester.GetAgentContextMessages(agent.ID, agentContext.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context messages: %w", err)
	}

	messages := cs.unionMessagesByHash(agentContext.Id, inputs, storedMessages)
	messages = cs.groupMessages(messages, computeContext.ConversationContext.Id)

	if !cs.isThereSystemPromptForThisConversation(messages, computeContext.ConversationContext.Id) {
		messages, err = cs.addSystemPrompt(agent, computeContext, messages, latestMessage, initialRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to add system prompt: %w", err)
		}
	}

	var threadFoundCompletion *string
	threadFoundCompletion, messages, err = cs.processContextFoundPrompt(agent, computeContext, agentContext, messages, latestMessage, initialRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to process thread found prompt: %w", err)
	}

	nonStoredMessages := cs.diffMessagesByHash(agentContext.Id, messages, storedMessages)

	err = cs.saveMessages(agent, agentContext.Key, nonStoredMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to save non-stored messages: %w", err)
	}

	responses = append(responses, *threadFoundCompletion)
	return responses, nil
}

func (cs *ExpertComputeStrategy) getAgentContexts(
	agent *models.Agent,
) ([]*models.AgentContext, error) {
	agentContexts, err := cs.agentContextRequester.GetAgentContexts(agent.ID, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent context keys: %w", err)
	}
	if len(agentContexts) == 0 {
		return nil, fmt.Errorf("no agent context keys found")
	}

	return agentContexts, nil
}

func (cs *ExpertComputeStrategy) processContextNotFoundPrompt(
	agent *models.Agent,
	keys []string,
	computeContext *ComputeContext,
	messages []*models.AgentContextMessage,
	latestRequest *models.AgentContextMessage,
	initialRequest *models.AgentContextMessage,
) (*string, []*models.AgentContextMessage, error) {
	keysWithoutContextDiscoveryKeys := lo.Filter(keys, func(key string, _ int) bool {
		return !strings.HasPrefix(key, ConversationContextDiscoveryKeyPrefix)
	})

	promptParameterBuilder := cs.promptParametersBuilderFactory.Create()
	promptParameters, err := promptParameterBuilder.
		WithBase(computeContext.PromptParameters).
		WithParam([]string{"ContextKeys"}, keysWithoutContextDiscoveryKeys, nil).
		WithParam([]string{"LatestRequest", "Content"}, latestRequest.Content, nil).
		WithParam([]string{"InitialRequest", "Content"}, initialRequest.Content, nil).
		WithParam([]string{"ConversationId"}, computeContext.ConversationContext.Id, nil).
		WithParam([]string{"LatestRequest", "Author"}, latestRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"InitialRequest", "Author"}, initialRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build prompt parameters: %w", err)
	}

	contextNotFoundPrompt, err := cs.promptFromPromptTemplateParser.ParseExpertContextNotFound(agent.ID, promptParameters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse expert context not found prompt: %w", err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Content:        *contextNotFoundPrompt,
		Author:         "",
		Role:           enums.AgentContextMessageRoleUser,
		Platform:       enums.AgentContextMessagePlatformNone,
		ConversationId: &computeContext.ConversationContext.Id,
	})

	threadNotFoundCompletion, err := cs.llmConnectionChatCompletionRequester.ChatCompletion(agent.ID, messages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get completion for prompt template key %v: %w", enums.AgentPromptTemplateKeyExpertContextNotFound, err)
	}

	return threadNotFoundCompletion, messages, nil
}

func (cs *ExpertComputeStrategy) addSystemPrompt(
	agent *models.Agent,
	computeContext *ComputeContext,
	messages []*models.AgentContextMessage,
	latestRequest *models.AgentContextMessage,
	initialRequest *models.AgentContextMessage,
) ([]*models.AgentContextMessage, error) {

	promptParameterBuilder := cs.promptParametersBuilderFactory.Create()
	promptParameters, err := promptParameterBuilder.
		WithBase(computeContext.PromptParameters).
		WithParam([]string{"LatestRequest", "Content"}, latestRequest.Content, nil).
		WithParam([]string{"InitialRequest", "Content"}, initialRequest.Content, nil).
		WithParam([]string{"ConversationId"}, computeContext.ConversationContext.Id, nil).
		WithParam([]string{"LatestRequest", "Author"}, latestRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"InitialRequest", "Author"}, initialRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt parameters: %w", err)
	}

	systemPrompt, err := cs.promptFromPromptTemplateParser.ParseExpertSystem(agent.ID, promptParameters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expert system prompt: %w", err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Content:        *systemPrompt,
		Author:         "",
		Role:           enums.AgentContextMessageRoleSystem,
		Platform:       enums.AgentContextMessagePlatformNone,
		ConversationId: &computeContext.ConversationContext.Id,
	})

	return messages, nil
}

func (cs *ExpertComputeStrategy) processContextFoundPrompt(
	agent *models.Agent,
	computeContext *ComputeContext,
	agentContext *models.AgentContext,
	messages []*models.AgentContextMessage,
	latestRequest *models.AgentContextMessage,
	initialRequest *models.AgentContextMessage,
) (*string, []*models.AgentContextMessage, error) {
	promptParameterBuilder := cs.promptParametersBuilderFactory.Create()
	promptParameters, err := promptParameterBuilder.
		WithBase(computeContext.PromptParameters).
		WithParam([]string{"ContextKey"}, agentContext.Key, nil).
		WithParam([]string{"LatestRequest", "Content"}, latestRequest.Content, nil).
		WithParam([]string{"InitialRequest", "Content"}, initialRequest.Content, nil).
		WithParam([]string{"ConversationId"}, computeContext.ConversationContext.Id, nil).
		WithParam([]string{"LatestRequest", "Author"}, latestRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"InitialRequest", "Author"}, initialRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build prompt parameters: %w", err)
	}

	contextFoundPrompt, err := cs.promptFromPromptTemplateParser.ParseExpertContextFound(agent.ID, promptParameters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse expert context found prompt: %w", err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Content:        *contextFoundPrompt,
		Author:         "",
		ConversationId: &computeContext.ConversationContext.Id,
		Role:           enums.AgentContextMessageRoleUser,
		Platform:       enums.AgentContextMessagePlatformNone,
	})

	contextFoundCompletion, err := cs.llmConnectionChatCompletionRequester.ChatCompletion(agent.ID, messages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get completion for prompt template key %v: %w", enums.AgentPromptTemplateKeyExpertContextFound, err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Author:         agent.EnsDomain,
		Content:        *contextFoundCompletion,
		ConversationId: &computeContext.ConversationContext.Id,
		Platform:       enums.AgentContextMessagePlatformNone,
		Role:           enums.AgentContextMessageRoleAssistant,
	})

	return contextFoundCompletion, messages, nil
}

func (cs *ExpertComputeStrategy) tryFindConversationContext(
	agent *models.Agent,
	computeContext *ComputeContext,
	inputs []*models.AgentContextMessage,
	latestRequest *models.AgentContextMessage,
	initialRequest *models.AgentContextMessage,
) (*string, *models.AgentContext, error) {
	agentContexts, err := cs.getAgentContexts(agent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent contexts: %w", err)
	}

	keys := lo.Map(agentContexts, func(agentContext *models.AgentContext, _ int) string {
		return agentContext.Key
	})

	contextDiscoveryKey := fmt.Sprintf(ConversationContextDiscoveryKeyTemplate, computeContext.ConversationContext.Id)

	if !lo.Contains(keys, contextDiscoveryKey) {
		cs.logger.Info("Discovery context not found, creating it",
			zap.String("agentId", agent.ID.String()),
			zap.String("contextDiscoveryKey", contextDiscoveryKey),
			zap.String("conversationId", computeContext.ConversationContext.Id),
		)

		err := cs.agentContextRequester.CreateAgentContext(agent.ID, contextDiscoveryKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create agent context: %w", err)
		}

		agentContexts, err = cs.getAgentContexts(agent)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get agent contexts after creation: %w", err)
		}
	}

	contextDiscovery, found := lo.Find(agentContexts, func(agentContext *models.AgentContext) bool {
		return agentContext.Key == contextDiscoveryKey
	})
	if !found {
		return nil, nil, fmt.Errorf("context discovery not found")
	}

	storedMessages, err := cs.agentContextMessageRequester.GetAgentContextMessages(agent.ID, contextDiscovery.Key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent context messages: %w", err)
	}

	messages := cs.unionMessagesByHash(contextDiscovery.Id, inputs, storedMessages)
	messages = cs.groupMessages(messages, computeContext.ConversationContext.Id)

	var contextNotFoundCompletion *string
	contextNotFoundCompletion, messages, err = cs.processContextNotFoundPrompt(agent, keys, computeContext, messages, latestRequest, initialRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process context not found prompt: %w", err)
	}

	cs.logger.Debug(
		"Context not found completion",
		zap.String("agentId", agent.ID.String()),
		zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
		zap.String("contextNotFoundCompletion", *contextNotFoundCompletion),
	)

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Author:         agent.EnsDomain,
		Content:        *contextNotFoundCompletion,
		ConversationId: &computeContext.ConversationContext.Id,
		Platform:       enums.AgentContextMessagePlatformNone,
		Role:           enums.AgentContextMessageRoleAssistant,
	})

	foundedAgentContext, _ := lo.Find(agentContexts, func(agentContext *models.AgentContext) bool {
		return strings.Contains(*contextNotFoundCompletion, agentContext.Key)
	})

	nonStoredMessages := cs.diffMessagesByHash(contextDiscovery.Id, messages, storedMessages)
	err = cs.saveMessages(agent, contextDiscovery.Key, nonStoredMessages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save non-stored messages: %w", err)
	}

	return contextNotFoundCompletion, foundedAgentContext, nil
}
