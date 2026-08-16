package services

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/mappings"
)

type AgentRegistryRequester struct {
	client clients.ClientWithResponsesInterface
}

func NewAgentRegistryRequester(client clients.ClientWithResponsesInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{
		client: client,
	}
}

func (r *AgentRegistryRequester) GetAgent(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error) {
	rsp, err := r.client.GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(
		context.Background(),
		collectionAddress,
		int(collectionTokenId),
	)
	if err != nil {
		return nil, err
	}
	if rsp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", rsp.StatusCode(), string(rsp.Body))
	}

	return (*mappings.ExternalAgentResponse)(rsp.JSON200).ToSlackAgent()
}
