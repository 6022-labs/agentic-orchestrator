package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenDeletingAgentDatabaseCredentialsTestSuite struct {
	sut                *services.AgentDatabaseCredentialsStorage
	vaultSecretStorage *services_mocks.MockVaultSecretStorageInterface
}

func WhenDeletingAgentDatabaseCredentialsBeforeEach(t *testing.T) *WhenDeletingAgentDatabaseCredentialsTestSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := services_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentDatabaseCredentialsStorage(vaultSecretStorageMock)

	return &WhenDeletingAgentDatabaseCredentialsTestSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenDeletingAgentDatabaseCredentials(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while deleting credentials", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingAgentDatabaseCredentialsBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Delete(gomock.Any()).Return(assert.AnError)

		err := suite.sut.Delete(id)

		assert.Error(t, err)
	})

	t.Run("Given credentials are successfully deleted", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingAgentDatabaseCredentialsBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Delete(gomock.Any()).Return(nil)

		err := suite.sut.Delete(id)

		assert.NoError(t, err)
	})
}
