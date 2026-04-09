package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenFetchingOrchestratorPrivateKeyTestingSuite struct {
	sut *services.OrchestratorPrivateKeyFetcher

	orchestratorSettings     *settings.OrchestratorSettings
	mockVaultSecretStorage   *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenFetchingOrchestratorPrivateKeyBeforeEach(t *testing.T) *WhenFetchingOrchestratorPrivateKeyTestingSuite {
	mockController := gomock.NewController(t)

	mockVaultSecretStorage := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common_tests.RandomEthAddress(),
	}

	sut := services.NewOrchestratorPrivateKeyFetcher(
		mockVaultSecretStorage,
		orchestratorSettings,
	)

	return &WhenFetchingOrchestratorPrivateKeyTestingSuite{
		sut:                    sut,
		orchestratorSettings:   orchestratorSettings,
		mockVaultSecretStorage: mockVaultSecretStorage,
	}
}

func TestWhenFetchingOrchestratorPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while reading from vault", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFetchingOrchestratorPrivateKeyTestingSuite) {
			expectedPath := fmt.Sprintf("orchestrator/%s/private_key", suite.orchestratorSettings.PublicKey)

			suite.mockVaultSecretStorage.EXPECT().
				Read(expectedPath).
				Return("", assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenFetchingOrchestratorPrivateKeyBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.FetchPrivateKey()

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFetchingOrchestratorPrivateKeyTestingSuite) {
			expectedPath := fmt.Sprintf("orchestrator/%s/private_key", suite.orchestratorSettings.PublicKey)

			suite.mockVaultSecretStorage.EXPECT().
				Read(expectedPath).
				Return("some-private-key", nil)
		}

		t.Run("Should return the private key", func(t *testing.T) {
			t.Parallel()

			suite := WhenFetchingOrchestratorPrivateKeyBeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.FetchPrivateKey()

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "some-private-key", *result)
		})
	})
}
