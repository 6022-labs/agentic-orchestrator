package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenFetchingPrivateKeyTestSuite struct {
	sut *services.OrchestratorPrivateKeyFetcher

	orchestratorSettings   *settings.OrchestratorSettings
	vaultSecretStorageMock *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenFetchingPrivateKeyBeforeEach(t *testing.T) WhenFetchingPrivateKeyTestSuite {
	mockController := gomock.NewController(t)

	orchestratorSettings := &settings.OrchestratorSettings{}
	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewOrchestratorPrivateKeyFetcher(
		vaultSecretStorageMock,
		orchestratorSettings,
	)

	return WhenFetchingPrivateKeyTestSuite{
		sut: sut,

		orchestratorSettings:   orchestratorSettings,
		vaultSecretStorageMock: vaultSecretStorageMock,
	}
}

func TestWhenFetchingPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given error while requesting persistent storage", func(t *testing.T) {
		t.Parallel()

		suite := WhenFetchingPrivateKeyBeforeEach(t)

		suite.vaultSecretStorageMock.EXPECT().Read(gomock.Any()).Return("", assert.AnError)

		privateKey, err := suite.sut.FetchPrivateKey()

		assert.Nil(t, privateKey)
		assert.NotNil(t, err)
	})

	t.Run("Given successful request to persistent storage", func(t *testing.T) {
		t.Parallel()

		suite := WhenFetchingPrivateKeyBeforeEach(t)

		expectedPrivateKey := "private"
		suite.vaultSecretStorageMock.EXPECT().Read(gomock.Any()).Return(expectedPrivateKey, nil)

		privateKey, err := suite.sut.FetchPrivateKey()

		assert.Equal(t, expectedPrivateKey, *privateKey)
		assert.Nil(t, err)
	})
}
