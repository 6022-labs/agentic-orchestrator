package services_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenIssuingWalletOwnershipChallengeTestingSuite struct {
	sut *services.WalletOwnershipChallenge

	mockCacheStorage interfaces_mocks.MockCacheStorageInterface
}

func WhenIssuingWalletOwnershipChallengeBeforeEach(t *testing.T) *WhenIssuingWalletOwnershipChallengeTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)

	sut := services.NewWalletOwnershipChallenge(mockCacheStorage)

	return &WhenIssuingWalletOwnershipChallengeTestingSuite{
		sut:              sut,
		mockCacheStorage: *mockCacheStorage,
	}
}

func TestWhenIssuingWalletOwnershipChallenge(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while storing challenge", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenIssuingWalletOwnershipChallengeTestingSuite, userId string) {
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), userId, gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenIssuingWalletOwnershipChallengeBeforeEach(t)

			userId := "userId"

			initSuite(suite, userId)

			_, err := suite.sut.IssueChallenge(context.Background(), userId)

			assert.Error(t, err)
		})
	})

	t.Run("Given challenge is stored successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenIssuingWalletOwnershipChallengeTestingSuite, userId string) {
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), userId, gomock.Any()).
				Return(nil)
		}

		t.Run("Should return nonce", func(t *testing.T) {
			t.Parallel()

			suite := WhenIssuingWalletOwnershipChallengeBeforeEach(t)

			userId := "userId"

			initSuite(suite, userId)

			nonce, err := suite.sut.IssueChallenge(context.Background(), userId)

			assert.NoError(t, err)
			assert.NotNil(t, nonce)
		})
	})
}
