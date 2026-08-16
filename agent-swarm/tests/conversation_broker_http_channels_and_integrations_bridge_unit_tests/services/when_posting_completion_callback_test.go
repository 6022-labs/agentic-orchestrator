package services_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_http_channels_and_integrations_bridge_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenPostingCompletionCallbackTestingSuite struct {
	sut *services.BridgeCallbackRequester

	mockClient      *clients_mocks.MockClientWithResponsesInterface
	capturedBaseUrl *string
	factoryErr      error

	bridge       *models.Bridge
	conversation *models.Conversation
	agent        *models.Participant
}

func WhenPostingCompletionCallbackBeforeEach(t *testing.T) *WhenPostingCompletionCallbackTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	suite := &WhenPostingCompletionCallbackTestingSuite{
		mockClient:      mockClient,
		capturedBaseUrl: new(string),
		bridge:          &models.Bridge{BaseUrl: "http://slack-bridge.local"},
		conversation:    &models.Conversation{ID: uuid.New()},
		agent: &models.Participant{
			CollectionAddress: lo.ToPtr("0xabc"),
			CollectionTokenId: lo.ToPtr(uint64(7)),
		},
	}

	suite.sut = services.NewBridgeCallbackRequester(zap.NewNop(), func(baseUrl string) (clients.ClientWithResponsesInterface, error) {
		*suite.capturedBaseUrl = baseUrl
		if suite.factoryErr != nil {
			return nil, suite.factoryErr
		}
		return mockClient, nil
	})

	return suite
}

func TestWhenPostingCompletionCallback(t *testing.T) {
	t.Parallel()

	t.Run("Given the bridge accepts the callback", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingCompletionCallbackTestingSuite) {
			suite.mockClient.EXPECT().
				PostBrokerCallbackWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				DoAndReturn(func(_ any, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostBrokerCallbackResponse, error) {
					var request clients.RequestsPostBrokerCompletionToSlackRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, suite.conversation.ID.String(), *request.ConversationId)
					assert.Equal(t, "0xabc", *request.CollectionAddress)
					assert.Equal(t, 7, *request.CollectionTokenId)
					assert.Equal(t, "hello", *request.Content)

					return &clients.PostBrokerCallbackResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusNoContent},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should build the client from the bridge base url and succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingCompletionCallbackBeforeEach(t)
			initSuite(suite)

			err := suite.sut.PostCompletion(suite.bridge, suite.conversation, suite.agent, "hello")

			assert.NoError(t, err)
			assert.Equal(t, "http://slack-bridge.local", *suite.capturedBaseUrl)
		})
	})

	t.Run("Given the bridge returns a non-2xx status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingCompletionCallbackTestingSuite) {
			suite.mockClient.EXPECT().
				PostBrokerCallbackWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostBrokerCallbackResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
					Body:         []byte("not found"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingCompletionCallbackBeforeEach(t)
			initSuite(suite)

			err := suite.sut.PostCompletion(suite.bridge, suite.conversation, suite.agent, "hello")

			assert.Error(t, err)
		})
	})

	t.Run("Given the client factory fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingCompletionCallbackTestingSuite) {
			suite.factoryErr = errors.New("bad base url")
		}

		t.Run("Should return the factory error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingCompletionCallbackBeforeEach(t)
			initSuite(suite)

			err := suite.sut.PostCompletion(suite.bridge, suite.conversation, suite.agent, "hello")

			assert.Error(t, err)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingCompletionCallbackTestingSuite) {
			suite.mockClient.EXPECT().
				PostBrokerCallbackWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("bridge down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingCompletionCallbackBeforeEach(t)
			initSuite(suite)

			err := suite.sut.PostCompletion(suite.bridge, suite.conversation, suite.agent, "hello")

			assert.Error(t, err)
		})
	})

	t.Run("Given the participant has no on-chain identity", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingCompletionCallbackTestingSuite) {
			suite.agent.CollectionAddress = nil
		}

		t.Run("Should return an error without building a client", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingCompletionCallbackBeforeEach(t)
			initSuite(suite)

			err := suite.sut.PostCompletion(suite.bridge, suite.conversation, suite.agent, "hello")

			assert.Error(t, err)
			assert.Empty(t, *suite.capturedBaseUrl)
		})
	})
}
