package use_cases_test

import (
	"context"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type WhenResolvingAgentUrlTestingSuite struct {
	sut *use_cases.ResolveAgentUrl

	mockAgentUrlResolver *interfaces_mocks.MockAgentUrlResolverInterface
}

func WhenResolvingAgentUrlBeforeEach(t *testing.T) *WhenResolvingAgentUrlTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentUrlResolver := interfaces_mocks.NewMockAgentUrlResolverInterface(mockController)

	sut := use_cases.NewResolveAgentUrl(mockAgentUrlResolver)

	return &WhenResolvingAgentUrlTestingSuite{
		sut: sut,

		mockAgentUrlResolver: mockAgentUrlResolver,
	}
}

func TestWhenResolvingAgentUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given ensDomain is empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)

			_, err := suite.sut.Resolve(context.Background(), &requests.ResolveAgentUrlRequest{EnsDomain: "   "})

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given resolver fails", func(t *testing.T) {
		t.Parallel()

		cause := errors.New("boom")

		initSuite := func(suite *WhenResolvingAgentUrlTestingSuite) {
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.6022.eth").
				Return("", cause).
				Times(1)
		}

		t.Run("Should return an unavailable error carrying the cause", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Resolve(
				context.Background(),
				&requests.ResolveAgentUrlRequest{EnsDomain: "agent.6022.eth"},
			)

			var unavailableErr *cb_errors.UnavailableError
			require.ErrorAs(t, err, &unavailableErr)
			assert.Equal(t, "agent_url_unresolved", unavailableErr.Code)
			assert.Equal(t, "failed to resolve agent url", unavailableErr.Message)
			assert.ErrorIs(t, err, cause)
		})
	})

	t.Run("Given resolver succeeds", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingAgentUrlTestingSuite) {
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.6022.eth").
				Return("https://agent.example.com", nil).
				Times(1)
		}

		t.Run("Should return the resolved url", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.Resolve(
				context.Background(),
				&requests.ResolveAgentUrlRequest{EnsDomain: " agent.6022.eth "},
			)

			assert.NoError(t, err)
			assert.Equal(t, "https://agent.example.com", got.Url)
		})
	})
}
