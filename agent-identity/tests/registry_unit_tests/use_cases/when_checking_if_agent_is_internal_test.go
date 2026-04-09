package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCheckingIfAgentIsInternalTestingSuite struct {
	sut *use_cases.IsInternalAgent

	mockAgentRepository *repositories_mocks.MockAgentRepositoryInterface
}

func WhenCheckingIfAgentIsInternalBeforeEach(t *testing.T) *WhenCheckingIfAgentIsInternalTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := use_cases.NewIsInternalAgent(mockAgentRepository)

	return &WhenCheckingIfAgentIsInternalTestingSuite{
		sut:                 sut,
		mockAgentRepository: mockAgentRepository,
	}
}

func TestWhenCheckingIfAgentIsInternal(t *testing.T) {
	t.Parallel()

	t.Run("Given error happens while fetching agent from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIfAgentIsInternalTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenCheckingIfAgentIsInternalBeforeEach(t)
			initSuite(basics)

			req := &requests.AgentByCollectionTokenIdRequest{
				AgentCollectionAddress: "0x1234567890abcdef1234567890abcdef12345678",
				AgentCollectionTokenId: 1,
			}
			resp, err := basics.sut.IsInternal(req)

			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingIfAgentIsInternalTestingSuite) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return false for IsInternalAgent", func(t *testing.T) {
			t.Parallel()
			basics := WhenCheckingIfAgentIsInternalBeforeEach(t)
			initSuite(basics)

			req := &requests.AgentByCollectionTokenIdRequest{
				AgentCollectionAddress: "0x1234567890abcdef1234567890abcdef12345678",
				AgentCollectionTokenId: 1,
			}
			resp, err := basics.sut.IsInternal(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.False(t, resp.IsInternalAgent)
		})
	})

	t.Run("Given agent is internal", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return true for IsInternalAgent", func(t *testing.T) {
			t.Parallel()
			basics := WhenCheckingIfAgentIsInternalBeforeEach(t)

			agent := model_builders.NewAgentBuilder().
				WithIsInternal(true).
				Build()

			basics.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(agent.CollectionAddress, *agent.CollectionTokenId).
				Return(agent, nil).
				Times(1)

			req := &requests.AgentByCollectionTokenIdRequest{
				AgentCollectionAddress: agent.CollectionAddress,
				AgentCollectionTokenId: *agent.CollectionTokenId,
			}
			resp, err := basics.sut.IsInternal(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.True(t, resp.IsInternalAgent)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false for IsInternalAgent", func(t *testing.T) {
			t.Parallel()
			basics := WhenCheckingIfAgentIsInternalBeforeEach(t)

			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				Build()

			basics.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(agent.CollectionAddress, *agent.CollectionTokenId).
				Return(agent, nil).
				Times(1)

			req := &requests.AgentByCollectionTokenIdRequest{
				AgentCollectionAddress: agent.CollectionAddress,
				AgentCollectionTokenId: *agent.CollectionTokenId,
			}
			resp, err := basics.sut.IsInternal(req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.False(t, resp.IsInternalAgent)
		})
	})
}
