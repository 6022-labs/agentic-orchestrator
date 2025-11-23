package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenDeletingAgentPrivateKeyTestSuite struct {
	sut                *services.AgentPrivateKeyStorage
	vaultSecretStorage *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenDeletingAgentPrivateKeyBeforeEach(t *testing.T) *WhenDeletingAgentPrivateKeyTestSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentPrivateKeyStorage(vaultSecretStorageMock)

	return &WhenDeletingAgentPrivateKeyTestSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenDeletingAgentPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while deleting the private key", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingAgentPrivateKeyBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Delete(gomock.Any()).Return(assert.AnError)

		err := suite.sut.Delete(id)

		assert.Error(t, err)
	})

	t.Run("Given the private key is successfully deleted", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingAgentPrivateKeyBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Delete(gomock.Any()).Return(nil)

		err := suite.sut.Delete(id)

		assert.NoError(t, err)
	})
}
