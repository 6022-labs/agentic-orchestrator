package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentPrivateKeyTestSuite struct {
	sut                *services.AgentPrivateKeyStorage
	vaultSecretStorage *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenGettingAgentPrivateKeyBeforeEach(t *testing.T) *WhenGettingAgentPrivateKeyTestSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentPrivateKeyStorage(vaultSecretStorageMock)

	return &WhenGettingAgentPrivateKeyTestSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenGettingAgentPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while retrieving the private key", func(t *testing.T) {
		t.Parallel()

		suite := WhenGettingAgentPrivateKeyBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Read(gomock.Any()).Return("", assert.AnError)

		privateKey, err := suite.sut.Get(id)

		assert.Nil(t, privateKey)
		assert.Error(t, err)
	})

	t.Run("Given the private key is successfully retrieved", func(t *testing.T) {
		t.Parallel()

		suite := WhenGettingAgentPrivateKeyBeforeEach(t)

		id := uuid.New()
		expectedPrivateKey := "test-private-key"
		suite.vaultSecretStorage.EXPECT().Read(gomock.Any()).Return(expectedPrivateKey, nil)

		privateKey, err := suite.sut.Get(id)

		assert.NoError(t, err)
		assert.Equal(t, expectedPrivateKey, *privateKey)
	})
}
