package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenConsumingWalletOwnershipChallengeTestingSuite struct {
	sut *services.WalletOwnershipChallenge

	mockCacheStorage interfaces_mocks.MockCacheStorageInterface
}

func WhenConsumingWalletOwnershipChallengeBeforeEach(t *testing.T) *WhenConsumingWalletOwnershipChallengeTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)

	sut := services.NewWalletOwnershipChallenge(mockCacheStorage)

	return &WhenConsumingWalletOwnershipChallengeTestingSuite{
		sut:              sut,
		mockCacheStorage: *mockCacheStorage,
	}
}

func TestWhenConsumingWalletOwnershipChallenge(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while storing challenge", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConsumingWalletOwnershipChallengeTestingSuite, nonce uuid.UUID) {
			suite.mockCacheStorage.EXPECT().
				ReadString(fmt.Sprintf(services.WALLET_OWNERSHIP_CHALLENGE_PATH_TEMPLATE, nonce.String())).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenConsumingWalletOwnershipChallengeBeforeEach(t)

			uuid := uuid.New()

			initSuite(suite, uuid)

			_, err := suite.sut.ConsumeChallenge(&uuid)

			assert.Error(t, err)
		})
	})

	t.Run("Given challenge does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConsumingWalletOwnershipChallengeTestingSuite, nonce uuid.UUID) {
			suite.mockCacheStorage.EXPECT().
				ReadString(fmt.Sprintf(services.WALLET_OWNERSHIP_CHALLENGE_PATH_TEMPLATE, nonce.String())).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenConsumingWalletOwnershipChallengeBeforeEach(t)

			uuid := uuid.New()

			initSuite(suite, uuid)

			_, err := suite.sut.ConsumeChallenge(&uuid)

			assert.Error(t, err)
		})
	})

	t.Run("Given challenge exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConsumingWalletOwnershipChallengeTestingSuite, nonce uuid.UUID, userId string) {
			suite.mockCacheStorage.EXPECT().
				ReadString(fmt.Sprintf(services.WALLET_OWNERSHIP_CHALLENGE_PATH_TEMPLATE, nonce.String())).
				Return(&userId, nil)
		}

		t.Run("Should return userId", func(t *testing.T) {
			t.Parallel()

			suite := WhenConsumingWalletOwnershipChallengeBeforeEach(t)

			uuid := uuid.New()
			userId := "userId"

			initSuite(suite, uuid, userId)

			result, err := suite.sut.ConsumeChallenge(&uuid)

			assert.NoError(t, err)
			assert.Equal(t, &userId, result)
		})
	})
}
