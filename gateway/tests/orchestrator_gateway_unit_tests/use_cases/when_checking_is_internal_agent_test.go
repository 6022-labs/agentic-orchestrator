package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenCheckingIsInternalAgentTestSuite struct {
	agentRepositoryMock *repositories_mocks.MockAgentRepositoryInterface
	sut                 *use_cases.IsInternalAgent
}

func WhenCheckingIsInternalAgentBeforeEach(t *testing.T) *WhenCheckingIsInternalAgentTestSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	sut := use_cases.NewIsInternalAgent(agentRepositoryMock)

	return &WhenCheckingIsInternalAgentTestSuite{
		agentRepositoryMock: agentRepositoryMock,
		sut:                 sut,
	}
}

func TestWhenCheckingIsInternalAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIsInternalAgentTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, nil)
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

	t.Run("Given error occurs while fetching agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIsInternalAgentTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, gorm.ErrRecordNotFound)
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
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(true).
				Build()
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
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
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				Build()
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
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
