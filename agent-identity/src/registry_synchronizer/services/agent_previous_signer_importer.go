package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
)

type AgentPreviousSignerImporter struct {
}

func NewAgentPreviousSignerImporter() *AgentPreviousSignerImporter {
	return &AgentPreviousSignerImporter{}
}

func (a *AgentPreviousSignerImporter) Import(previousSignerGatewayDomain string, agent *models.Agent) error {
	return fmt.Errorf("not implemented")
}
