package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/responses"
)

type ListMCPConnectionsInterface interface {
	List(request *requests.ListMCPConnectionsRequest) (*responses.ListMCPConnectionsResponse, error)
}

type ListMCPConnections struct {
	mcpConnectionRepository repositories.MCPConnectionRepositoryInterface
}

func NewListMCPConnections(mcpConnectionRepository repositories.MCPConnectionRepositoryInterface) *ListMCPConnections {
	return &ListMCPConnections{
		mcpConnectionRepository: mcpConnectionRepository,
	}
}

func (l *ListMCPConnections) List(request *requests.ListMCPConnectionsRequest) (*responses.ListMCPConnectionsResponse, error) {
	err := request.SanitizeAndValidate()
	if err != nil {
		return nil, err
	}

	searchParams := repositories.MCPConnectionSearchParams{
		AgentId:          request.AgentId,
		Search:           request.Search,
		OrderBy:          request.SortBy.ToDatabaseField(),
		OrderByDirection: request.SortDirection.ToDatabaseDirection(),
	}
	if request.Page != nil && request.ElementsPerPage != nil {
		offset := (*request.Page - 1) * (*request.ElementsPerPage)
		limit := *request.ElementsPerPage
		searchParams.Offset = &offset
		searchParams.Limit = &limit
	}

	mcpConnections, err := l.mcpConnectionRepository.Search(searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list MCP connections: %w", err)
	}

	var results []*responses.MCPConnectionResponse
	for _, conn := range mcpConnections {
		response := &responses.MCPConnectionResponse{}
		response.FromModel(conn)

		results = append(results, response)
	}

	return &responses.ListMCPConnectionsResponse{
		Page:            *request.Page,            // Sanitize ensures these are not nil
		ElementsPerPage: *request.ElementsPerPage, // Sanitize ensures these are not nil
		TotalElements:   uint64(len(results)),
		MCPConnections:  results,
	}, nil
}
