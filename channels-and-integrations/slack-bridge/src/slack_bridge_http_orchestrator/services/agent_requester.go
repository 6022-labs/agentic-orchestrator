package services

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/clients"
	external_models "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/models"
)

type AgentRequester struct {
	client clients.OrchestratorGatewayClientInterface
}

func NewAgentRequester(client clients.OrchestratorGatewayClientInterface) *AgentRequester {
	return &AgentRequester{
		client: client,
	}
}

func (a *AgentRequester) GetAgent(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error) {
	agent, err := a.client.GetAgent(&external_models.ExternalAgentByCollectionTokenIdRequest{
		AgentCollectionAddress: collectionAddress,
		AgentCollectionTokenId: collectionTokenId,
	})
	if err != nil {
		return nil, err
	}

	slackAgent, err := agent.ToSlackAgent()
	if err != nil {
		return nil, err
	}

	return slackAgent, nil
}
