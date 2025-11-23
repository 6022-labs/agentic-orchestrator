package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_vault/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_vault/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenReadingKeySecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenReadingKeySecretBeforeEach() *WhenReadingKeySecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	vaultSettings := &settings.VaultSettings{
		KvV1MountPath: "secret",
	}

	sut := services.NewVaultSecretStorage(
		vaultClient,
		vaultSettings,
	)

	return &WhenReadingKeySecretTestSuite{
		vaultClient: vaultClient,
		sut:         sut,
	}
}

func WhenReadingKeySecretAfterEach(suite *WhenReadingKeySecretTestSuite) {
	container_builders.CloseVaultConnection(suite.vaultClient)
}

func TestWhenReadingKeySecret(t *testing.T) {
	t.Parallel()

	t.Run("Should read secret", func(t *testing.T) {
		t.Parallel()

		suite := WhenReadingKeySecretBeforeEach()
		defer WhenReadingKeySecretAfterEach(suite)

		err := suite.sut.StoreKey("test-path", "test-key", "test-value")
		assert.Nil(t, err)

		secret, err := suite.sut.ReadKey("test-path", "test-key")
		assert.Nil(t, err)
		assert.Equal(t, "test-value", secret)
	})
}
