package services_test

import (
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

type WhenRegisteringUserTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenRegisteringUserBeforeEach(t *testing.T) *WhenRegisteringUserTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	return &WhenRegisteringUserTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenRegisteringUser(t *testing.T) {
	t.Parallel()

	t.Run("Given the broker creates the user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringUserTestingSuite) uuid.UUID {
			participantId := uuid.New()

			suite.mockClient.EXPECT().
				PostParticipantsUsersWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				DoAndReturn(func(_ any, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostParticipantsUsersResponse, error) {
					var request clients.RequestsRegisterUserRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, "Pierre", *request.DisplayName)

					return &clients.PostParticipantsUsersResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusCreated},
						JSON201:      &clients.ResponsesParticipantResponse{Id: lo.ToPtr(participantId.String())},
					}, nil
				}).
				Times(1)

			return participantId
		}

		t.Run("Should return the participant id", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringUserBeforeEach(t)
			participantId := initSuite(suite)

			got, err := suite.sut.RegisterUser("Pierre")

			assert.NoError(t, err)
			assert.Equal(t, participantId, got)
		})
	})

	t.Run("Given the broker returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringUserTestingSuite) {
			suite.mockClient.EXPECT().
				PostParticipantsUsersWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostParticipantsUsersResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusBadRequest},
					Body:         []byte("invalid"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringUserBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.RegisterUser("Pierre")

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringUserTestingSuite) {
			suite.mockClient.EXPECT().
				PostParticipantsUsersWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("broker down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringUserBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.RegisterUser("Pierre")

			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, got)
		})
	})
}
