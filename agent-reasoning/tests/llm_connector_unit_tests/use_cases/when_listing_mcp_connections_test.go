package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingMCPConnectionTestingSuite struct {
	sut *use_cases.ListMCPConnections

	mcpConnectionRepositoryMock *repositories_mocks.MockMCPConnectionRepositoryInterface
}

func WhenListingMCPConnectionBeforeEach(t *testing.T) *WhenListingMCPConnectionTestingSuite {
	mockController := gomock.NewController(t)

	mcpConnectionRepositoryMock := repositories_mocks.NewMockMCPConnectionRepositoryInterface(mockController)

	sut := use_cases.NewListMCPConnections(mcpConnectionRepositoryMock)

	return &WhenListingMCPConnectionTestingSuite{
		sut: sut,

		mcpConnectionRepositoryMock: mcpConnectionRepositoryMock,
	}
}

func generateValidListMCPConnectionsRequest() *requests.ListMCPConnectionsRequest {
	return &requests.ListMCPConnectionsRequest{
		PaginateRequest: requests.PaginateRequest{},
	}
}

func TestWhenListingMCPConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because page is less than 1, should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingMCPConnectionBeforeEach(t)
			request := generateValidListMCPConnectionsRequest()
			invalidPage := uint(0)
			request.Page = &invalidPage

			_, err := suite.sut.List(request)

			assert.Error(t, err)
		})

		t.Run("Because elementsPerPage is less than 1, should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingMCPConnectionBeforeEach(t)
			request := generateValidListMCPConnectionsRequest()
			invalidElementsPerPage := uint(0)
			request.ElementsPerPage = &invalidElementsPerPage

			_, err := suite.sut.List(request)

			assert.Error(t, err)
		})

		t.Run("Because elementsPerPage is greater than 100, should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingMCPConnectionBeforeEach(t)
			request := generateValidListMCPConnectionsRequest()
			invalidElementsPerPage := uint(101)
			request.ElementsPerPage = &invalidElementsPerPage

			_, err := suite.sut.List(request)

			assert.Error(t, err)
		})

		t.Run("Because sort direction is invalid, should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingMCPConnectionBeforeEach(t)
			request := generateValidListMCPConnectionsRequest()
			invalidSortDirection := requests.SortDirection("invalid_direction")
			request.SortDirection = &invalidSortDirection

			_, err := suite.sut.List(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while searching MCP connections in repository", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingMCPConnectionBeforeEach(t)
			agentId := uuid.New()
			request := generateValidListMCPConnectionsRequest()
			request.AgentId = &agentId

			suite.mcpConnectionRepositoryMock.
				EXPECT().
				Search(gomock.Any()).
				Return(nil, assert.AnError)

			result, err := suite.sut.List(request)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Given no MCP connections found in repository", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return empty result", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingMCPConnectionBeforeEach(t)
			agentId := uuid.New()
			request := generateValidListMCPConnectionsRequest()
			request.AgentId = &agentId

			suite.mcpConnectionRepositoryMock.
				EXPECT().
				Search(gomock.Any()).
				Return([]*models.MCPConnection{}, nil)

			response, err := suite.sut.List(request)

			assert.NoError(t, err)
			assert.Empty(t, response.MCPConnections)
		})
	})

	t.Run("Given MCP connections found in repository", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return mapped results", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingMCPConnectionBeforeEach(t)
			agentId := uuid.New()
			request := generateValidListMCPConnectionsRequest()
			request.AgentId = &agentId

			conn := model_builders.NewMCPConnectionBuilder().Build()
			suite.mcpConnectionRepositoryMock.
				EXPECT().
				Search(gomock.Any()).
				Return([]*models.MCPConnection{conn}, nil)

			response, err := suite.sut.List(request)

			assert.NoError(t, err)
			assert.Len(t, response.MCPConnections, 1)
			assert.Equal(t, conn.ID, response.MCPConnections[0].ID)
			assert.Equal(t, conn.Url, response.MCPConnections[0].Url)
			assert.Equal(t, conn.Name, response.MCPConnections[0].Name)
		})
	})
}
