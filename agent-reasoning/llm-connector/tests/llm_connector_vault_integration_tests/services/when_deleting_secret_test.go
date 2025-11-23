package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_vault/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_vault/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenDeletingSecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenDeletingSecretBeforeEach() *WhenDeletingSecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	vaultSettings := &settings.VaultSettings{
		KvV1MountPath: "secret",
	}

	sut := services.NewVaultSecretStorage(
		vaultClient,
		vaultSettings,
	)

	return &WhenDeletingSecretTestSuite{
		vaultClient: vaultClient,
		sut:         sut,
	}
}

func WhenDeletingSecretAfterEach(suite *WhenDeletingSecretTestSuite) {
	container_builders.CloseVaultConnection(suite.vaultClient)
}

func TestWhenDeletingSecret(t *testing.T) {
	t.Parallel()

	t.Run("Should delete secret when only one key exists", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingSecretBeforeEach()
		defer WhenDeletingSecretAfterEach(suite)

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

	t.Run("Should delete secret when multiple keys exist", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingSecretBeforeEach()
		defer WhenDeletingSecretAfterEach(suite)

		// Store a secret with multiple keys
		err := suite.sut.StoreKey("test-path", "key1", "value1")
		assert.Nil(t, err)
		err = suite.sut.StoreKey("test-path", "key2", "value2")
		assert.Nil(t, err)

		// Delete one key
		err = suite.sut.Delete("test-path")
		assert.Nil(t, err)

		// Try to read the secret, should not exist
		secret, err := suite.vaultClient.KVv1("secret").Get(context.Background(), "test-path")
		assert.NotNil(t, err)
		assert.Nil(t, secret)
	})
}
