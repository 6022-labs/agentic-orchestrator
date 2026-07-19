package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenFetchingOrchestratorPrivateKeyTestingSuite struct {
	sut *services.OrchestratorPrivateKeyFetcher

	vaultSecretStorage   *interfaces_mocks.MockVaultSecretStorageInterface
	orchestratorSettings *settings.OrchestratorSettings
}

func WhenFetchingOrchestratorPrivateKeyBeforeEach(t *testing.T) *WhenFetchingOrchestratorPrivateKeyTestingSuite {
	vaultSecretStorage := interfaces_mocks.NewMockVaultSecretStorageInterface(gomock.NewController(t))
	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common.HexToAddress("0x99065eCD99aE9fBA4dC467C1588b9205Aa71Aa16"),
	}

	return &WhenFetchingOrchestratorPrivateKeyTestingSuite{
		sut:                  services.NewOrchestratorPrivateKeyFetcher(vaultSecretStorage, orchestratorSettings),
		vaultSecretStorage:   vaultSecretStorage,
		orchestratorSettings: orchestratorSettings,
	}
}

func TestWhenFetchingOrchestratorPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given the vault returns the private key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should read it at the orchestrator path", func(t *testing.T) {
			t.Parallel()

			suite := WhenFetchingOrchestratorPrivateKeyBeforeEach(t)

			expectedPrivateKey := "abcdef1234567890"
			expectedPath := fmt.Sprintf("orchestrator/%s/private_key", suite.orchestratorSettings.PublicKey)
			suite.vaultSecretStorage.EXPECT().
				Read(expectedPath).
				Return(expectedPrivateKey, nil)

			result, err := suite.sut.FetchPrivateKey()

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expectedPrivateKey, *result)
		})
	})

	t.Run("Given the vault returns an error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenFetchingOrchestratorPrivateKeyBeforeEach(t)

			suite.vaultSecretStorage.EXPECT().
				Read(gomock.Any()).
				Return("", assert.AnError)

			result, err := suite.sut.FetchPrivateKey()

			assert.Nil(t, result)
			assert.Equal(t, assert.AnError, err)
		})
	})
}
