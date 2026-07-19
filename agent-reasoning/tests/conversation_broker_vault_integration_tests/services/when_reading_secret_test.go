package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenReadingSecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenReadingSecretBeforeEach() *WhenReadingSecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	sut := services.NewVaultSecretStorage(vaultClient, &settings.VaultSettings{
		KvV1MountPath: "secret",
	})

	return &WhenReadingSecretTestSuite{
		vaultClient: vaultClient,
		sut:         sut,
	}
}

func WhenReadingSecretAfterEach(suite *WhenReadingSecretTestSuite) {
	container_builders.CloseVaultConnection(suite.vaultClient)
}

func TestWhenReadingSecret(t *testing.T) {
	t.Parallel()

	t.Run("Should read secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenReadingSecretBeforeEach()
		defer WhenReadingSecretAfterEach(suite)

		err := suite.sut.Store("read-test-key", "test-value")
		assert.Nil(t, err)

		secret, err := suite.sut.Read("read-test-key")
		assert.Nil(t, err)
		assert.Equal(t, "test-value", secret)
	})
}
