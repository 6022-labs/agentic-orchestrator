package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/mappings"
	"github.com/samber/lo"
)

type AgentRegistryRequester struct {
	client clients.ClientWithResponsesInterface
}

func NewAgentRegistryRequester(client clients.ClientWithResponsesInterface) *AgentRegistryRequester {
	return &AgentRegistryRequester{client: client}
}

func (r *AgentRegistryRequester) ListAgents(request *interfaces.ListRegistryAgentsRequest) (*responses.ListAgentsResponse, error) {
	params := &clients.GetAgentsParams{
		SortBy:        request.SortBy,
		SortDirection: request.SortDirection,
		IsMinted:      request.IsMinted,
	}
	if request.Page != nil {
		params.Page = lo.ToPtr(int(*request.Page))
	}
	if request.ElementsPerPage != nil {
		params.ElementsPerPage = lo.ToPtr(int(*request.ElementsPerPage))
	}

	rsp, err := r.client.GetAgentsWithResponse(context.Background(), params)
	if err != nil {
		return nil, err
	}
	if rsp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", rsp.StatusCode(), string(rsp.Body))
	}

	return (*mappings.ExternalListAgentsResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}

func (r *AgentRegistryRequester) GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*responses.AgentResponse, error) {
	rsp, err := r.client.GetAgentAgentCollectionAddressAgentCollectionTokenIdWithResponse(
		context.Background(),
		agentCollectionAddress,
		int(agentCollectionTokenId),
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

	return (*mappings.ExternalAgentResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}

func (r *AgentRegistryRequester) ListAgentCollections() ([]*responses.AgentCollectionResponse, error) {
	rsp, err := r.client.GetAgentCollectionsWithResponse(context.Background())
	if err != nil {
		return nil, err
	}
	if rsp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", rsp.StatusCode(), string(rsp.Body))
	}

	collections := make([]*responses.AgentCollectionResponse, len(*rsp.JSON200))
	for i := range *rsp.JSON200 {
		external := mappings.ExternalAgentCollectionResponse((*rsp.JSON200)[i])
		collections[i] = external.ToUseCaseResponse()
	}
	return collections, nil
}
