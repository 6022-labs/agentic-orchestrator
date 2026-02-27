package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_vault/services"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_vault/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenDeletingSecretTestingSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenDeletingSecretBeforeEach() *WhenDeletingSecretTestingSuite {
	vaultClient := container_builders.GetVaultConnection()

	vaultSettings := &settings.VaultSettings{
		KvV1MountPath: "secret",
	}

	sut := services.NewVaultSecretStorage(
		vaultClient,
		vaultSettings,
	)

	return &WhenDeletingSecretTestingSuite{
		vaultClient: vaultClient,
		sut:         sut,
	}
}

func WhenDeletingSecretAfterEach(suite *WhenDeletingSecretTestingSuite) {
	container_builders.CloseVaultConnection(suite.vaultClient)
}

func TestWhenDeletingSecret(t *testing.T) {
	t.Parallel()

	t.Run("Given secret does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingSecretBeforeEach()
			defer WhenDeletingSecretAfterEach(suite)

			err := suite.sut.Delete("non-existing-key")
			assert.Nil(t, err)
		})
	})

	t.Run("Given secret exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should delete secret", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingSecretBeforeEach()
			defer WhenDeletingSecretAfterEach(suite)

			err := suite.sut.Store("test-key", "test-value")
			assert.Nil(t, err)

			err = suite.sut.Delete("test-key")
			assert.Nil(t, err)

			secret, err := suite.sut.Read("test-key")
			assert.NotNil(t, err)
			assert.Equal(t, "", secret)
		})
	})
}
