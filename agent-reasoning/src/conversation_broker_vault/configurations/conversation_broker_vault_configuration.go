package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_vault/settings"
	vault "github.com/hashicorp/vault/api"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddConversationBrokerVaultConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewVaultSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewVaultSecretStorage,
		dig.As(new(interfaces.VaultSecretStorageInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(vaultSettings *settings.VaultSettings, logger *zap.Logger) *vault.Client {
		vaultConfig := vault.DefaultConfig()
		if vaultSettings.Address != "" {
			vaultConfig.Address = vaultSettings.Address
		}

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
