package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenDeletingSecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenDeletingSecretBeforeEach() *WhenDeletingSecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	sut := services.NewVaultSecretStorage(vaultClient, &settings.VaultSettings{
		KvV1MountPath: "secret",
	})

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

	t.Run("Should delete secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenDeletingSecretBeforeEach()
		defer WhenDeletingSecretAfterEach(suite)

		err := suite.sut.Store("delete-test-key", "test-value")
		assert.Nil(t, err)

		err = suite.sut.Delete("delete-test-key")
		assert.Nil(t, err)

		_, err = suite.sut.Read("delete-test-key")
		assert.Error(t, err)
	})
}
