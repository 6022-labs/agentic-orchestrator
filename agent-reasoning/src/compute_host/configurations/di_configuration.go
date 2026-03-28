package configurations

import (
	orchestrator_agent_compute_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/configurations"
	orchestrator_agent_compute_db_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/configurations"
	orchestrator_agent_compute_host_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_host/settings"
	orchestrator_agent_compute_http_agent_identity_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_identity/configurations"
	orchestrator_agent_compute_http_agent_knowledge_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/configurations"
	orchestrator_agent_compute_http_inference_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_inference/configurations"
	orchestrator_agent_compute_http_mcp_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_mcp/configurations"
	orchestrator_agent_compute_mvc_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_mvc/configurations"
	orchestrator_agent_compute_vault_configuration "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_vault/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(orchestrator_agent_compute_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	orchestrator_agent_compute_configuration.AddComputeConfiguration(container)
	orchestrator_agent_compute_db_configuration.AddComputeDbConfiguration(container)
	orchestrator_agent_compute_mvc_configuration.AddComputeMvcConfiguration(container)
	orchestrator_agent_compute_vault_configuration.AddComputeVaultConfiguration(container)
	orchestrator_agent_compute_http_mcp_configuration.AddComputeHttpMcpConfiguration(container)
	orchestrator_agent_compute_http_inference_configuration.AddComputeHttpInferenceConfiguration(container)
	orchestrator_agent_compute_http_agent_knowledge_configuration.AddComputeHttpAgentKnowledgeConfiguration(container)
	orchestrator_agent_compute_http_agent_identity_configuration.AddComputeHttpAgentIdentityConfiguration(container)

	return container
}
