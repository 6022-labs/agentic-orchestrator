package services_test

import (
	"context"

	"errors"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentUrlTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenGettingAgentUrlBeforeEach(t *testing.T) *WhenGettingAgentUrlTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenGettingAgentUrlTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenGettingAgentUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker returns the url", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockClient.EXPECT().
				GetRemoteAgentsEnsDomainUrlWithResponse(gomock.Any(), "ada.experts.80002.6022.eth").
				Return(&clients.GetRemoteAgentsEnsDomainUrlResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200: &clients.ResponsesAgentUrlResponse{
						Url: lo.ToPtr("https://ada.example"),
					},
				}, nil).
				Times(1)
		}

		t.Run("Should return the mapped url", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentUrl(context.Background(), "ada.experts.80002.6022.eth")

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, "https://ada.example", result.Url)
			}
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockClient.EXPECT().
				GetRemoteAgentsEnsDomainUrlWithResponse(gomock.Any(), "ada.experts.80002.6022.eth").
				Return(&clients.GetRemoteAgentsEnsDomainUrlResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusBadGateway},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentUrl(context.Background(), "ada.experts.80002.6022.eth")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentUrlTestingSuite) {
			suite.mockClient.EXPECT().
				GetRemoteAgentsEnsDomainUrlWithResponse(gomock.Any(), "ada.experts.80002.6022.eth").
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentUrlBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentUrl(context.Background(), "ada.experts.80002.6022.eth")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
