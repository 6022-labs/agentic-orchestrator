package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenStoringAgentLLMSecretTestingSuite struct {
	sut                *services.AgentLLMSecretStorage
	vaultSecretStorage *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenStoringAgentLLMSecretBeforeEach(t *testing.T) *WhenStoringAgentLLMSecretTestingSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentLLMSecretStorage(vaultSecretStorageMock)

	return &WhenStoringAgentLLMSecretTestingSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenStoringAgentLLMSecret(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while storing secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringAgentLLMSecretBeforeEach(t)

		id := uuid.New()
		secret := "test-secret"
		suite.vaultSecretStorage.EXPECT().StoreKey(gomock.Any(), "value", secret).Return(assert.AnError)

		err := suite.sut.Store(id, secret)

		assert.Error(t, err)
	})

	t.Run("Given secret is successfully stored", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringAgentLLMSecretBeforeEach(t)

		id := uuid.New()
		secret := "test-secret"
		suite.vaultSecretStorage.EXPECT().StoreKey(gomock.Any(), "value", secret).Return(nil)

		err := suite.sut.Store(id, secret)

		assert.NoError(t, err)
	})
}
