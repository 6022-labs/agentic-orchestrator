package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type ListAgentMCPConnectionsInterface interface {
	List(requests *requests.ListAgentMCPConnectionsRequest) (*responses.ListAgentMCPConnectionsResponse, error)
}

type ListAgentMCPConnections struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface
}

func NewListAgentMCPConnections(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface,
) *ListAgentMCPConnections {
	return &ListAgentMCPConnections{
		agentRegistryRequester: agentRegistryRequester,
		mcpConnectionRequester: mcpConnectionRequester,
	}
}

func (g *ListAgentMCPConnections) List(request *requests.ListAgentMCPConnectionsRequest) (*responses.ListAgentMCPConnectionsResponse, error) {
	err := request.SanitizeAndValidate()
	if err != nil {
		return nil, err
	}

	agent, err := g.agentRegistryRequester.GetAgent(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	if agent == nil {
		return nil, errors.NewValidationError("agentCollectionTokenId", "agent not found")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return nil, errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}
	if agent.Role == enums.AgentRoleHuman {
		return nil, errors.NewValidationError("agentCollectionTokenId", "Agent is human, cannot have MCP settings")
	}

	paginatedMCPConnections, err := g.mcpConnectionRequester.List(&interfaces.ListMCPConnectionsRequest{
		AgentId:         &agent.ID,
		Search:          request.Search,
		Page:            request.Page,
		ElementsPerPage: request.ElementsPerPage,
		SortBy:          request.SortBy.ToDatabaseField(),
		SortDirection:   request.SortDirection.ToDatabaseDirection(),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*responses.MCPConnectionResponse, len(paginatedMCPConnections.MCPConnections))
	for i, mcpConnection := range paginatedMCPConnections.MCPConnections {
		result[i] = &responses.MCPConnectionResponse{
			ID:      mcpConnection.ID,
			AgentID: mcpConnection.AgentID,
			Name:    mcpConnection.Name,
			Url:     mcpConnection.Url,
		}
	}

	return &responses.ListAgentMCPConnectionsResponse{
		Page:            paginatedMCPConnections.Page,
		ElementsPerPage: paginatedMCPConnections.ElementsPerPage,
		TotalElements:   paginatedMCPConnections.TotalElements,
		MCPConnections:  result,
	}, nil
}
