package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentTestingSuite struct {
	sut *use_cases.GetAgent

	agentRegistryRequesterMock *interfaces_mocks.MockAgentRegistryRequesterInterface
}

func WhenGettingAgentBeforeEach(t *testing.T) *WhenGettingAgentTestingSuite {
	mockController := gomock.NewController(t)

	agentRegistryRequesterMock := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)

	sut := use_cases.NewGetAgent(agentRegistryRequesterMock)

	return &WhenGettingAgentTestingSuite{
		sut: sut,

		agentRegistryRequesterMock: agentRegistryRequesterMock,
	}
}

func TestWhenGettingAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agent from the registry", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.agentRegistryRequesterMock.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			req := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			resp, err := suite.sut.Get(req)

			assert.Nil(t, resp)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.agentRegistryRequesterMock.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Return(nil, nil)
		}

		t.Run("Should return nil response and nil error", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			req := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			resp, err := suite.sut.Get(req)

			assert.Nil(t, resp)
			assert.Nil(t, err)
		})
	})

	t.Run("Given agent is fetched successfully", func(t *testing.T) {
		t.Parallel()

		agent := model_builders.NewAgentBuilder().Build()

		initSuite := func(suite *WhenGettingAgentTestingSuite) {
			suite.agentRegistryRequesterMock.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Return(agent, nil)
		}

		t.Run("Should return the agent response", func(t *testing.T) {
			t.Parallel()
			suite := WhenGettingAgentBeforeEach(t)
			initSuite(suite)

			req := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			resp, err := suite.sut.Get(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, agent.ID, resp.ID)
			assert.Equal(t, agent.Name, resp.Name)
		})
	})
}
