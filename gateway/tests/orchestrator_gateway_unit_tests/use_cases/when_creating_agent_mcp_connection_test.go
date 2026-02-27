package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingAgentMCPConnectionTestingSuite struct {
	sut *use_cases.CreateAgentMCPConnection

	agentRepositoryMock        *repositories_mocks.MockAgentRepositoryInterface
	mcpConnectionRequesterMock *interfaces_mocks.MockMCPConnectionRequesterInterface
}

func WhenCreatingAgentMCPConnectionBeforeEach(t *testing.T) *WhenCreatingAgentMCPConnectionTestingSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mcpConnectionRequesterMock := interfaces_mocks.NewMockMCPConnectionRequesterInterface(mockController)

	sut := use_cases.NewCreateAgentMCPConnection(
		agentRepositoryMock,
		mcpConnectionRequesterMock,
	)

	return &WhenCreatingAgentMCPConnectionTestingSuite{
		sut: sut,

		agentRepositoryMock:        agentRepositoryMock,
		mcpConnectionRequesterMock: mcpConnectionRequesterMock,
	}
}

func generateValidCreateAgentMCPConnectionRequest() *requests.CreateAgentMCPConnectionRequest {
	wallet := common.HexToAddress(common_tests.RandomEthAddress().Hex())
	return &requests.CreateAgentMCPConnectionRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
			AgentCollectionTokenId: 1,
			AuthenticatedWallet:    &wallet,
		},
		Name: "mcp-server-1",
		Url:  "https://mcp.example.com",
	}
}

func TestWhenCreatingAgentMCPConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because authenticated wallet is missing", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.AuthenticatedWallet = nil

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name is lower than 3 characters", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Name = "ab"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name does not start with a letter", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Name = "123Invalids"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name contains underscore", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Name = "test_name"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name contains invalid characters", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Name = "inv@lid"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Url = ""

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is not a valid url hostname", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Url = "https://invalid-host"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is not a valid url ipv4", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Url = "https://256.256.256.256"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is not a valid url ipv6", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			req.Url = "https://[::1]"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because only header auth key provided", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			key := "Authorization"
			req.HeaderAuth = map[string]string{
				key: "",
			}

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because only header auth value provided", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			val := "Bearer <<token>>"
			req.HeaderAuth = map[string]string{
				"": val,
			}

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error occurs while getting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentMCPConnectionTestingSuite, request *requests.CreateAgentMCPConnectionRequest) {
			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Create(req)
			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to get agent")
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentMCPConnectionTestingSuite, request *requests.CreateAgentMCPConnectionRequest) {
			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()
			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Create(req)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given caller is not the owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentMCPConnectionTestingSuite, request *requests.CreateAgentMCPConnectionRequest) {
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
			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Create(req)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentMCPConnectionTestingSuite, request *requests.CreateAgentMCPConnectionRequest) {
			owner := request.AuthenticatedWallet.Hex()
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				WithRole(enums.AgentRoleHuman).
				Build()

			suite.agentRepositoryMock.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()
			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Create(req)
			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error occurs while creating mcp connection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentMCPConnectionTestingSuite, request *requests.CreateAgentMCPConnectionRequest) {
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
				Create(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()

			initSuite(suite, req)

			err := suite.sut.Create(req)
			assert.Error(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given all steps are successful", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentMCPConnectionTestingSuite, request *requests.CreateAgentMCPConnectionRequest) *models.Agent {
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

			return agent
		}

		t.Run("Should not return error and call requester with expected payload", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentMCPConnectionBeforeEach(t)
			req := generateValidCreateAgentMCPConnectionRequest()
			headerKey := "X-Auth"
			headerVal := "Bearer <<token>>"
			req.HeaderAuth = map[string]string{
				headerKey: headerVal,
			}

			agent := initSuite(suite, req)

			var captured *interfaces.CreateMCPConnectionRequest
			suite.mcpConnectionRequesterMock.EXPECT().
				Create(gomock.Any()).
				DoAndReturn(func(r *interfaces.CreateMCPConnectionRequest) error {
					captured = r
					return nil
				})

			err := suite.sut.Create(req)
			assert.NoError(t, err)
			if assert.NotNil(t, captured) {
				assert.Equal(t, req.Url, captured.Url)
				assert.Equal(t, req.Name, captured.Name)
				assert.Equal(t, agent.ID, captured.AgentId)
				assert.Equal(t, req.HeaderAuth, captured.HeaderAuth)
			}
		})
	})
}
