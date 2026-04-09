package compute

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CloneComputeStrategyInterface interface {
	Compute(agent *models.Agent, inputs []*models.AgentContextMessage, context *ComputeContext) ([]string, error)
}

type CloneComputeStrategy struct {
	*AbstractAgentSetScopedComputeStrategy
	promptParametersBuilderFactory           interfaces.PromptParametersBuilderFactoryInterface
	promptFromPromptTemplateParser           interfaces.PromptFromPromptTemplateParserInterface
	llmConnectionChatCompletionRequester     interfaces.LLMConnectionChatCompletionRequesterInterface
	agentContextConversationMappingRequester interfaces.AgentContextConversationMappingRequesterInterface
}

func NewCloneComputeStrategy(
	logger *zap.Logger,
	agentContextRequester interfaces.AgentContextRequesterInterface,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
	promptParametersBuilderFactory interfaces.PromptParametersBuilderFactoryInterface,
	promptFromPromptTemplateParser interfaces.PromptFromPromptTemplateParserInterface,
	llmConnectionChatCompletionRequester interfaces.LLMConnectionChatCompletionRequesterInterface,
	agentContextConversationMappingRequester interfaces.AgentContextConversationMappingRequesterInterface,
) *CloneComputeStrategy {
	return &CloneComputeStrategy{
		AbstractAgentSetScopedComputeStrategy: NewAbstractAgentSetScopedComputeStrategy(
			logger,
			agentContextRequester,
			agentContextMessageRequester,
		),
		promptParametersBuilderFactory:           promptParametersBuilderFactory,
		promptFromPromptTemplateParser:           promptFromPromptTemplateParser,
		llmConnectionChatCompletionRequester:     llmConnectionChatCompletionRequester,
		agentContextConversationMappingRequester: agentContextConversationMappingRequester,
	}
}

func (cs *CloneComputeStrategy) Compute(
	agent *models.Agent,
	inputs []*models.AgentContextMessage,
	computeContext *ComputeContext,
) ([]string, error) {
	cs.logger.Debug("Computing clone strategy",
		zap.String("agentId", agent.ID.String()),
		zap.String("latestMessageAuthor", inputs[len(inputs)-1].Author),
		zap.String("conversationContext.Id", computeContext.ConversationContext.Id),
	)

	key := cs.getContextKey(computeContext.ConversationContext.Agents)

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

	err = cs.mapContextToConversationIfNotMapped(agent.ID, computeContext.ConversationContext.Id, agentContext.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to map context to conversation: %w", err)
	}

	latestMessage := inputs[len(inputs)-1]
	initialRequest := inputs[0]

	var messages []*models.AgentContextMessage

	if !cs.isThereSystemPromptForThisConversation(messages, computeContext.ConversationContext.Id) {
		messages, err = cs.addSystemPrompt(agent, computeContext, storedMessages, latestMessage, initialRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to add system prompt: %w", err)
		}
	}

	messages = cs.unionMessagesByHash(agentContext.Id, inputs, messages)

	var defaultCompletion *string
	defaultCompletion, messages, err = cs.processDefaultPrompt(agent, computeContext, messages, latestMessage, initialRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to process default prompt: %w", err)
	}

	nonStoredMessages := cs.diffMessagesByHash(agentContext.Id, messages, storedMessages)

	err = cs.saveMessages(agent, agentContext.Key, nonStoredMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to save non-stored messages: %w", err)
	}

	return []string{*defaultCompletion}, nil
}

func (cs *CloneComputeStrategy) mapContextToConversationIfNotMapped(
	agentId uuid.UUID,
	conversationId string,
	agentContextKey string,
) error {
	cs.logger.Debug("Mapping context to conversation",
		zap.String("agentId", agentId.String()),
		zap.String("conversationContext.Id", conversationId),
	)

	_, found, err := cs.agentContextConversationMappingRequester.GetLinkedAgentContextKeyOfConversation(agentId, conversationId)
	if err != nil {
		return fmt.Errorf("failed to get linked agent context key of conversation: %w", err)
	}

	if !*found {
		err = cs.agentContextConversationMappingRequester.LinkConversationToAgentContext(agentId, conversationId, agentContextKey)
		if err != nil {
			return fmt.Errorf("failed to link agent context key to conversation: %w", err)
		}
	}

	return nil
}

func (cs *CloneComputeStrategy) addSystemPrompt(
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

func (cs *CloneComputeStrategy) processDefaultPrompt(
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
		WithBase(computeContext.PromptParameters).
		WithParam([]string{"LatestRequest", "Content"}, latestRequest.Content, nil).
		WithParam([]string{"InitialRequest", "Content"}, initialRequest.Content, nil).
		WithParam([]string{"ConversationId"}, computeContext.ConversationContext.Id, nil).
		WithParam([]string{"LatestRequest", "Author"}, latestRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		WithParam([]string{"InitialRequest", "Author"}, initialRequest.Author, interfaces.PromptParameterTypePotentialAgent.ToPointer()).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build prompt parameters: %w", err)
	}

	defaultPrompt, err := cs.promptFromPromptTemplateParser.ParseCloneDefault(agent.ID, promptParameters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute prompt template: %w", err)
	}

	messages = append(messages, &models.AgentContextMessage{
		CreatedAt:      time.Now(),
		Content:        *defaultPrompt,
		Author:         "",
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
