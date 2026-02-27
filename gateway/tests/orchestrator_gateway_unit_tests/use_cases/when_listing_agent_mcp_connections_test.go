package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentMCPConnectionsTestingSuite struct {
	sut *use_cases.ListAgentMCPConnections

	agentRepositoryMock        *repositories_mocks.MockAgentRepositoryInterface
	mcpConnectionRequesterMock *interfaces_mocks.MockMCPConnectionRequesterInterface
}

func WhenListingAgentMCPConnectionsBeforeEach(t *testing.T) *WhenListingAgentMCPConnectionsTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mcpConnectionRequesterMock := interfaces_mocks.NewMockMCPConnectionRequesterInterface(mockController)

	sut := use_cases.NewListAgentMCPConnections(
		agentRepositoryMock,
		mcpConnectionRequesterMock,
	)

	return &WhenListingAgentMCPConnectionsTestingSuite{
		sut: sut,

		agentRepositoryMock:        agentRepositoryMock,
		mcpConnectionRequesterMock: mcpConnectionRequesterMock,
	}
}

func generateValidListAgentMCPConnectionsRequest() *requests.ListAgentMCPConnectionsRequest {
	authenticatedWallet := common_tests.RandomEthAddress()
	collectionAddress := common_tests.RandomEthAddress()

	return &requests.ListAgentMCPConnectionsRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionTokenId: 1,
			AgentCollectionAddress: collectionAddress.Hex(),
			AuthenticatedWallet:    &authenticatedWallet,
		},
	}
}

func TestWhenListingAgentMCPConnections(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because authenticated wallet is missing", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentMCPConnectionsBeforeEach(t)
			req := generateValidListAgentMCPConnectionsRequest()
			req.AuthenticatedWallet = nil

			connections, err := suite.sut.List(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Nil(t, connections)
		})
	})

	t.Run("Given error occurs while getting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentMCPConnectionsTestingSuite, request *requests.ListAgentMCPConnectionsRequest) {
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentMCPConnectionsBeforeEach(t)
			req := generateValidListAgentMCPConnectionsRequest()

			initSuite(suite, req)

			connections, err := suite.sut.List(req)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent")
			assert.Nil(t, connections)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentMCPConnectionsTestingSuite, request *requests.ListAgentMCPConnectionsRequest) {
			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentMCPConnectionsBeforeEach(t)
			req := generateValidListAgentMCPConnectionsRequest()

			initSuite(suite, req)

			connections, err := suite.sut.List(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Nil(t, connections)
		})
	})

	t.Run("Given caller is not the owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentMCPConnectionsTestingSuite, request *requests.ListAgentMCPConnectionsRequest) {
			otherOwner := common_tests.RandomEthAddress().Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(otherOwner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentMCPConnectionsBeforeEach(t)
			req := generateValidListAgentMCPConnectionsRequest()

			initSuite(suite, req)

			connections, err := suite.sut.List(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Nil(t, connections)
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentMCPConnectionsTestingSuite, request *requests.ListAgentMCPConnectionsRequest) {
			owner := request.AuthenticatedWallet.Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				WithRole(enums.AgentRoleHuman).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentMCPConnectionsBeforeEach(t)
			req := generateValidListAgentMCPConnectionsRequest()

			initSuite(suite, req)

			connections, err := suite.sut.List(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Nil(t, connections)
		})
	})

	t.Run("Given error occurs while listing mcp connections", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentMCPConnectionsTestingSuite, request *requests.ListAgentMCPConnectionsRequest) {
			owner := request.AuthenticatedWallet.Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mcpConnectionRequesterMock.
				EXPECT().
				List(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentMCPConnectionsBeforeEach(t)
			req := generateValidListAgentMCPConnectionsRequest()

			initSuite(suite, req)

			connections, err := suite.sut.List(req)
			assert.Error(t, err)
			assert.Equal(t, assert.AnError, err)
			assert.Nil(t, connections)
		})
	})

	t.Run("Given all steps are successful", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentMCPConnectionsTestingSuite, request *requests.ListAgentMCPConnectionsRequest) *interfaces.ListMCPConnectionsResponse {
			owner := request.AuthenticatedWallet.Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				Build()

			suite.agentRepositoryMock.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			expected := &interfaces.ListMCPConnectionsResponse{
				Page:            1,
				ElementsPerPage: 10,
				TotalElements:   2,
				MCPConnections: []*interfaces.MCPConnection{
					{ID: "1", AgentID: agent.ID.String(), Name: "conn-1"},
					{ID: "2", AgentID: agent.ID.String(), Name: "conn-2"},
				},
			}

			suite.mcpConnectionRequesterMock.
				EXPECT().
				List(gomock.Any()).
				Return(expected, nil)

			return expected
		}

		t.Run("Should return mapped connections", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentMCPConnectionsBeforeEach(t)
			req := generateValidListAgentMCPConnectionsRequest()

			expected := initSuite(suite, req)

			response, err := suite.sut.List(req)

			assert.NoError(t, err)
			if assert.Len(t, response.MCPConnections, 2) {
				assert.Equal(t, expected.MCPConnections[0].ID, response.MCPConnections[0].ID)
				assert.Equal(t, expected.MCPConnections[0].AgentID, response.MCPConnections[0].AgentID)
				assert.Equal(t, expected.MCPConnections[0].Name, response.MCPConnections[0].Name)
				assert.Equal(t, expected.MCPConnections[1].ID, response.MCPConnections[1].ID)
				assert.Equal(t, expected.MCPConnections[1].AgentID, response.MCPConnections[1].AgentID)
				assert.Equal(t, expected.MCPConnections[1].Name, response.MCPConnections[1].Name)
			}
		})
	})
}
