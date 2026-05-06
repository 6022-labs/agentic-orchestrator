package computes_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/compute"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenComputingExpertStrategyTestingSuite struct {
	sut *compute.ExpertComputeStrategy

	mockAgentContextRequester                    *interfaces_mocks.MockAgentContextRequesterInterface
	mockAgentContextMessageRequester             *interfaces_mocks.MockAgentContextMessageRequesterInterface
	mockAgentPromptTemplateRequester             *interfaces_mocks.MockAgentPromptTemplateRequesterInterface
	mockPromptParametersBuilderFactory           *interfaces_mocks.MockPromptParametersBuilderFactoryInterface
	mockPromptParametersBuilder                  *interfaces_mocks.MockPromptParametersBuilderInterface
	mockPromptFromPromptTemplateParser           *interfaces_mocks.MockPromptFromPromptTemplateParserInterface
	mockLLMConnectionChatCompletionRequester     *interfaces_mocks.MockLLMConnectionChatCompletionRequesterInterface
	mockAgentContextConversationMappingRequester *interfaces_mocks.MockAgentContextConversationMappingRequesterInterface
}

func WhenComputingExpertStrategyBeforeEach(t *testing.T) *WhenComputingExpertStrategyTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentContextRequester := interfaces_mocks.NewMockAgentContextRequesterInterface(mockController)
	mockAgentContextMessageRequester := interfaces_mocks.NewMockAgentContextMessageRequesterInterface(mockController)
	mockAgentPromptTemplateRequester := interfaces_mocks.NewMockAgentPromptTemplateRequesterInterface(mockController)
	mockPromptParametersBuilderFactory := interfaces_mocks.NewMockPromptParametersBuilderFactoryInterface(mockController)
	mockPromptParametersBuilder := interfaces_mocks.NewMockPromptParametersBuilderInterface(mockController)
	mockPromptFromPromptTemplateParser := interfaces_mocks.NewMockPromptFromPromptTemplateParserInterface(mockController)
	mockLLMConnectionChatCompletionRequester := interfaces_mocks.NewMockLLMConnectionChatCompletionRequesterInterface(mockController)
	mockAgentContextConversationMappingRequester := interfaces_mocks.NewMockAgentContextConversationMappingRequesterInterface(mockController)

	logger := zap.NewNop()

	sut := compute.NewExpertComputeStrategy(
		logger,
		mockAgentContextRequester,
		mockAgentContextMessageRequester,
		mockAgentPromptTemplateRequester,
		mockPromptParametersBuilderFactory,
		mockPromptFromPromptTemplateParser,
		mockLLMConnectionChatCompletionRequester,
		mockAgentContextConversationMappingRequester,
	)

	return &WhenComputingExpertStrategyTestingSuite{
		sut:                                          sut,
		mockAgentContextRequester:                    mockAgentContextRequester,
		mockAgentContextMessageRequester:             mockAgentContextMessageRequester,
		mockAgentPromptTemplateRequester:             mockAgentPromptTemplateRequester,
		mockPromptParametersBuilderFactory:           mockPromptParametersBuilderFactory,
		mockPromptParametersBuilder:                  mockPromptParametersBuilder,
		mockPromptFromPromptTemplateParser:           mockPromptFromPromptTemplateParser,
		mockLLMConnectionChatCompletionRequester:     mockLLMConnectionChatCompletionRequester,
		mockAgentContextConversationMappingRequester: mockAgentContextConversationMappingRequester,
	}
}

func (suite *WhenComputingExpertStrategyTestingSuite) setupPromptParametersBuilder() {
	suite.mockPromptParametersBuilderFactory.EXPECT().
		Create().
		Return(suite.mockPromptParametersBuilder).
		AnyTimes()

	suite.mockPromptParametersBuilder.EXPECT().
		WithBase(gomock.Any()).
		Return(suite.mockPromptParametersBuilder).
		AnyTimes()

	suite.mockPromptParametersBuilder.EXPECT().
		WithParam(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(suite.mockPromptParametersBuilder).
		AnyTimes()

	suite.mockPromptParametersBuilder.EXPECT().
		Build().
		Return(map[string]any{"test": "params"}, nil).
		AnyTimes()
}

func TestWhenComputingExpertStrategy(t *testing.T) {
	t.Parallel()

	agentId := uuid.New()
	agentContextId := uuid.New()
	contextKey := "existing-context-key"
	conversationId := "conv-1"
	contextDiscoveryKey := fmt.Sprintf("conversation-context-discovery/%s", conversationId)

	agent := &models.Agent{
		ID:        agentId,
		EnsDomain: "expert.protocol6022.eth",
		Role:      enums.AgentRoleExpert,
	}

	computeContext := &compute.ComputeContext{
		Platform:         enums.AgentContextMessagePlatformSlack,
		PromptParameters: map[string]any{},
		ConversationContext: &compute.ConversationContext{
			Id: conversationId,
			Agents: []*compute.ConversationContextAgent{
				{CollectionAddress: "0x123", CollectionTokenId: 1},
			},
		},
	}

	t.Run("Given conversation is already linked to a context", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are an expert"
		contextFoundPrompt := "Context found prompt"
		contextFoundCompletion := "Expert response from found context"

		initSuite := func(suite *WhenComputingExpertStrategyTestingSuite) {
			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertContextFound(agentId, gomock.Any()).
				Return(&contextFoundPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&contextFoundCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should use that context and process context-found prompt", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingExpertStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Ask the expert",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Contains(t, result, contextFoundCompletion)
		})
	})

	t.Run("Given conversation is NOT linked and context discovery finds a matching key", func(t *testing.T) {
		t.Parallel()

		isNotFound := false
		discoveryContextId := uuid.New()
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		discoveryContext := &models.AgentContext{Id: discoveryContextId, Key: contextDiscoveryKey}
		systemPrompt := "You are an expert"
		contextNotFoundPrompt := "Context not found prompt"
		// The LLM response contains the context key, which is how the strategy discovers it
		contextNotFoundCompletion := fmt.Sprintf("I found the context: %s", contextKey)
		contextFoundPrompt := "Context found prompt"
		contextFoundCompletion := "Expert response after discovery"

		initSuite := func(suite *WhenComputingExpertStrategyTestingSuite) {
			// Step 1: conversation is not linked
			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(nil, &isNotFound, nil)

			// Step 2: get all agent contexts for discovery
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContexts(agentId, nil, nil, nil).
				Return([]*models.AgentContext{agentContext, discoveryContext}, nil).
				Times(1)

			// Step 3: get stored messages for discovery context
			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextDiscoveryKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.setupPromptParametersBuilder()

			// Step 4: process context-not-found prompt (discovery)
			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertContextNotFound(agentId, gomock.Any()).
				Return(&contextNotFoundPrompt, nil)

			// First LLM call: discovery - returns completion containing the context key
			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&contextNotFoundCompletion, nil).
				Times(1)

			// Save discovery messages
			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)

			// Step 5: link conversation to discovered context
			suite.mockAgentContextConversationMappingRequester.EXPECT().
				LinkConversationToAgentContext(agentId, conversationId, contextKey).
				Return(nil)

			// Step 6: get stored messages for the found context
			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			// Step 7: add system prompt and process context-found prompt
			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertContextFound(agentId, gomock.Any()).
				Return(&contextFoundPrompt, nil)

			// Second LLM call: context-found completion
			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&contextFoundCompletion, nil).
				Times(1)

			// Save final messages
			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should discover context via LLM, link conversation, and process context-found prompt", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingExpertStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Ask the expert something",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, 2)
			assert.Equal(t, contextNotFoundCompletion, result[0])
			assert.Equal(t, contextFoundCompletion, result[1])
		})
	})

	t.Run("Given conversation is NOT linked and context discovery finds no matching key", func(t *testing.T) {
		t.Parallel()

		isNotFound := false
		discoveryContextId := uuid.New()
		discoveryContext := &models.AgentContext{Id: discoveryContextId, Key: contextDiscoveryKey}
		contextNotFoundPrompt := "Context not found prompt"
		// The LLM response does NOT contain any existing context key
		contextNotFoundCompletion := "I could not find a matching context for this conversation"

		initSuite := func(suite *WhenComputingExpertStrategyTestingSuite) {
			// Step 1: conversation is not linked
			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(nil, &isNotFound, nil)

			// Step 2: get all agent contexts - only the discovery context exists
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContexts(agentId, nil, nil, nil).
				Return([]*models.AgentContext{discoveryContext}, nil).
				Times(1)

			// Step 3: get stored messages for discovery context
			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextDiscoveryKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.setupPromptParametersBuilder()

			// Step 4: process context-not-found prompt
			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertContextNotFound(agentId, gomock.Any()).
				Return(&contextNotFoundPrompt, nil)

			// LLM call: returns completion without any context key
			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&contextNotFoundCompletion, nil)

			// Save discovery messages
			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should return only the discovery response", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingExpertStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Ask the expert something unrelated",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, 1)
			assert.Equal(t, contextNotFoundCompletion, result[0])
		})
	})

	t.Run("Given conversation is NOT linked and discovery context does not exist yet", func(t *testing.T) {
		t.Parallel()

		isNotFound := false
		discoveryContextId := uuid.New()
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		discoveryContext := &models.AgentContext{Id: discoveryContextId, Key: contextDiscoveryKey}
		contextNotFoundPrompt := "Context not found prompt"
		contextNotFoundCompletion := "No matching context found"

		initSuite := func(suite *WhenComputingExpertStrategyTestingSuite) {
			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(nil, &isNotFound, nil)

			// First call: contexts do not include discovery key
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContexts(agentId, nil, nil, nil).
				Return([]*models.AgentContext{agentContext}, nil).
				Times(1)

			// Create discovery context
			suite.mockAgentContextRequester.EXPECT().
				CreateAgentContext(agentId, contextDiscoveryKey).
				Return(nil)

			// Second call: now includes discovery context
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContexts(agentId, nil, nil, nil).
				Return([]*models.AgentContext{agentContext, discoveryContext}, nil).
				Times(1)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextDiscoveryKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertContextNotFound(agentId, gomock.Any()).
				Return(&contextNotFoundPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&contextNotFoundCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should create discovery context then proceed with discovery", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingExpertStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "First time asking expert",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, 1)
			assert.Equal(t, contextNotFoundCompletion, result[0])
		})
	})

	t.Run("Given the LLM emits unicode emoji in its completion", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are an expert"
		contextFoundPrompt := "Context found prompt"
		unicodeCompletion := "🔄 status update"
		expectedDecoded := ":arrows_counterclockwise: status update"

		storedSystem := &models.AgentContextMessage{
			Content:        systemPrompt,
			Author:         "",
			Role:           enums.AgentContextMessageRoleSystem,
			Platform:       enums.AgentContextMessagePlatformNone,
			ConversationId: &conversationId,
			CreatedAt:      time.Now().Add(-2 * time.Hour),
		}

		var capturedRequest *interfaces.AddMessageToAgentContextRequest

		initSuite := func(suite *WhenComputingExpertStrategyTestingSuite) {
			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{storedSystem}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertContextFound(agentId, gomock.Any()).
				Return(&contextFoundPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&unicodeCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				DoAndReturn(func(request *interfaces.AddMessageToAgentContextRequest) error {
					capturedRequest = request
					return nil
				})
		}

		t.Run("Should save the assistant content in shortcode form so it matches the conversation broker round-trip", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingExpertStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Ask the expert",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.Contains(t, result, expectedDecoded)

			assert.NotNil(t, capturedRequest)
			assistantMessages := lo.Filter(capturedRequest.Messages, func(m *interfaces.AddMessageToAgentContextRequestMessage, _ int) bool {
				return m.Role == enums.AgentContextMessageRoleAssistant
			})
			assert.Len(t, assistantMessages, 1)
			assert.Equal(t, expectedDecoded, assistantMessages[0].Content)
		})
	})

	t.Run("Given the assistant message replayed as input by the conversation broker matches a stored row", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are an expert"
		contextFoundPrompt := "Context found prompt"
		previousAssistantContent := "status :fire: shipping"
		newCompletion := "Next turn"

		storedSystem := &models.AgentContextMessage{
			Content:        systemPrompt,
			Author:         "",
			Role:           enums.AgentContextMessageRoleSystem,
			Platform:       enums.AgentContextMessagePlatformNone,
			ConversationId: &conversationId,
			CreatedAt:      time.Now().Add(-2 * time.Hour),
		}
		storedAssistant := &models.AgentContextMessage{
			Content:        previousAssistantContent,
			Author:         agent.EnsDomain,
			Role:           enums.AgentContextMessageRoleAssistant,
			Platform:       enums.AgentContextMessagePlatformSlack,
			ConversationId: &conversationId,
			CreatedAt:      time.Now().Add(-time.Hour),
		}

		var capturedRequest *interfaces.AddMessageToAgentContextRequest

		initSuite := func(suite *WhenComputingExpertStrategyTestingSuite) {
			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{storedSystem, storedAssistant}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseExpertContextFound(agentId, gomock.Any()).
				Return(&contextFoundPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&newCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				DoAndReturn(func(request *interfaces.AddMessageToAgentContextRequest) error {
					capturedRequest = request
					return nil
				})
		}

		t.Run("Should not re-save the replay when content+platform already in stored history", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingExpertStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        previousAssistantContent,
					Author:         agent.EnsDomain,
					Role:           enums.AgentContextMessageRoleAssistant,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now().Add(-30 * time.Minute),
				},
				{
					Content:        "follow-up question",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			_, err := suite.sut.Compute(agent, inputs, computeContext)
			assert.NoError(t, err)

			assert.NotNil(t, capturedRequest)
			savedReplays := lo.Filter(capturedRequest.Messages, func(m *interfaces.AddMessageToAgentContextRequestMessage, _ int) bool {
				return m.Role == enums.AgentContextMessageRoleAssistant && m.Content == previousAssistantContent
			})
			assert.Len(t, savedReplays, 0, "an assistant message already in stored history must not be re-saved when replayed as input")
		})
	})
}
