package services

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_identity_registry/clients"
	external_models "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_identity_registry/models"
)

type AgentRegistryRequester struct {
	client clients.AgentIdentityRegistryClientInterface
}

func NewAgentRegistryRequester(client clients.AgentIdentityRegistryClientInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{
		client: client,
	}
}

func (r *AgentRegistryRequester) GetAgent(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error) {
	agent, err := r.client.GetAgent(&external_models.ExternalAgentByCollectionTokenIdRequest{
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
