package configurations

import (
	orchestrator_agent_llm_connector_configuration "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/configurations"
	orchestrator_agent_llm_connector_db_configuration "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_db/configurations"
	orchestrator_agent_llm_connector_host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_host/settings"
	orchestrator_agent_llm_connector_http_inference_configuration "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/configurations"
	orchestrator_agent_llm_connector_http_mcp_configuration "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_mcp/configurations"
	orchestrator_agent_llm_connector_mvc_configuration "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_mvc/configurations"
	orchestrator_agent_llm_connector_vault_configuration "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_vault/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(orchestrator_agent_llm_connector_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	orchestrator_agent_llm_connector_configuration.AddLLMConnectorConfiguration(container)
	orchestrator_agent_llm_connector_db_configuration.AddLLMConnectorDbConfiguration(container)
	orchestrator_agent_llm_connector_mvc_configuration.AddLLMConnectorMvcConfiguration(container)
	orchestrator_agent_llm_connector_vault_configuration.AddLLMConnectorVaultConfiguration(container)
	orchestrator_agent_llm_connector_http_mcp_configuration.AddLLMConnectorHttpMcpConfiguration(container)
	orchestrator_agent_llm_connector_http_inference_configuration.AddLLMConnectorHttpInferenceConfiguration(container)

	return container
}
