package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenStoringMCPConnectionHeaderAuthValueTestingSuite struct {
	sut                *services.MCPConnectionHeaderAuthValueStorage
	vaultSecretStorage *interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenStoringMCPConnectionHeaderAuthValueBeforeEach(t *testing.T) *WhenStoringMCPConnectionHeaderAuthValueTestingSuite {
	mockController := gomock.NewController(t)

	vaultSecretStorageMock := interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewMCPConnectionHeaderAuthValueStorage(vaultSecretStorageMock)

	return &WhenStoringMCPConnectionHeaderAuthValueTestingSuite{
		sut:                sut,
		vaultSecretStorage: vaultSecretStorageMock,
	}
}

func TestWhenStoringMCPConnectionHeaderAuthValue(t *testing.T) {
	t.Parallel()

	t.Run("Store", func(t *testing.T) {
		t.Parallel()

		t.Run("Given an error occurs while storing header auth", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringMCPConnectionHeaderAuthValueBeforeEach(t)

			mcpConnection := model_builders.NewMCPConnectionBuilder().Build()
			headerAuth := map[string]string{"authorization": "Bearer token"}
			expectedPath := fmt.Sprintf(models.MCPAuthHeaderSecretVaultPathTemplate, mcpConnection.AgentID, mcpConnection.Name)

			suite.vaultSecretStorage.EXPECT().
				Store(expectedPath, headerAuth).
				Return(assert.AnError)

			err := suite.sut.Store(mcpConnection, headerAuth)

			assert.Error(t, err)
		})

		t.Run("Given header auth is successfully stored", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringMCPConnectionHeaderAuthValueBeforeEach(t)

			mcpConnection := model_builders.NewMCPConnectionBuilder().Build()
			headerAuth := map[string]string{"authorization": "Bearer token"}
			expectedPath := fmt.Sprintf(models.MCPAuthHeaderSecretVaultPathTemplate, mcpConnection.AgentID, mcpConnection.Name)

			suite.vaultSecretStorage.EXPECT().
				Store(expectedPath, headerAuth).
				Return(nil)

			err := suite.sut.Store(mcpConnection, headerAuth)

			assert.NoError(t, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Parallel()

		t.Run("Given an error occurs while deleting header auth", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringMCPConnectionHeaderAuthValueBeforeEach(t)

			mcpConnection := model_builders.NewMCPConnectionBuilder().Build()
			expectedPath := fmt.Sprintf(models.MCPAuthHeaderSecretVaultPathTemplate, mcpConnection.AgentID, mcpConnection.Name)

			suite.vaultSecretStorage.EXPECT().
				Delete(expectedPath).
				Return(assert.AnError)

			err := suite.sut.Delete(mcpConnection)

			assert.Error(t, err)
		})

		t.Run("Given header auth is successfully deleted", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringMCPConnectionHeaderAuthValueBeforeEach(t)

			mcpConnection := model_builders.NewMCPConnectionBuilder().Build()
			expectedPath := fmt.Sprintf(models.MCPAuthHeaderSecretVaultPathTemplate, mcpConnection.AgentID, mcpConnection.Name)

			suite.vaultSecretStorage.EXPECT().
				Delete(expectedPath).
				Return(nil)

			err := suite.sut.Delete(mcpConnection)

			assert.NoError(t, err)
		})
	})
}
