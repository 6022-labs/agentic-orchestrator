package configurations

import (
	orchestrator_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/configurations"
	orchestrator_blockchain_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_blockchain/configurations"
	orchestrator_db_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/configurations"
	orchestrator_http_orchestrator_gateway_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_http_orchestrator_gateway/configurations"
	orchestrator_sc_sync_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/configurations"
	orchestrator_sc_sync_event_listeners_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync_event_listeners/configurations"
	orchestrator_sc_sync_host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync_host/settings"
	orchestrator_sc_sync_mvc_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync_mvc/configurations"
	orchestrator_vault_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_vault/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(orchestrator_sc_sync_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	orchestrator_configurations.AddOrchestratorConfiguration(container)
	orchestrator_db_configurations.AddOrchestratorDbConfiguration(container)
	orchestrator_vault_configurations.AddOrchestratorVaultConfiguration(container)
	orchestrator_sc_sync_configurations.AddOrchestratorSCSyncConfiguration(container)
	orchestrator_blockchain_configurations.AddOrchestratorBlockchainConfiguration(container)
	orchestrator_sc_sync_mvc_configurations.AddOrchestratorScSyncMvcConfiguration(container)
	orchestrator_sc_sync_event_listeners_configurations.AddOrchestratorScSyncEventListenersConfiguration(container)
	orchestrator_http_orchestrator_gateway_configurations.AddOrchestratorHttpOrchestratorGatewayConfiguration(container)

	return container
}
