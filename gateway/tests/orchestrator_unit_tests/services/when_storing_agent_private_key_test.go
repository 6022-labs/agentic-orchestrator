package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenStoringAgentPrivateKeyTestSuite struct {
	sut                *services.AgentPrivateKeyStorage
	vaultSecretStorage *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenStoringAgentPrivateKeyBeforeEach(t *testing.T) *WhenStoringAgentPrivateKeyTestSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentPrivateKeyStorage(vaultSecretStorageMock)

	return &WhenStoringAgentPrivateKeyTestSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenStoringAgentPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while storing the private key", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringAgentPrivateKeyBeforeEach(t)

		id := uuid.New()
		privateKey := "test-private-key"
		suite.vaultSecretStorage.EXPECT().Store(gomock.Any(), privateKey).Return(assert.AnError)

		err := suite.sut.Store(id, privateKey)

		assert.Error(t, err)
	})

	t.Run("Given the private key is successfully stored", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringAgentPrivateKeyBeforeEach(t)

		id := uuid.New()
		privateKey := "test-private-key"
		suite.vaultSecretStorage.EXPECT().Store(gomock.Any(), privateKey).Return(nil)

		err := suite.sut.Store(id, privateKey)

		assert.NoError(t, err)
	})
}
