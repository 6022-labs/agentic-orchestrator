package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/tests/database_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenStoringAgentDatabaseCredentialsTestSuite struct {
	sut                *services.AgentDatabaseCredentialsStorage
	vaultSecretStorage *services_mocks.MockVaultSecretStorageInterface
}

func WhenStoringAgentDatabaseCredentialsBeforeEach(t *testing.T) *WhenStoringAgentDatabaseCredentialsTestSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := services_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentDatabaseCredentialsStorage(vaultSecretStorageMock)

	return &WhenStoringAgentDatabaseCredentialsTestSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenStoringAgentDatabaseCredentials(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while storing credentials", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringAgentDatabaseCredentialsBeforeEach(t)

		id := uuid.New()
		credentials := "test-credentials"
		suite.vaultSecretStorage.EXPECT().Store(gomock.Any(), credentials).Return(assert.AnError)

		err := suite.sut.Store(id, credentials)

		assert.Error(t, err)
	})

	t.Run("Given credentials are successfully stored", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringAgentDatabaseCredentialsBeforeEach(t)

		id := uuid.New()
		credentials := "test-credentials"
		suite.vaultSecretStorage.EXPECT().Store(gomock.Any(), credentials).Return(nil)

		err := suite.sut.Store(id, credentials)

		assert.NoError(t, err)
	})
}
