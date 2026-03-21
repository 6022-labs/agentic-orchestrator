package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCheckingIsInternalAgentTestSuite struct {
	agentRegistryRequesterMock *interfaces_mocks.MockAgentRegistryRequesterInterface
	sut                        *use_cases.IsInternalAgent
}

func WhenCheckingIsInternalAgentBeforeEach(t *testing.T) *WhenCheckingIsInternalAgentTestSuite {
	mockController := gomock.NewController(t)

	agentRegistryRequesterMock := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	sut := use_cases.NewIsInternalAgent(agentRegistryRequesterMock)

	return &WhenCheckingIsInternalAgentTestSuite{
		agentRegistryRequesterMock: agentRegistryRequesterMock,
		sut:                        sut,
	}
}

func TestWhenCheckingIsInternalAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIsInternalAgentTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRegistryRequesterMock.EXPECT().IsInternalAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(&responses.IsInternalAgentResponse{IsInternalAgent: false}, nil)
		}

		t.Run("Should return false with no error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			suite := WhenCheckingIsInternalAgentBeforeEach(t)
			initSuite(suite, request)

			response, err := suite.sut.IsInternal(request)

			assert.NoError(t, err)
			assert.False(t, response.IsInternalAgent)
		})
	})

	t.Run("Given error occurs while fetching internal status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIsInternalAgentTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRegistryRequesterMock.EXPECT().IsInternalAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			suite := WhenCheckingIsInternalAgentBeforeEach(t)
			initSuite(suite, request)

			response, err := suite.sut.IsInternal(request)

			assert.Error(t, err)
			assert.Nil(t, response)
		})
	})

	t.Run("Given agent is found and is internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIsInternalAgentTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			_ = model_builders.NewAgentBuilder().WithIsInternal(true).Build()
			suite.agentRegistryRequesterMock.EXPECT().IsInternalAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(&responses.IsInternalAgentResponse{IsInternalAgent: true}, nil)
		}

		t.Run("Should return true with no error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			suite := WhenCheckingIsInternalAgentBeforeEach(t)
			initSuite(suite, request)

			response, err := suite.sut.IsInternal(request)

			assert.NoError(t, err)
			assert.True(t, response.IsInternalAgent)
		})
	})

	t.Run("Given agent is found and is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIsInternalAgentTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			_ = model_builders.NewAgentBuilder().WithIsInternal(false).Build()
			suite.agentRegistryRequesterMock.EXPECT().IsInternalAgent(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(&responses.IsInternalAgentResponse{IsInternalAgent: false}, nil)
		}

		t.Run("Should return false with no error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			suite := WhenCheckingIsInternalAgentBeforeEach(t)
			initSuite(suite, request)

			response, err := suite.sut.IsInternal(request)

			assert.NoError(t, err)
			assert.False(t, response.IsInternalAgent)
		})
	})
}
