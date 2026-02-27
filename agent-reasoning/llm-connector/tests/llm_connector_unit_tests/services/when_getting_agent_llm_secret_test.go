package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/tests/llm_connector_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentLLMSecretTestingSuite struct {
	sut                *services.AgentLLMSecretStorage
	vaultSecretStorage *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenGettingAgentLLMSecretBeforeEach(t *testing.T) *WhenGettingAgentLLMSecretTestingSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentLLMSecretStorage(vaultSecretStorageMock)

	return &WhenGettingAgentLLMSecretTestingSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenGettingAgentLLMSecret(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while retrieving secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenGettingAgentLLMSecretBeforeEach(t)

		id := uuid.New()
		suite.vaultSecretStorage.EXPECT().ReadKey(gomock.Any(), "value").Return("", assert.AnError)

		credentials, err := suite.sut.Get(id)

		assert.Nil(t, credentials)
		assert.Error(t, err)
	})

	t.Run("Given secret is successfully retrieved", func(t *testing.T) {
		t.Parallel()

		suite := WhenGettingAgentLLMSecretBeforeEach(t)

		id := uuid.New()
		expectedCredentials := "test-credentials"
		suite.vaultSecretStorage.EXPECT().ReadKey(gomock.Any(), "value").Return(expectedCredentials, nil)

		credentials, err := suite.sut.Get(id)

		assert.NoError(t, err)
		assert.Equal(t, expectedCredentials, *credentials)
	})
}
