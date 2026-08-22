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
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_http_agent_swarm_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRegisteringAgentTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenRegisteringAgentBeforeEach(t *testing.T) *WhenRegisteringAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenRegisteringAgentTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenRegisteringAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker creates the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite) uuid.UUID {
			participantId := uuid.New()

			suite.mockClient.EXPECT().
				PostParticipantsAgentsWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				DoAndReturn(func(_ any, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostParticipantsAgentsResponse, error) {
					var request clients.RequestsRegisterAgentRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, "0xabc", *request.CollectionAddress)
					assert.Equal(t, "7", *request.CollectionTokenId)

					return &clients.PostParticipantsAgentsResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusCreated},
						JSON201:      &clients.ResponsesParticipantResponse{Id: lo.ToPtr(participantId.String())},
					}, nil
				}).
				Times(1)

			return participantId
		}

		t.Run("Should return the participant id", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			participantId := initSuite(suite)

			got, err := suite.sut.RegisterAgent(context.Background(), "0xabc", 7)

			assert.NoError(t, err)
			assert.Equal(t, participantId, got)
		})
	})

	t.Run("Given the broker already knows the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite) uuid.UUID {
			participantId := uuid.New()

			suite.mockClient.EXPECT().
				PostParticipantsAgentsWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostParticipantsAgentsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					JSON200:      &clients.ResponsesParticipantResponse{Id: lo.ToPtr(participantId.String())},
				}, nil).
				Times(1)

			return participantId
		}

		t.Run("Should return the existing participant id", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			participantId := initSuite(suite)

			got, err := suite.sut.RegisterAgent(context.Background(), "0xabc", 7)

			assert.NoError(t, err)
			assert.Equal(t, participantId, got)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite) {
			suite.mockClient.EXPECT().
				PostParticipantsAgentsWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostParticipantsAgentsResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
					Body:         []byte("boom"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.RegisterAgent(context.Background(), "0xabc", 7)

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringAgentTestingSuite) {
			suite.mockClient.EXPECT().
				PostParticipantsAgentsWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringAgentBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.RegisterAgent(context.Background(), "0xabc", 7)

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})
}
