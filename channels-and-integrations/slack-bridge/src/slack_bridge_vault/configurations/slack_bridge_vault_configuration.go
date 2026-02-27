package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	vault_services "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_vault/services"
	vault_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_vault/settings"
	vault "github.com/hashicorp/vault/api"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddSlackBridgeVaultConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(vault_settings.NewVaultSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		vault_services.NewVaultSecretStorage,
		dig.As(new(interfaces.VaultSecretStorageInterface)),
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
