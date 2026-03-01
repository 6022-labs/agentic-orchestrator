package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenDeletingAgentLLMSecretTestingSuite struct {
	sut                *services.AgentLLMSecretStorage
	vaultSecretStorage *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenDeletingAgentLLMSecretBeforeEach(t *testing.T) *WhenDeletingAgentLLMSecretTestingSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentLLMSecretStorage(vaultSecretStorageMock)

	return &WhenDeletingAgentLLMSecretTestingSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenDeletingAgentLLMSecret(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while deleting secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingAgentLLMSecretBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Delete(gomock.Any()).Return(assert.AnError)

		err := suite.sut.Delete(id)

		assert.Error(t, err)
	})

	t.Run("Given secret is successfully deleted", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingAgentLLMSecretBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Delete(gomock.Any()).Return(nil)

		err := suite.sut.Delete(id)

		assert.NoError(t, err)
	})
}
