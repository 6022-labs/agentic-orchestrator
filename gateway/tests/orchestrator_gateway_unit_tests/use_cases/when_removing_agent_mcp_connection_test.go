package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRemovingAgentMCPConnectionTestingSuite struct {
	sut *use_cases.RemoveAgentMCPConnection

	agentRepositoryMock        *repositories_mocks.MockAgentRepositoryInterface
	mcpConnectionRequesterMock *interfaces_mocks.MockMCPConnectionRequesterInterface
}

func WhenRemovingAgentMCPConnectionBeforeEach(t *testing.T) *WhenRemovingAgentMCPConnectionTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mcpConnectionRequesterMock := interfaces_mocks.NewMockMCPConnectionRequesterInterface(mockController)

	sut := use_cases.NewRemoveAgentMCPConnection(
		agentRepositoryMock,
		mcpConnectionRequesterMock,
	)

	return &WhenRemovingAgentMCPConnectionTestingSuite{
		sut: sut,

		agentRepositoryMock:        agentRepositoryMock,
		mcpConnectionRequesterMock: mcpConnectionRequesterMock,
	}
}

func generateValidRemoveAgentMCPConnectionRequest() *requests.RemoveAgentMCPConnectionRequest {
	wallet := common.HexToAddress(common_tests.RandomEthAddress().Hex())
	return &requests.RemoveAgentMCPConnectionRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
			AgentCollectionTokenId: 1,
			AuthenticatedWallet:    &wallet,
		},
		MCPConnectionId: uuid.New(),
	}
}

func TestWhenRemovingAgentMCPConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because authenticated wallet is missing", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()
			req.AuthenticatedWallet = nil

			err := suite.sut.Remove(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because mcp connection id is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()
			req.MCPConnectionId = uuid.Nil

			err := suite.sut.Remove(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error occurs while getting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentMCPConnectionTestingSuite, request *requests.RemoveAgentMCPConnectionRequest) {
			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Remove(req)
			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent")
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentMCPConnectionTestingSuite, request *requests.RemoveAgentMCPConnectionRequest) {
			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()
			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Remove(req)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given caller is not the owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentMCPConnectionTestingSuite, request *requests.RemoveAgentMCPConnectionRequest) {
			otherOwner := common_tests.RandomEthAddress().Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(otherOwner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				Build()

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()
			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Remove(req)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given agent is hunan", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentMCPConnectionTestingSuite, request *requests.RemoveAgentMCPConnectionRequest) {
			agent := model_builders.NewAgentBuilder().
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithRole(enums.AgentRoleHuman).
				Build()

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()
			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Remove(req)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error occurs while removing mcp connection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentMCPConnectionTestingSuite, request *requests.RemoveAgentMCPConnectionRequest) {
			owner := request.AuthenticatedWallet.Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				Build()

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mcpConnectionRequesterMock.EXPECT().
				Remove(request.MCPConnectionId).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Remove(req)
			assert.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given all steps are successful", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentMCPConnectionTestingSuite, request *requests.RemoveAgentMCPConnectionRequest) {
			owner := request.AuthenticatedWallet.Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				Build()

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mcpConnectionRequesterMock.EXPECT().
				Remove(request.MCPConnectionId).
				Return(nil)
		}

		t.Run("Should not return error and call requester with expected payload", func(t *testing.T) {
			t.Parallel()
			suite := WhenRemovingAgentMCPConnectionBeforeEach(t)
			req := generateValidRemoveAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Remove(req)
			assert.NoError(t, err)
		})
	})
}
