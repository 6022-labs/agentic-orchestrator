package services_test

import (
	"errors"
	"testing"

	external_models "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_http_conversation_broker_mocks/clients_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenEnsuringBridgeRegistrationTestingSuite struct {
	sut *services.BridgeRegistrar

	mockClient *clients_mocks.MockConversationBrokerClientInterface
	bridgeId   uuid.UUID
}

func WhenEnsuringBridgeRegistrationBeforeEach(t *testing.T) *WhenEnsuringBridgeRegistrationTestingSuite {
	ctrl := gomock.NewController(t)

	mockClient := clients_mocks.NewMockConversationBrokerClientInterface(ctrl)
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
				RegisterBridge(gomock.Any()).
				DoAndReturn(func(req *external_models.ExternalRegisterBridgeRequest) (*external_models.ExternalBridgeResponse, error) {
					assert.NotNil(t, req.ID)
					assert.Equal(t, suite.bridgeId, *req.ID)
					assert.Equal(t, "slack", req.Kind)
					assert.Equal(t, "http://slack.local/cb", req.CallbackUrl)
					return &external_models.ExternalBridgeResponse{ID: *req.ID, Kind: "slack", CallbackUrl: req.CallbackUrl}, nil
				}).
				Times(1)
		}

		t.Run("Should call the broker with the configured id", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			got, err := suite.sut.Ensure()
			assert.NoError(t, err)
			assert.Equal(t, suite.bridgeId, got)
		})
	})

	t.Run("Given the broker call fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringBridgeRegistrationTestingSuite) {
			suite.mockClient.EXPECT().RegisterBridge(gomock.Any()).Return(nil, errors.New("boom")).Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Ensure()

			assert.Error(t, err)
		})
	})

	t.Run("Given the broker returns a different id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringBridgeRegistrationTestingSuite) {
			suite.mockClient.EXPECT().
				RegisterBridge(gomock.Any()).
				Return(&external_models.ExternalBridgeResponse{ID: uuid.New(), Kind: "slack", CallbackUrl: "http://slack.local/cb"}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringBridgeRegistrationBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Ensure()

			assert.Error(t, err)
		})
	})
}
