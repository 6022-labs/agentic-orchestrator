package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentDatabaseCredentialsTestSuite struct {
	sut                *services.AgentDatabaseCredentialsStorage
	vaultSecretStorage *services_mocks.MockVaultSecretStorageInterface
}

func WhenGettingAgentDatabaseCredentialsBeforeEach(t *testing.T) *WhenGettingAgentDatabaseCredentialsTestSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := services_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentDatabaseCredentialsStorage(vaultSecretStorageMock)

	return &WhenGettingAgentDatabaseCredentialsTestSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenGettingAgentDatabaseCredentials(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while retrieving credentials", func(t *testing.T) {
		t.Parallel()

		suite := WhenGettingAgentDatabaseCredentialsBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().Read(gomock.Any()).Return("", assert.AnError)

		credentials, err := suite.sut.Get(id)

		assert.Nil(t, credentials)
		assert.Error(t, err)
	})

	t.Run("Given credentials are successfully retrieved", func(t *testing.T) {
		t.Parallel()

		suite := WhenGettingAgentDatabaseCredentialsBeforeEach(t)

		id := uuid.New()
		expectedCredentials := "test-credentials"
		suite.vaultSecretStorage.EXPECT().Read(gomock.Any()).Return(expectedCredentials, nil)

		credentials, err := suite.sut.Get(id)

		assert.NoError(t, err)
		assert.Equal(t, expectedCredentials, *credentials)
	})
}
