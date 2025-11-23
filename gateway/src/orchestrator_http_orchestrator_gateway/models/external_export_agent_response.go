package models

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"

type ExternalExportAgentResponse struct {
	PrivateKey string `json:"privateKey"`
}

func (m *ExternalExportAgentResponse) ToInternalImportAgentInformation() *interfaces.ImportAgentInformation {
	return &interfaces.ImportAgentInformation{
		PrivateKey: m.PrivateKey,
	}
}
