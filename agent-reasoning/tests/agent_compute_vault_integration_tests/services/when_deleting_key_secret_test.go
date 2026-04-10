package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_vault/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_vault/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenDeletingKeySecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenDeletingKeySecretBeforeEach() *WhenDeletingKeySecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	vaultSettings := &settings.VaultSettings{
		KvV1MountPath: "secret",
	}

	sut := services.NewVaultSecretStorage(
		vaultClient,
		vaultSettings,
	)

	return &WhenDeletingKeySecretTestSuite{
		vaultClient: vaultClient,
		sut:         sut,
	}
}

func WhenDeletingKeySecretAfterEach(suite *WhenDeletingKeySecretTestSuite) {
	container_builders.CloseVaultConnection(suite.vaultClient)
}

func TestWhenDeletingKeySecret(t *testing.T) {
	t.Parallel()

	t.Run("Should delete secret when only one key exists", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingKeySecretBeforeEach()
		defer WhenDeletingKeySecretAfterEach(suite)

		// Store a secret with one key
		err := suite.sut.StoreKey("test-path", "test-key", "test-value")
		assert.Nil(t, err)

		// Delete the only key
		err = suite.sut.DeleteKey("test-path", "test-key")
		assert.Nil(t, err)

		// Try to read the secret, should not exist
		secret, err := suite.vaultClient.KVv1("secret").Get(context.Background(), "test-path")
		assert.NotNil(t, err)
		assert.Nil(t, secret)
	})

	t.Run("Should delete only the specified key when multiple keys exist", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingKeySecretBeforeEach()
		defer WhenDeletingKeySecretAfterEach(suite)

		// Store a secret with multiple keys
		err := suite.sut.StoreKey("test-path", "key1", "value1")
		assert.Nil(t, err)
		err = suite.sut.StoreKey("test-path", "key2", "value2")
		assert.Nil(t, err)

		// Delete one key
		err = suite.sut.DeleteKey("test-path", "key1")
		assert.Nil(t, err)

		// Read the secret, key1 should be gone, key2 should remain
		secret, err := suite.vaultClient.KVv1("secret").Get(context.Background(), "test-path")
		assert.Nil(t, err)
		_, exists := secret.Data["key1"]
		assert.False(t, exists)
		val, exists := secret.Data["key2"]
		assert.True(t, exists)
		assert.Equal(t, "value2", val)
	})
}
