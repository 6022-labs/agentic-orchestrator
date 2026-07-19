package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenStoringSecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenStoringSecretBeforeEach() *WhenStoringSecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	sut := services.NewVaultSecretStorage(vaultClient, &settings.VaultSettings{
		KvV1MountPath: "secret",
	})

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

		err := suite.sut.Store("store-test-key", "test-value")
		assert.Nil(t, err)

		secret, err := suite.vaultClient.KVv1("secret").Get(context.Background(), "store-test-key")
		assert.Nil(t, err)
		assert.Equal(t, "test-value", secret.Data["value"])
	})
}
