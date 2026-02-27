package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/services"
	vault_services "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_vault/services"
	vault_settings "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_vault/settings"
	vault "github.com/hashicorp/vault/api"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddDatabaseConnectorVaultConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(vault_settings.NewVaultSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		vault_services.NewVaultSecretStorage,
		dig.As(new(services.VaultSecretStorageInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Provide the vault client
	err = container.Provide(func(vaultSettings *vault_settings.VaultSettings, logger *zap.Logger) *vault.Client {
		vaultConfig := vault.DefaultConfig()
		vaultConfig.Address = vaultSettings.Address

		client, err := vault.NewClient(vaultConfig)
		if err != nil {
			logger.Fatal("Failed to create vault client", zap.Error(err))
		}

		client.SetToken(vaultSettings.Token)

		return client
	})
	if err != nil {
		panic(err)
	}
}
