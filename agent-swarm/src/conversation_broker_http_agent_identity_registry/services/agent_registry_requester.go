package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/mappings"
)

type AgentRegistryRequester struct {
	client clients.ClientWithResponsesInterface
}

func NewAgentRegistryRequester(client clients.ClientWithResponsesInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{client: client}
}

func (r *AgentRegistryRequester) GetAgent(
	ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
) (*models.RegistryAgent, error) {
	rsp, err := r.client.GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(
		ctx,
		collectionAddress,
		int(collectionTokenId),
	)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode() == http.StatusNotFound {
		return nil, nil
	}
	if rsp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", rsp.StatusCode(), string(rsp.Body))
	}

	return (*mappings.ExternalAgentResponse)(rsp.JSON200).ToRegistryAgent(), nil
}

func (r *AgentRegistryRequester) GetMintedAgentByEnsDomain(
	ctx context.Context,
	ensDomain string,
) (*models.RegistryAgent, error) {
	rsp, err := r.client.GetAgentEnsEnsDomainWithResponse(ctx, ensDomain)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode() == http.StatusNotFound {
		return nil, nil
	}
	if rsp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", rsp.StatusCode(), string(rsp.Body))
	}

	return (*mappings.ExternalAgentResponse)(rsp.JSON200).ToRegistryAgent(), nil
}
