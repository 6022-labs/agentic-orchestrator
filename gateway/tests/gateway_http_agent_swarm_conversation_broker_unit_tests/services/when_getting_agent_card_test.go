package services_test

import (
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentCardTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenGettingAgentCardBeforeEach(t *testing.T) *WhenGettingAgentCardTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenGettingAgentCardTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenGettingAgentCard(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker returns the card", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentCardTestingSuite) {
			suite.mockClient.EXPECT().GetRemoteAgentsEnsDomainCardWithResponse(gomock.Any(), "ada.experts.80002.6022.eth").Return(&clients.GetRemoteAgentsEnsDomainCardResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &clients.ResponsesRemoteAgentCardResponse{
					Protocol:        lo.ToPtr("6022"),
					ProtocolVersion: lo.ToPtr("1"),
					Name:            lo.ToPtr("Ada"),
					EnsDomain:       lo.ToPtr("ada.experts.80002.6022.eth"),
					Endpoints: &clients.ResponsesRemoteAgentCardEndpointsResponse{
						Responses: lo.ToPtr("https://ada.example/responses"),
					},
					PaymentMethods: lo.ToPtr([]string{"x402"}),
				},
			}, nil).Times(1)
		}

		t.Run("Should return the mapped card", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentCard("ada.experts.80002.6022.eth")

			assert.NoError(t, err)
			if assert.NotNil(t, result) {
				assert.Equal(t, "6022", result.Protocol)
				assert.Equal(t, "Ada", result.Name)
				assert.Equal(t, "https://ada.example/responses", result.Endpoints.Responses)
				assert.Equal(t, []string{"x402"}, result.PaymentMethods)
			}
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentCardTestingSuite) {
			suite.mockClient.EXPECT().GetRemoteAgentsEnsDomainCardWithResponse(gomock.Any(), "ada.experts.80002.6022.eth").Return(&clients.GetRemoteAgentsEnsDomainCardResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusBadGateway},
				Body:         []byte("boom"),
			}, nil).Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentCardBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.GetAgentCard("ada.experts.80002.6022.eth")

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
