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

type WhenStoringKeySecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenStoringKeySecretBeforeEach() *WhenStoringKeySecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	vaultSettings := &settings.VaultSettings{
		KvV1MountPath: "secret",
	}

	sut := services.NewVaultSecretStorage(
		vaultClient,
		vaultSettings,
	)

	return &WhenStoringKeySecretTestSuite{
		vaultClient: vaultClient,
		sut:         sut,
	}
}

func WhenStoringKeySecretAfterEach(suite *WhenStoringKeySecretTestSuite) {
	container_builders.CloseVaultConnection(suite.vaultClient)
}

func TestWhenStoringKeySecret(t *testing.T) {
	t.Parallel()

	t.Run("Should store secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenStoringKeySecretBeforeEach()
		defer WhenStoringKeySecretAfterEach(suite)

		err := suite.sut.StoreKey("test-path", "test-key", "test-value")
		assert.Nil(t, err)

		secret, err := suite.vaultClient.KVv1("secret").Get(context.Background(), "test-path")
		assert.Nil(t, err)
		assert.Equal(t, "test-value", secret.Data["test-key"])
	})
}
