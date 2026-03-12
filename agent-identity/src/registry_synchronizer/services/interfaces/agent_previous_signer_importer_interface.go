package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"

type AgentPreviousSignerImporterInterface interface {
	Import(previousSignerGatewayDomain string, agent *models.Agent) error
}
