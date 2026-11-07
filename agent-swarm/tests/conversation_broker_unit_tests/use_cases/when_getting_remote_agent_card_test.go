package use_cases_test

import (
	"context"
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type WhenGettingRemoteAgentCardTestingSuite struct {
	sut *use_cases.GetRemoteAgentCard

	mockRemoteAgentCardResolver *interfaces_mocks.MockRemoteAgentCardResolverInterface
}

func WhenGettingRemoteAgentCardBeforeEach(t *testing.T) *WhenGettingRemoteAgentCardTestingSuite {
	mockController := gomock.NewController(t)

	mockRemoteAgentCardResolver := interfaces_mocks.NewMockRemoteAgentCardResolverInterface(mockController)

	sut := use_cases.NewGetRemoteAgentCard(mockRemoteAgentCardResolver)

	return &WhenGettingRemoteAgentCardTestingSuite{
		sut: sut,

		mockRemoteAgentCardResolver: mockRemoteAgentCardResolver,
	}
}

func TestWhenGettingRemoteAgentCard(t *testing.T) {
	t.Parallel()

	t.Run("Given ensDomain is empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingRemoteAgentCardBeforeEach(t)

			_, err := suite.sut.Get(context.Background(), "   ")

			var validationErr *cb_errors.ValidationError
			require.ErrorAs(t, err, &validationErr)
			assert.Equal(t, "ensDomain", validationErr.Field)
		})
	})

	t.Run("Given resolver fails", func(t *testing.T) {
		t.Parallel()

		cause := errors.New("dial tcp 10.0.0.1:443: connect: connection refused")

		initSuite := func(suite *WhenGettingRemoteAgentCardTestingSuite) {
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), "agent.6022.eth").
				Return(nil, cause).
				Times(1)
		}

		t.Run("Should return an unavailable error carrying the cause", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(context.Background(), "agent.6022.eth")

			var unavailableErr *cb_errors.UnavailableError
			require.ErrorAs(t, err, &unavailableErr)
			assert.Equal(t, "agent_node_unreachable", unavailableErr.Code)
			assert.Equal(t, "failed to reach agent node", unavailableErr.Message)
			assert.ErrorIs(t, err, cause)
		})
	})

	t.Run("Given resolver succeeds", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingRemoteAgentCardTestingSuite) {
			suite.mockRemoteAgentCardResolver.EXPECT().
				ResolveCard(gomock.Any(), "agent.6022.eth").
				Return(&models.RemoteAgentCard{Name: "agent", EnsDomain: "agent.6022.eth"}, nil).
				Times(1)
		}

		t.Run("Should return the card", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.Get(context.Background(), " agent.6022.eth ")

			require.NoError(t, err)
			assert.Equal(t, "agent", got.Name)
			assert.Equal(t, "agent.6022.eth", got.EnsDomain)
		})
	})
}
