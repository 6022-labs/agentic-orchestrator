package computes_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/compute"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenComputingCloneStrategyTestingSuite struct {
	sut *compute.CloneComputeStrategy

	mockAgentContextRequester                    *interfaces_mocks.MockAgentContextRequesterInterface
	mockAgentContextMessageRequester             *interfaces_mocks.MockAgentContextMessageRequesterInterface
	mockPromptParametersBuilderFactory           *interfaces_mocks.MockPromptParametersBuilderFactoryInterface
	mockPromptParametersBuilder                  *interfaces_mocks.MockPromptParametersBuilderInterface
	mockPromptFromPromptTemplateParser           *interfaces_mocks.MockPromptFromPromptTemplateParserInterface
	mockLLMConnectionChatCompletionRequester     *interfaces_mocks.MockLLMConnectionChatCompletionRequesterInterface
	mockAgentContextConversationMappingRequester *interfaces_mocks.MockAgentContextConversationMappingRequesterInterface
}

func WhenComputingCloneStrategyBeforeEach(t *testing.T) *WhenComputingCloneStrategyTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentContextRequester := interfaces_mocks.NewMockAgentContextRequesterInterface(mockController)
	mockAgentContextMessageRequester := interfaces_mocks.NewMockAgentContextMessageRequesterInterface(mockController)
	mockPromptParametersBuilderFactory := interfaces_mocks.NewMockPromptParametersBuilderFactoryInterface(mockController)
	mockPromptParametersBuilder := interfaces_mocks.NewMockPromptParametersBuilderInterface(mockController)
	mockPromptFromPromptTemplateParser := interfaces_mocks.NewMockPromptFromPromptTemplateParserInterface(mockController)
	mockLLMConnectionChatCompletionRequester := interfaces_mocks.NewMockLLMConnectionChatCompletionRequesterInterface(mockController)
	mockAgentContextConversationMappingRequester := interfaces_mocks.NewMockAgentContextConversationMappingRequesterInterface(mockController)

	logger := zap.NewNop()

	sut := compute.NewCloneComputeStrategy(
		logger,
		mockAgentContextRequester,
		mockAgentContextMessageRequester,
		mockPromptParametersBuilderFactory,
		mockPromptFromPromptTemplateParser,
		mockLLMConnectionChatCompletionRequester,
		mockAgentContextConversationMappingRequester,
	)

	return &WhenComputingCloneStrategyTestingSuite{
		sut:                                          sut,
		mockAgentContextRequester:                    mockAgentContextRequester,
		mockAgentContextMessageRequester:             mockAgentContextMessageRequester,
		mockPromptParametersBuilderFactory:           mockPromptParametersBuilderFactory,
		mockPromptParametersBuilder:                  mockPromptParametersBuilder,
		mockPromptFromPromptTemplateParser:           mockPromptFromPromptTemplateParser,
		mockLLMConnectionChatCompletionRequester:     mockLLMConnectionChatCompletionRequester,
		mockAgentContextConversationMappingRequester: mockAgentContextConversationMappingRequester,
	}
}

func (suite *WhenComputingCloneStrategyTestingSuite) setupPromptParametersBuilder() {
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

func TestWhenComputingCloneStrategy(t *testing.T) {
	t.Parallel()

	agentId := uuid.New()
	agentContextId := uuid.New()
	contextKey := "test-context-key"
	conversationId := "conv-1"

	agent := &models.Agent{
		ID:        agentId,
		EnsDomain: "clone.protocol6022.eth",
		Role:      enums.AgentRoleClone,
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

	t.Run("Given agent context already exists and conversation is already mapped", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are a clone agent"
		defaultPrompt := "Respond to the user"
		defaultCompletion := "Hello, I am a clone"

		initSuite := func(suite *WhenComputingCloneStrategyTestingSuite) {
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(agentContext, &isFound, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&defaultCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should get or create agent context and map conversation if not mapped", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingCloneStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Hello clone",
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
			assert.Equal(t, []string{defaultCompletion}, result)
		})
	})

	t.Run("Given agent context exists but conversation is not mapped", func(t *testing.T) {
		t.Parallel()

		isFound := true
		isNotFound := false
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are a clone agent"
		defaultPrompt := "Respond to the user"
		defaultCompletion := "Hello, I am a clone"

		t.Run("Should link the conversation to the agent context", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingCloneStrategyBeforeEach(t)

			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(nil, &isNotFound, nil)

			suite.mockAgentContextConversationMappingRequester.EXPECT().
				LinkConversationToAgentContext(agentId, conversationId, contextKey).
				Return(nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&defaultCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Hello clone",
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
			assert.Equal(t, []string{defaultCompletion}, result)
		})
	})

	t.Run("Given a new conversation with no system prompt", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "You are a clone system prompt"
		defaultPrompt := "Default prompt for clone"
		defaultCompletion := "Clone response to user"

		t.Run("Should call prompt template parser for system prompt", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingCloneStrategyBeforeEach(t)

			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(agentContext, &isFound, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&defaultCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Hello clone",
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
		})
	})

	t.Run("Given LLM returns a completion for the default prompt", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}
		systemPrompt := "System prompt"
		defaultPrompt := "Default prompt"
		defaultCompletion := "This is the LLM response"

		t.Run("Should save messages and return the completion", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingCloneStrategyBeforeEach(t)

			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.mockAgentContextConversationMappingRequester.EXPECT().
				GetLinkedAgentContextKeyOfConversation(agentId, conversationId).
				Return(agentContext, &isFound, nil)

			suite.setupPromptParametersBuilder()

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneSystem(agentId, gomock.Any()).
				Return(&systemPrompt, nil)

			suite.mockPromptFromPromptTemplateParser.EXPECT().
				ParseCloneDefault(agentId, gomock.Any()).
				Return(&defaultPrompt, nil)

			suite.mockLLMConnectionChatCompletionRequester.EXPECT().
				ChatCompletion(agentId, gomock.Any()).
				Return(&defaultCompletion, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Tell me something",
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
			assert.Equal(t, "This is the LLM response", result[0])
		})
	})
}
