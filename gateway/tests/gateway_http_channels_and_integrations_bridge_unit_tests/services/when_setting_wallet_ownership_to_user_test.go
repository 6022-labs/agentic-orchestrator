package services_test

import (
	"context"

	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/services"
	clients_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_http_channels_and_integrations_bridge_mocks/clients_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSettingWalletOwnershipToUserTestingSuite struct {
	sut *services.PlatformBridgeRequester

	mockClient *clients_mocks.MockClientWithResponsesInterface
}

func WhenSettingWalletOwnershipToUserBeforeEach(t *testing.T) *WhenSettingWalletOwnershipToUserTestingSuite {
	mockController := gomock.NewController(t)

	mockClient := clients_mocks.NewMockClientWithResponsesInterface(mockController)

	sut := services.NewPlatformBridgeRequester(map[enums.BridgePlatform]clients.ClientWithResponsesInterface{
		enums.BridgePlatformSlack: mockClient,
	})

	return &WhenSettingWalletOwnershipToUserTestingSuite{
		sut:        sut,
		mockClient: mockClient,
	}
}

func TestWhenSettingWalletOwnershipToUser(t *testing.T) {
	t.Parallel()

	t.Run("Given the slack bridge accepts the request", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingWalletOwnershipToUserTestingSuite) {
			suite.mockClient.EXPECT().
				PostWalletSetOwnershipToSlackUserWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				DoAndReturn(func(_ any, _ string, body io.Reader, _ ...clients.RequestEditorFn) (*clients.PostWalletSetOwnershipToSlackUserResponse, error) {
					var request clients.RequestsSetWalletOwnershipToSlackUserRequest
					assert.NoError(t, json.NewDecoder(body).Decode(&request))
					assert.Equal(t, "challenge-1", *request.Challenge)
					assert.Equal(t, "0xwallet", *request.WalletAddress)

					return &clients.PostWalletSetOwnershipToSlackUserResponse{
						HTTPResponse: &http.Response{StatusCode: http.StatusOK},
					}, nil
				}).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToUserBeforeEach(t)
			initSuite(suite)

			assert.NoError(
				t,
				suite.sut.SetWalletOwnershipToUser(
					context.Background(),
					enums.BridgePlatformSlack,
					"challenge-1",
					"0xwallet",
				),
			)
		})
	})

	t.Run("Given the slack bridge returns an unexpected status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingWalletOwnershipToUserTestingSuite) {
			suite.mockClient.EXPECT().
				PostWalletSetOwnershipToSlackUserWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(&clients.PostWalletSetOwnershipToSlackUserResponse{
					HTTPResponse: &http.Response{StatusCode: http.StatusBadRequest},
					Body:         []byte("invalid"),
				}, nil).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToUserBeforeEach(t)
			initSuite(suite)

			assert.Error(
				t,
				suite.sut.SetWalletOwnershipToUser(
					context.Background(),
					enums.BridgePlatformSlack,
					"challenge-1",
					"0xwallet",
				),
			)
		})
	})

	t.Run("Given the transport fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingWalletOwnershipToUserTestingSuite) {
			suite.mockClient.EXPECT().
				PostWalletSetOwnershipToSlackUserWithBodyWithResponse(gomock.Any(), "application/json", gomock.Any()).
				Return(nil, errors.New("bridge down")).
				Times(1)
		}

		t.Run("Should return the transport error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToUserBeforeEach(t)
			initSuite(suite)

			assert.Error(
				t,
				suite.sut.SetWalletOwnershipToUser(
					context.Background(),
					enums.BridgePlatformSlack,
					"challenge-1",
					"0xwallet",
				),
			)
		})
	})

	t.Run("Given an unsupported platform", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error without calling any client", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingWalletOwnershipToUserBeforeEach(t)

			err := suite.sut.SetWalletOwnershipToUser(
				context.Background(),
				enums.BridgePlatform("teams"),
				"challenge-1",
				"0xwallet",
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "unsupported platform")
		})
	})
}
