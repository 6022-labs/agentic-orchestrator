package services_test

import (
	"context"

	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenEnsuringBridgeRegistrationTestingSuite struct {
	sut *services.BridgeRegistrar

	mockClient *clients_mocks.MockClientWithResponsesInterface
	bridgeId   uuid.UUID
}

func WhenEnsuringBridgeRegistrationBeforeEach(t *testing.T) *WhenEnsuringBridgeRegistrationTestingSuite {
	ctrl := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(ctrl)
	bridgeId := uuid.New()

	sut := services.NewBridgeRegistrar(
		zap.NewNop(),
		mockClient,
		&settings.ConversationBrokerSettings{
			BaseUrl:     "http://broker.local",
			CallbackUrl: "http://slack.local/cb",
			BridgeID:    bridgeId,
		},
	)

	return &WhenEnsuringBridgeRegistrationTestingSuite{
		sut:        sut,
		mockClient: mockClient,
		bridgeId:   bridgeId,
	}
}

func TestWhenEnsuringBridgeRegistration(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker confirms the bridge id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringBridgeRegistrationTestingSuite) {
			suite.mockClient.EXPECT().
				PostBridgesWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				DoAndReturn(func(_ any, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostBridgesResponse, error) {
					var request clients.RequestsRegisterBridgeRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, suite.bridgeId.String(), *request.Id)
					assert.Equal(t, clients.BridgeKindSlack, *request.Kind)
					assert.Equal(t, "http://slack.local/cb", *request.BaseUrl)

					return &clients.PostBridgesResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusOK},
						JSON200: &clients.ResponsesBridgeResponse{
							Id:      request.Id,
							Kind:    request.Kind,
							BaseUrl: request.BaseUrl,
						},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should call the broker with the configured id", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.Ensure(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, suite.bridgeId, got)
		})
	})

	t.Run("Given the broker newly creates the bridge", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringBridgeRegistrationTestingSuite) {
			suite.mockClient.EXPECT().
				PostBridgesWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostBridgesResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusCreated},
					JSON201: &clients.ResponsesBridgeResponse{
						Id: lo.ToPtr(suite.bridgeId.String()),
					},
				}, nil).
				Times(1)
		}

		t.Run("Should accept the created bridge", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.Ensure(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, suite.bridgeId, got)
		})
	})

	t.Run("Given the broker call fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringBridgeRegistrationTestingSuite) {
			suite.mockClient.EXPECT().
				PostBridgesWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Ensure(context.Background())

			assert.Error(t, err)
		})
	})

	t.Run("Given the broker returns a different id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringBridgeRegistrationTestingSuite) {
			suite.mockClient.EXPECT().
				PostBridgesWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostBridgesResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200: &clients.ResponsesBridgeResponse{
						Id: lo.ToPtr(uuid.New().String()),
					},
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Ensure(context.Background())

			assert.Error(t, err)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringBridgeRegistrationTestingSuite) {
			suite.mockClient.EXPECT().
				PostBridgesWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostBridgesResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Ensure(context.Background())

			assert.Error(t, err)
		})
	})
}
