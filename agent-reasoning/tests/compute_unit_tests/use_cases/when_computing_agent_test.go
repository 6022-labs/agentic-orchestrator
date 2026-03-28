package use_cases_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/compute_mocks/services_mocks/compute_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/compute_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenComputingAgentTestingSuite struct {
	sut *use_cases.AgentCompute

	mockAgentRegistryRequester     *interfaces_mocks.MockAgentRegistryRequesterInterface
	mockCloneComputeStrategy       *compute_mocks.MockCloneComputeStrategyInterface
	mockHumanComputeStrategy       *compute_mocks.MockHumanComputeStrategyInterface
	mockExpertComputeStrategy      *compute_mocks.MockExpertComputeStrategyInterface
	mockFacilitatorComputeStrategy *compute_mocks.MockFacilitatorComputeStrategyInterface
}

func WhenComputingAgentBeforeEach(t *testing.T) *WhenComputingAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	mockCloneComputeStrategy := compute_mocks.NewMockCloneComputeStrategyInterface(mockController)
	mockHumanComputeStrategy := compute_mocks.NewMockHumanComputeStrategyInterface(mockController)
	mockExpertComputeStrategy := compute_mocks.NewMockExpertComputeStrategyInterface(mockController)
	mockFacilitatorComputeStrategy := compute_mocks.NewMockFacilitatorComputeStrategyInterface(mockController)

	logger := zap.NewNop()

	sut := use_cases.NewAgentCompute(
		logger,
		mockAgentRegistryRequester,
		mockCloneComputeStrategy,
		mockHumanComputeStrategy,
		mockExpertComputeStrategy,
		mockFacilitatorComputeStrategy,
	)

	return &WhenComputingAgentTestingSuite{
		sut:                            sut,
		mockAgentRegistryRequester:     mockAgentRegistryRequester,
		mockCloneComputeStrategy:       mockCloneComputeStrategy,
		mockHumanComputeStrategy:       mockHumanComputeStrategy,
		mockExpertComputeStrategy:      mockExpertComputeStrategy,
		mockFacilitatorComputeStrategy: mockFacilitatorComputeStrategy,
	}
}

func generateValidAgentComputeRequest() *requests.AgentComputeRequest {
	return &requests.AgentComputeRequest{
		Platform:               enums.AgentContextMessagePlatformSlack,
		AgentCollectionAddress: "0x1234567890abcdef",
		AgentCollectionTokenId: 42,
		Inputs: []*requests.AgentComputeInputRequest{
			{
				Author:    "requester.protocol6022.eth",
				Content:   "Hello, how are you?",
				Timestamp: time.Now(),
			},
		},
		ConversationContext: &requests.AgentComputeConversationContextRequest{
			Id: "conversation-123",
			Agents: []*requests.AgentComputeConversationContextAgentRequest{
				{
					CollectionAddress: "0x1234567890abcdef",
					CollectionTokenId: 42,
				},
			},
		},
		PromptParameters: map[string]any{},
	}
}

func TestWhenComputingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because inputs are empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)

			request := generateValidAgentComputeRequest()
			request.Inputs = []*requests.AgentComputeInputRequest{}

			result, err := suite.sut.Compute(request)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "inputs")
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenComputingAgentTestingSuite, request *requests.AgentComputeRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)
			request := generateValidAgentComputeRequest()
			initSuite(suite, request)

			result, err := suite.sut.Compute(request)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "agent with collection token ID 42 not found")
		})
	})

	t.Run("Given agent registry returns error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenComputingAgentTestingSuite, request *requests.AgentComputeRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)
			request := generateValidAgentComputeRequest()
			initSuite(suite, request)

			result, err := suite.sut.Compute(request)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given agent is an expert", func(t *testing.T) {
		t.Parallel()

		agentResponse := model_builders.NewAgentResponseBuilder().
			WithRole(enums.AgentRoleExpert).
			WithEnsDomain("expert.protocol6022.eth").
			Build()

		initSuite := func(suite *WhenComputingAgentTestingSuite, request *requests.AgentComputeRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agentResponse, nil)

			suite.mockExpertComputeStrategy.EXPECT().
				Compute(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]string{"expert response"}, nil)
		}

		t.Run("Should call expert strategy and return responses", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)
			request := generateValidAgentComputeRequest()
			initSuite(suite, request)

			result, err := suite.sut.Compute(request)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsFinal)
			assert.Equal(t, []string{"expert response"}, result.Responses)
		})
	})

	t.Run("Given agent is a human", func(t *testing.T) {
		t.Parallel()

		agentResponse := model_builders.NewAgentResponseBuilder().
			WithRole(enums.AgentRoleHuman).
			WithEnsDomain("human.protocol6022.eth").
			Build()

		initSuite := func(suite *WhenComputingAgentTestingSuite, request *requests.AgentComputeRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agentResponse, nil)

			suite.mockHumanComputeStrategy.EXPECT().
				Compute(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should call human strategy", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)
			request := generateValidAgentComputeRequest()
			initSuite(suite, request)

			result, err := suite.sut.Compute(request)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsFinal)
		})
	})

	t.Run("Given agent is a facilitator", func(t *testing.T) {
		t.Parallel()

		agentResponse := model_builders.NewAgentResponseBuilder().
			WithRole(enums.AgentRoleFacilitator).
			WithEnsDomain("facilitator.protocol6022.eth").
			Build()

		isFinal := true

		initSuite := func(suite *WhenComputingAgentTestingSuite, request *requests.AgentComputeRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agentResponse, nil)

			suite.mockFacilitatorComputeStrategy.EXPECT().
				Compute(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]string{"facilitator response", "conclude response"}, &isFinal, nil)
		}

		t.Run("Should call facilitator strategy and return isFinal", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)
			request := generateValidAgentComputeRequest()
			initSuite(suite, request)

			result, err := suite.sut.Compute(request)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, result.IsFinal)
			assert.Equal(t, []string{"facilitator response", "conclude response"}, result.Responses)
		})
	})

	t.Run("Given agent is a clone", func(t *testing.T) {
		t.Parallel()

		agentResponse := model_builders.NewAgentResponseBuilder().
			WithRole(enums.AgentRoleClone).
			WithEnsDomain("clone.protocol6022.eth").
			Build()

		initSuite := func(suite *WhenComputingAgentTestingSuite, request *requests.AgentComputeRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agentResponse, nil)

			suite.mockCloneComputeStrategy.EXPECT().
				Compute(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]string{"clone response"}, nil)
		}

		t.Run("Should call clone strategy", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)
			request := generateValidAgentComputeRequest()
			initSuite(suite, request)

			result, err := suite.sut.Compute(request)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsFinal)
			assert.Equal(t, []string{"clone response"}, result.Responses)
		})
	})

	t.Run("Given expert strategy returns error", func(t *testing.T) {
		t.Parallel()

		agentResponse := model_builders.NewAgentResponseBuilder().
			WithRole(enums.AgentRoleExpert).
			WithEnsDomain("expert-error.protocol6022.eth").
			Build()

		initSuite := func(suite *WhenComputingAgentTestingSuite, request *requests.AgentComputeRequest) {
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agentResponse, nil)

			suite.mockExpertComputeStrategy.EXPECT().
				Compute(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenComputingAgentBeforeEach(t)
			request := generateValidAgentComputeRequest()
			initSuite(suite, request)

			result, err := suite.sut.Compute(request)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "failed to compute agent")
		})
	})
}
