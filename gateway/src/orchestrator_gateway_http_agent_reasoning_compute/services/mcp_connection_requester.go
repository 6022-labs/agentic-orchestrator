package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_reasoning_compute/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_reasoning_compute/models"
	"github.com/google/uuid"
)

type MCPConnectionRequester struct {
	client clients.AgentReasoningComputeClientInterface
}

func NewMCPConnectionRequester(client clients.AgentReasoningComputeClientInterface) *MCPConnectionRequester {
	return &MCPConnectionRequester{
		client: client,
	}
}

func (r *MCPConnectionRequester) Create(request *interfaces.CreateMCPConnectionRequest) error {
	err := r.client.CreateMCPConnection(&models.ExternalCreateMCPConnectionRequest{
		Url:        request.Url,
		Name:       request.Name,
		AgentId:    request.AgentId,
		HeaderAuth: request.HeaderAuth,
	})
	if err != nil {
		return fmt.Errorf("failed to create MCP connection: %w", err)
	}

	return nil
}

func (r *MCPConnectionRequester) Remove(id uuid.UUID) error {
	err := r.client.RemoveMCPConnection(&models.ExternalIdRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("failed to remove MCP connection: %w", err)
	}

	return nil
}

func (r *MCPConnectionRequester) List(request *interfaces.ListMCPConnectionsRequest) (*interfaces.ListMCPConnectionsResponse, error) {
	response, err := r.client.ListMCPConnections(&models.ExternalListMCPConnectionsRequest{
		AgentId:         request.AgentId,
		SortBy:          request.SortBy,
		SortDirection:   request.SortDirection,
		Page:            request.Page,
		ElementsPerPage: request.ElementsPerPage,
		Search:          request.Search,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list MCP connections: %w", err)
	}

	return response.ToInternal(), nil
}
