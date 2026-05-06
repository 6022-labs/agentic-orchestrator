package computes_test

import (
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

type WhenComputingFacilitatorStrategyTestingSuite struct {
	sut *compute.FacilitatorComputeStrategy

	mockAgentContextRequester                *interfaces_mocks.MockAgentContextRequesterInterface
	mockAgentContextMessageRequester         *interfaces_mocks.MockAgentContextMessageRequesterInterface
	mockAgentPromptTemplateRequester         *interfaces_mocks.MockAgentPromptTemplateRequesterInterface
	mockPromptParametersBuilderFactory       *interfaces_mocks.MockPromptParametersBuilderFactoryInterface
	mockPromptParametersBuilder              *interfaces_mocks.MockPromptParametersBuilderInterface
	mockPromptFromPromptTemplateParser       *interfaces_mocks.MockPromptFromPromptTemplateParserInterface
	mockLLMConnectionChatCompletionRequester *interfaces_mocks.MockLLMConnectionChatCompletionRequesterInterface
}

func WhenComputingFacilitatorStrategyBeforeEach(t *testing.T) *WhenComputingFacilitatorStrategyTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentContextRequester := interfaces_mocks.NewMockAgentContextRequesterInterface(mockController)
	mockAgentContextMessageRequester := interfaces_mocks.NewMockAgentContextMessageRequesterInterface(mockController)
	mockAgentPromptTemplateRequester := interfaces_mocks.NewMockAgentPromptTemplateRequesterInterface(mockController)
	mockPromptParametersBuilderFactory := interfaces_mocks.NewMockPromptParametersBuilderFactoryInterface(mockController)
	mockPromptParametersBuilder := interfaces_mocks.NewMockPromptParametersBuilderInterface(mockController)
	mockPromptFromPromptTemplateParser := interfaces_mocks.NewMockPromptFromPromptTemplateParserInterface(mockController)
	mockLLMConnectionChatCompletionRequester := interfaces_mocks.NewMockLLMConnectionChatCompletionRequesterInterface(mockController)

	logger := zap.NewNop()

	sut := compute.NewFacilitatorComputeStrategy(
		logger,
		mockAgentContextRequester,
		mockAgentContextMessageRequester,
		mockAgentPromptTemplateRequester,
		mockPromptParametersBuilderFactory,
		mockPromptFromPromptTemplateParser,
		mockLLMConnectionChatCompletionRequester,
	)

	return &WhenComputingFacilitatorStrategyTestingSuite{
		sut:                                      sut,
		mockAgentContextRequester:                mockAgentContextRequester,
		mockAgentContextMessageRequester:         mockAgentContextMessageRequester,
		mockAgentPromptTemplateRequester:         mockAgentPromptTemplateRequester,
		mockPromptParametersBuilderFactory:       mockPromptParametersBuilderFactory,
		mockPromptParametersBuilder:              mockPromptParametersBuilder,
		mockPromptFromPromptTemplateParser:       mockPromptFromPromptTemplateParser,
		mockLLMConnectionChatCompletionRequester: mockLLMConnectionChatCompletionRequester,
	}
}

func (suite *WhenComputingFacilitatorStrategyTestingSuite) setupPromptParametersBuilder() {
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

func TestWhenComputingFacilitatorStrategy(t *testing.T) {
	t.Parallel()

	agentId := uuid.New()
	agentContextId := uuid.New()
	contextKey := "test-context-key"
	conversationId := "conv-1"

	agent := &models.Agent{
		ID:        agentId,
		EnsDomain: "facilitator.protocol6022.eth",
		Role:      enums.AgentRoleFacilitator,
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

	t.Run("Given a standard conversation without end signal", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are a facilitator"
		defaultPrompt := "Facilitate the conversation"
		defaultCompletion := "Here is my facilitation response"

		initSuite := func(suite *WhenComputingFacilitatorStrategyTestingSuite) {
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&defaultCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should compute default prompt and return completion with isFinal=false", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingFacilitatorStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Please facilitate this",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, isFinal, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotNil(t, isFinal)
			assert.False(t, *isFinal)
			assert.Equal(t, []string{defaultCompletion}, result)
		})
	})

	t.Run("Given messages contain :end: signal", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are a facilitator"
		defaultPrompt := "Facilitate the conversation"
		defaultCompletion := "Here is my response"
		concludePrompt := "Conclude the conversation"
		concludeCompletion := "Here is the conclusion"

		// Pre-existing messages that include :end: signal
		// We need at least 2 messages with assistant/user role for isAgentReadyToConclude
		// Include a system prompt so no new system prompt is parsed
		storedSystemMessage := &models.AgentContextMessage{
			Content:        systemPrompt,
			Author:         "",
			Role:           enums.AgentContextMessageRoleSystem,
			Platform:       enums.AgentContextMessagePlatformNone,
			ConversationId: &conversationId,
			CreatedAt:      time.Now().Add(-3 * time.Hour),
		}
		storedAssistantMessage := &models.AgentContextMessage{
			Content:        "Previous assistant response",
			Author:         "facilitator.protocol6022.eth",
			Role:           enums.AgentContextMessageRoleAssistant,
			Platform:       enums.AgentContextMessagePlatformNone,
			ConversationId: &conversationId,
			CreatedAt:      time.Now().Add(-2 * time.Hour),
		}
		storedUserMessage := &models.AgentContextMessage{
			Content:        "User says :end:",
			Author:         "user.protocol6022.eth",
			Role:           enums.AgentContextMessageRoleUser,
			Platform:       enums.AgentContextMessagePlatformSlack,
			ConversationId: &conversationId,
			CreatedAt:      time.Now().Add(-time.Hour),
		}

		initSuite := func(suite *WhenComputingFacilitatorStrategyTestingSuite) {
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{storedSystemMessage, storedAssistantMessage, storedUserMessage}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorConclude(agentId, gomock.Any()).
				Return(&concludePrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&defaultCompletion, nil).
				Times(1)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&concludeCompletion, nil).
				Times(1)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should process conclude prompt and return isFinal=true", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingFacilitatorStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Continue :end:",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, isFinal, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotNil(t, isFinal)
			assert.True(t, *isFinal)
			assert.Len(t, result, 2)
			assert.Equal(t, defaultCompletion, result[0])
			assert.Equal(t, concludeCompletion, result[1])
		})
	})

	t.Run("Given the LLM emits unicode emoji in its completion", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are a facilitator"
		defaultPrompt := "Facilitate the conversation"
		unicodeCompletion := "🔄 status update"
		expectedDecoded := ":arrows_counterclockwise: status update"

		var capturedRequest *interfaces.AddMessageToAgentContextRequest

		initSuite := func(suite *WhenComputingFacilitatorStrategyTestingSuite) {
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

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

			suite := WhenComputingFacilitatorStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Please facilitate this",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, _, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.Equal(t, []string{expectedDecoded}, result)

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
		systemPrompt := "You are a facilitator"
		defaultPrompt := "Facilitate the conversation"
		previousAssistantContent := "status :fire: shipping"
		newDefaultCompletion := "Next turn"

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

		initSuite := func(suite *WhenComputingFacilitatorStrategyTestingSuite) {
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{storedSystem, storedAssistant}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&newDefaultCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				DoAndReturn(func(request *interfaces.AddMessageToAgentContextRequest) error {
					capturedRequest = request
					return nil
				})
		}

		t.Run("Should not re-save the replay when content+platform already in stored history", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingFacilitatorStrategyBeforeEach(t)
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

			_, _, err := suite.sut.Compute(agent, inputs, computeContext)
			assert.NoError(t, err)

			assert.NotNil(t, capturedRequest)
			savedReplays := lo.Filter(capturedRequest.Messages, func(m *interfaces.AddMessageToAgentContextRequestMessage, _ int) bool {
				return m.Role == enums.AgentContextMessageRoleAssistant && m.Content == previousAssistantContent
			})
			assert.Len(t, savedReplays, 0, "an assistant message already in stored history must not be re-saved when replayed as input")
		})
	})

	t.Run("Given ForceConclusion is true", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are a facilitator"
		defaultPrompt := "Facilitate"
		defaultCompletion := "Default response"
		concludePrompt := "Conclude"
		concludeCompletion := "Forced conclusion"

		forceConcludeContext := &compute.ComputeContext{
			Platform:         enums.AgentContextMessagePlatformSlack,
			ForceConclusion:  true,
			PromptParameters: map[string]any{},
			ConversationContext: &compute.ConversationContext{
				Id: conversationId,
				Agents: []*compute.ConversationContextAgent{
					{CollectionAddress: "0x123", CollectionTokenId: 1},
				},
			},
		}

		initSuite := func(suite *WhenComputingFacilitatorStrategyTestingSuite) {
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseFacilitatorConclude(agentId, gomock.Any()).
				Return(&concludePrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&defaultCompletion, nil).
				Times(1)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&concludeCompletion, nil).
				Times(1)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should process conclude prompt regardless of end signal", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingFacilitatorStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Normal message without end signal",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, isFinal, err := suite.sut.Compute(agent, inputs, forceConcludeContext)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotNil(t, isFinal)
			assert.True(t, *isFinal)
			assert.Len(t, result, 2)
			assert.Equal(t, defaultCompletion, result[0])
			assert.Equal(t, concludeCompletion, result[1])
		})
	})
}
