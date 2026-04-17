package services_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_http_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRequestingConversationBrokerTestingSuite struct {
	sut *services.ConversationBrokerRequester

	mockClient *clients_mocks.MockConversationBrokerClientInterface
}

func WhenRequestingConversationBrokerBeforeEach(t *testing.T) *WhenRequestingConversationBrokerTestingSuite {
	ctrl := gomock.NewController(t)
	mockClient := clients_mocks.NewMockConversationBrokerClientInterface(ctrl)

	return &WhenRequestingConversationBrokerTestingSuite{
		sut:        services.NewConversationBrokerRequester(mockClient),
		mockClient: mockClient,
	}
}

func TestWhenRequestingConversationBroker(t *testing.T) {
	t.Parallel()

	t.Run("RegisterAgent", func(t *testing.T) {
		t.Parallel()

		t.Run("Given the client succeeds", func(t *testing.T) {
			t.Parallel()

			expectedId := uuid.New()

			initSuite := func(suite *WhenRequestingConversationBrokerTestingSuite) {
				suite.mockClient.EXPECT().
					RegisterAgent(&models.ExternalRegisterAgentRequest{CollectionAddress: "0xabc", CollectionTokenId: "42"}).
					Return(&models.ExternalParticipantResponse{ID: expectedId}, nil)
			}

			t.Run("Should return the broker participant id", func(t *testing.T) {
				t.Parallel()

				suite := WhenRequestingConversationBrokerBeforeEach(t)
				initSuite(suite)

				id, err := suite.sut.RegisterAgent("0xabc", 42)

				assert.NoError(t, err)
				assert.Equal(t, expectedId, id)
			})
		})

		t.Run("Given the client fails", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenRequestingConversationBrokerTestingSuite) {
				suite.mockClient.EXPECT().RegisterAgent(gomock.Any()).Return(nil, errors.New("boom"))
			}

			t.Run("Should propagate the error", func(t *testing.T) {
				t.Parallel()

				suite := WhenRequestingConversationBrokerBeforeEach(t)
				initSuite(suite)

				_, err := suite.sut.RegisterAgent("0xabc", 1)
				assert.Error(t, err)
			})
		})
	})

	t.Run("RegisterUser", func(t *testing.T) {
		t.Parallel()

		t.Run("Given the client succeeds", func(t *testing.T) {
			t.Parallel()

			expectedId := uuid.New()

			initSuite := func(suite *WhenRequestingConversationBrokerTestingSuite) {
				suite.mockClient.EXPECT().
					RegisterUser(&models.ExternalRegisterUserRequest{DisplayName: "bob"}).
					Return(&models.ExternalParticipantResponse{ID: expectedId}, nil)
			}

			t.Run("Should return the broker participant id", func(t *testing.T) {
				t.Parallel()

				suite := WhenRequestingConversationBrokerBeforeEach(t)
				initSuite(suite)

				id, err := suite.sut.RegisterUser("bob")

				assert.NoError(t, err)
				assert.Equal(t, expectedId, id)
			})
		})
	})

	t.Run("IngestMessage", func(t *testing.T) {
		t.Parallel()

		t.Run("Given the client succeeds", func(t *testing.T) {
			t.Parallel()

			convId := uuid.New()
			participantId := uuid.New()
			expectedMsgId := uuid.New()

			initSuite := func(suite *WhenRequestingConversationBrokerTestingSuite) {
				suite.mockClient.EXPECT().
					IngestMessage(convId, &models.ExternalIngestMessageRequest{ParticipantID: participantId, Content: "hello"}).
					Return(&models.ExternalConversationMessageResponse{ID: expectedMsgId}, nil)
			}

			t.Run("Should return the broker message id", func(t *testing.T) {
				t.Parallel()

				suite := WhenRequestingConversationBrokerBeforeEach(t)
				initSuite(suite)

				id, err := suite.sut.IngestMessage(convId, participantId, "hello")

				assert.NoError(t, err)
				assert.Equal(t, expectedMsgId, id)
			})
		})
	})

	t.Run("ForceConcludeConversation", func(t *testing.T) {
		t.Parallel()

		t.Run("Given the client succeeds", func(t *testing.T) {
			t.Parallel()

			convId := uuid.New()

			initSuite := func(suite *WhenRequestingConversationBrokerTestingSuite) {
				suite.mockClient.EXPECT().ForceConcludeConversation(convId).Return(nil)
			}

			t.Run("Should not return an error", func(t *testing.T) {
				t.Parallel()

				suite := WhenRequestingConversationBrokerBeforeEach(t)
				initSuite(suite)

				err := suite.sut.ForceConcludeConversation(convId)
				assert.NoError(t, err)
			})
		})
	})

	t.Run("AttachParticipant", func(t *testing.T) {
		t.Parallel()

		t.Run("Given the client succeeds", func(t *testing.T) {
			t.Parallel()

			convId := uuid.New()
			participantId := uuid.New()

			initSuite := func(suite *WhenRequestingConversationBrokerTestingSuite) {
				suite.mockClient.EXPECT().
					AttachParticipant(convId, &models.ExternalConversationParticipantRequest{ParticipantID: participantId}).
					Return(nil)
			}

			t.Run("Should not return an error", func(t *testing.T) {
				t.Parallel()

				suite := WhenRequestingConversationBrokerBeforeEach(t)
				initSuite(suite)

				err := suite.sut.AttachParticipant(convId, participantId)
				assert.NoError(t, err)
			})
		})
	})

	t.Run("DetachParticipant", func(t *testing.T) {
		t.Parallel()

		t.Run("Given the client succeeds", func(t *testing.T) {
			t.Parallel()

			convId := uuid.New()
			participantId := uuid.New()

			initSuite := func(suite *WhenRequestingConversationBrokerTestingSuite) {
				suite.mockClient.EXPECT().DetachParticipant(convId, participantId).Return(nil)
			}

			t.Run("Should not return an error", func(t *testing.T) {
				t.Parallel()

				suite := WhenRequestingConversationBrokerBeforeEach(t)
				initSuite(suite)

				err := suite.sut.DetachParticipant(convId, participantId)
				assert.NoError(t, err)
			})
		})
	})
}
