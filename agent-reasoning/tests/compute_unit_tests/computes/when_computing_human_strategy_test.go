package computes_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/compute"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/compute_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenComputingHumanStrategyTestingSuite struct {
	sut *compute.HumanComputeStrategy

	mockAgentContextRequester        *interfaces_mocks.MockAgentContextRequesterInterface
	mockAgentContextMessageRequester *interfaces_mocks.MockAgentContextMessageRequesterInterface
}

func WhenComputingHumanStrategyBeforeEach(t *testing.T) *WhenComputingHumanStrategyTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentContextRequester := interfaces_mocks.NewMockAgentContextRequesterInterface(mockController)
	mockAgentContextMessageRequester := interfaces_mocks.NewMockAgentContextMessageRequesterInterface(mockController)

	logger := zap.NewNop()

	sut := compute.NewHumanComputeStrategy(
		logger,
		mockAgentContextRequester,
		mockAgentContextMessageRequester,
	)

	return &WhenComputingHumanStrategyTestingSuite{
		sut:                              sut,
		mockAgentContextRequester:        mockAgentContextRequester,
		mockAgentContextMessageRequester: mockAgentContextMessageRequester,
	}
}

func TestWhenComputingHumanStrategy(t *testing.T) {
	t.Parallel()

	agentId := uuid.New()
	agentContextId := uuid.New()
	contextKey := "test-context-key"
	conversationId := "conv-1"

	agent := &models.Agent{
		ID:        agentId,
		EnsDomain: "human.protocol6022.eth",
		Role:      enums.AgentRoleHuman,
	}

	computeContext := &compute.ComputeContext{
		Platform: enums.AgentContextMessagePlatformSlack,
		ConversationContext: &compute.ConversationContext{
			Id: conversationId,
			Agents: []*compute.ConversationContextAgent{
				{CollectionAddress: "0x123", CollectionTokenId: 1},
			},
		},
	}

	t.Run("Given agent context already exists", func(t *testing.T) {
		t.Parallel()

		isFound := true
		agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}

		initSuite := func(suite *WhenComputingHumanStrategyTestingSuite) {
			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should get or create agent context by key derived from conversation agents", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingHumanStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Hello from human",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.Nil(t, result)
		})

		t.Run("Should store all messages (union of input + stored) to context", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingHumanStrategyBeforeEach(t)

			storedMessage := &models.AgentContextMessage{
				Content:        "Previously stored message",
				Author:         "other.protocol6022.eth",
				Role:           enums.AgentContextMessageRoleUser,
				Platform:       enums.AgentContextMessagePlatformSlack,
				ConversationId: &conversationId,
				CreatedAt:      time.Now().Add(-time.Hour),
			}

			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{storedMessage}, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "New message from human",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.Nil(t, result)
		})

		t.Run("Should return nil completions because human agents do not produce completions", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingHumanStrategyBeforeEach(t)
			initSuite(suite)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "Human message",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given agent context does not exist yet", func(t *testing.T) {
		t.Parallel()

		t.Run("Should create agent context and proceed", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingHumanStrategyBeforeEach(t)

			isNotFound := false
			isFound := true
			agentContext := &models.AgentContext{Id: agentContextId, Key: contextKey}

			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(nil, &isNotFound, nil)

			suite.mockAgentContextRequester.EXPECT().
				CreateAgentContext(agentId, gomock.Any()).
				Return(nil)

			suite.mockAgentContextRequester.EXPECT().
				GetAgentContextByKey(agentId, gomock.Any()).
				Return(agentContext, &isFound, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				GetAgentContextMessages(agentId, contextKey).
				Return([]*models.AgentContextMessage{}, nil)

			suite.mockAgentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)

			inputs := []*models.AgentContextMessage{
				{
					Content:        "First human message",
					Author:         "user.protocol6022.eth",
					Role:           enums.AgentContextMessageRoleUser,
					Platform:       enums.AgentContextMessagePlatformSlack,
					ConversationId: &conversationId,
					CreatedAt:      time.Now(),
				},
			}

			result, err := suite.sut.Compute(agent, inputs, computeContext)

			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	})
}
