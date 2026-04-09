package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	common_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenFetchingOrchestratorPrivateKeyTestingSuite struct {
	sut *services.OrchestratorPrivateKeyFetcher

	mockVaultSecretStorage *common_interfaces_mocks.MockVaultSecretStorageInterface
	orchestratorSettings   *settings.OrchestratorSettings
}

func WhenFetchingOrchestratorPrivateKeyBeforeEach(t *testing.T) *WhenFetchingOrchestratorPrivateKeyTestingSuite {
	mockController := gomock.NewController(t)

	mockVaultSecretStorage := common_interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)
	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common_tests.RandomEthAddress(),
	}

	sut := services.NewOrchestratorPrivateKeyFetcher(mockVaultSecretStorage, orchestratorSettings)

	return &WhenFetchingOrchestratorPrivateKeyTestingSuite{
		sut:                    sut,
		mockVaultSecretStorage: mockVaultSecretStorage,
		orchestratorSettings:   orchestratorSettings,
	}
}

func TestWhenFetchingOrchestratorPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given the vault returns the private key successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the orchestrator private key", func(t *testing.T) {
			t.Parallel()
			basics := WhenFetchingOrchestratorPrivateKeyBeforeEach(t)

			expectedPrivateKey := "abcdef1234567890"
			expectedPath := fmt.Sprintf("orchestrator/%s/private_key", basics.orchestratorSettings.PublicKey)

			basics.mockVaultSecretStorage.EXPECT().
				Read(expectedPath).
				Return(expectedPrivateKey, nil).
				Times(1)

			result, err := basics.sut.FetchPrivateKey()

			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expectedPrivateKey, *result)
		})
	})

	t.Run("Given the vault returns an error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()
			basics := WhenFetchingOrchestratorPrivateKeyBeforeEach(t)

			basics.mockVaultSecretStorage.EXPECT().
				Read(gomock.Any()).
				Return("", assert.AnError).
				Times(1)

			result, err := basics.sut.FetchPrivateKey()

			assert.Nil(t, result)
			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})
}
