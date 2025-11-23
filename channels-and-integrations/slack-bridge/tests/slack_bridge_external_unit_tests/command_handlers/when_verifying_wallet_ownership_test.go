package command_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	external_interfaces_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenVerifyingWalletOwnershipTestingSuite struct {
	sut *command_handlers.VerifyWalletOwnership

	verifyWalletOwnershipSettings   *settings.VerifyWalletOwnershipSettings
	mockWalletOwnershipChallenge    *interfaces_mocks.MockWalletOwnershipChallengeInterface
	mockSlackUserSynchronizerCached *external_interfaces_mocks.MockSlackUserSynchronizerCachedInterface
}

func WhenVerifyingWalletOwnershipBeforeEach(t *testing.T) *WhenVerifyingWalletOwnershipTestingSuite {
	mockController := gomock.NewController(t)

	verifyWalletOwnershipSettings := &settings.VerifyWalletOwnershipSettings{
		OrchestratorFrontBaseUrl: "https://example.com/verify-wallet",
	}

	mockWalletOwnershipChallenge := interfaces_mocks.NewMockWalletOwnershipChallengeInterface(mockController)
	mockSlackUserSynchronizerCached := external_interfaces_mocks.NewMockSlackUserSynchronizerCachedInterface(mockController)

	sut := command_handlers.NewVerifyWalletOwnership(
		zap.NewNop(),
		verifyWalletOwnershipSettings,
		mockWalletOwnershipChallenge,
		mockSlackUserSynchronizerCached,
	)

	return &WhenVerifyingWalletOwnershipTestingSuite{
		sut: sut,

		mockWalletOwnershipChallenge:    mockWalletOwnershipChallenge,
		verifyWalletOwnershipSettings:   verifyWalletOwnershipSettings,
		mockSlackUserSynchronizerCached: mockSlackUserSynchronizerCached,
	}
}

func generateValidVerifyWalletOwnershipRequest() *requests.VerifyWalletOwnershipRequest {
	return &requests.VerifyWalletOwnershipRequest{
		UserId: "test-user-id",
	}
}

func TestWhenVerifyingWalletOwnership(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Due to empty user id, should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletOwnershipBeforeEach(t)
			request := generateValidVerifyWalletOwnershipRequest()
			request.UserId = ""

			_, err := suite.sut.Verify(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while syncing slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletOwnershipTestingSuite, request *requests.VerifyWalletOwnershipRequest) {
			suite.mockSlackUserSynchronizerCached.
				EXPECT().
				Sync(request.UserId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletOwnershipBeforeEach(t)
			request := generateValidVerifyWalletOwnershipRequest()

			initSuite(suite, request)

			_, err := suite.sut.Verify(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while issuing challenge", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletOwnershipTestingSuite, request *requests.VerifyWalletOwnershipRequest) {
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.UserId).
				Build()
			suite.mockSlackUserSynchronizerCached.
				EXPECT().
				Sync(request.UserId).
				Return(slackUser, nil)

			suite.mockWalletOwnershipChallenge.
				EXPECT().
				IssueChallenge(slackUser.Id).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletOwnershipBeforeEach(t)
			request := generateValidVerifyWalletOwnershipRequest()

			initSuite(suite, request)

			_, err := suite.sut.Verify(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given challenge is successfully issued", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenVerifyingWalletOwnershipTestingSuite, request *requests.VerifyWalletOwnershipRequest, challenge uuid.UUID) {
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.UserId).
				Build()
			suite.mockSlackUserSynchronizerCached.
				EXPECT().
				Sync(request.UserId).
				Return(slackUser, nil)

			suite.mockWalletOwnershipChallenge.
				EXPECT().
				IssueChallenge(request.UserId).
				Return(&challenge, nil)
		}

		t.Run("Should return a valid verification URL", func(t *testing.T) {
			t.Parallel()

			suite := WhenVerifyingWalletOwnershipBeforeEach(t)
			request := generateValidVerifyWalletOwnershipRequest()

			challenge := uuid.New()

			initSuite(suite, request, challenge)

			response, err := suite.sut.Verify(request)

			assert.NoError(t, err)
			expectedUrl := suite.verifyWalletOwnershipSettings.OrchestratorFrontBaseUrl +
				command_handlers.VERIFY_WALLET_FRONT_PAGE_PATH +
				"?challenge=" + challenge.String() +
				"&platform=slack"

			assert.Contains(t, *response, expectedUrl)
		})
	})
}
