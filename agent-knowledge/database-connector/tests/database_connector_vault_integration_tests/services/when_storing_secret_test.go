package services_test

import (
	"context"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_vault/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_vault/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/container_builders"
)

type WhenStoringSecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenStoringSecretBeforeEach() *WhenStoringSecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	vaultSettings := &settings.VaultSettings{
		KvV1MountPath: "secret",
	}

	sut := services.NewVaultSecretStorage(
		vaultClient,
		vaultSettings,
	)

	return &WhenStoringSecretTestSuite{
		vaultClient: vaultClient,
		sut:         sut,
	}
}

func WhenStoringSecretAfterEach(suite *WhenStoringSecretTestSuite) {
	container_builders.CloseVaultConnection(suite.vaultClient)
}

func TestWhenStoringSecret(t *testing.T) {
	t.Parallel()

	t.Run("Should store secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringSecretBeforeEach()
		defer WhenStoringSecretAfterEach(suite)

		err := suite.sut.Store("test-key", "test-value")
		assert.Nil(t, err)

		secret, err := suite.vaultClient.KVv1("secret").Get(context.Background(), "test-key")
		assert.Nil(t, err)
		assert.Equal(t, "test-value", secret.Data["value"])
	})
}
