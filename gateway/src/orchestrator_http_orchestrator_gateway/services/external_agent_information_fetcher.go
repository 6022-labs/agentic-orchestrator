package services

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_http_orchestrator_gateway/clients"
)

type ExternalAgentInformationFetcher struct {
	orchestratorMainApiClient clients.OrchestratorMainApiClientInterface
}

func NewExternalAgentInformationFetcher(
	orchestratorMainApiClient clients.OrchestratorMainApiClientInterface,
) *ExternalAgentInformationFetcher {
	return &ExternalAgentInformationFetcher{
		orchestratorMainApiClient: orchestratorMainApiClient,
	}
}

func (c *ExternalAgentInformationFetcher) IsInternalAgent(orchestratorDomain string, collectionTokenId uint64) (*bool, error) {
	externalIsInternalAgent, err := c.orchestratorMainApiClient.CallIsInternalAgent(orchestratorDomain, collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &externalIsInternalAgent.IsInternalAgent, nil
}

func (c *ExternalAgentInformationFetcher) ImportAgent(orchestratorDomain string, collectionTokenId uint64) (*interfaces.ImportAgentInformation, error) {
	externalExportAgent, err := c.orchestratorMainApiClient.CallExportAgent(orchestratorDomain, collectionTokenId)
	if err != nil {
		return nil, err
	}

	return externalExportAgent.ToInternalImportAgentInformation(), nil
}

func (c *ExternalAgentInformationFetcher) FetchAgentMetadataInformation(orchestratorDomain string, collectionTokenId uint64) (*interfaces.AgentMetadataInformation, error) {
	externalAgentMetadataInformation, err := c.orchestratorMainApiClient.CallAgentMetadata(orchestratorDomain, collectionTokenId)
	if err != nil {
		return nil, err
	}

	return externalAgentMetadataInformation.ToInternalAgentMetadataInformation(), nil
}
