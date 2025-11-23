package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_vault/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_vault/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/common_tests/container_builders"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type WhenReadingSecretTestSuite struct {
	vaultClient *vault.Client
	sut         *services.VaultSecretStorage
}

func WhenReadingSecretBeforeEach() *WhenReadingSecretTestSuite {
	vaultClient := container_builders.GetVaultConnection()

	vaultSettings := &settings.VaultSettings{
		KvV1MountPath: "secret",
	}

	sut := services.NewVaultSecretStorage(
		vaultClient,
		vaultSettings,
	)

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

		err := suite.sut.Store("test-path", map[string]string{"test-key": "test-value"})
		assert.Nil(t, err)

		secret, err := suite.sut.Read("test-path")
		assert.Nil(t, err)
		assert.Equal(t, "test-value", secret["test-key"])
	})
}
