package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
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

type WhenCreatingAgentContextTestingSuite struct {
	sut *use_cases.CreateAgentContext

	agentRepository       *repositories_mocks.MockAgentRepositoryInterface
	agentContextRequester *interfaces_mocks.MockAgentContextRequesterInterface
}

func WhenCreatingAgentContextBeforeEach(t *testing.T) *WhenCreatingAgentContextTestingSuite {
	mockController := gomock.NewController(t)

	agentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentContextRequester := interfaces_mocks.NewMockAgentContextRequesterInterface(mockController)

	sut := use_cases.NewCreateAgentContext(
		agentRepository,
		agentContextRequester,
	)

	return &WhenCreatingAgentContextTestingSuite{
		sut: sut,

		agentRepository:       agentRepository,
		agentContextRequester: agentContextRequester,
	}
}

func generateValidCreateAgentContextRequest() *requests.CreateAgentContextRequest {
	authenticatedWallet := common.HexToAddress("0x1234567890123456789012345678901234567890")
	return &requests.CreateAgentContextRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
			AgentCollectionTokenId: 1,
			AuthenticatedWallet:    &authenticatedWallet,
		},
		Key: "test-key",
	}
}

func TestWhenCreatingAgentContext(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because key is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenCreatingAgentContextBeforeEach(t)

			request := generateValidCreateAgentContextRequest()
			request.Key = ""

			err := sut.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite, request *requests.CreateAgentContextRequest) {
			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			request := generateValidCreateAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite, request *requests.CreateAgentContextRequest) {
			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			request := generateValidCreateAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite, request *requests.CreateAgentContextRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(owner).
				WithIsInternal(false).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenCreatingAgentContextBeforeEach(t)
			request := generateValidCreateAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite, request *requests.CreateAgentContextRequest) {
			owner := common_tests.RandomEthAddress().Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(owner).
				WithIsInternal(true).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			request := generateValidCreateAgentContextRequest()
			initSuite(suite, request)

			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not a expert", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite, request *requests.CreateAgentContextRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(owner).
				WithIsInternal(true).
				WithRole(enums.AgentRoleHuman).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			request := generateValidCreateAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while creating agent context", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite, request *requests.CreateAgentContextRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithIsInternal(true).
				Build()
			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.agentContextRequester.EXPECT().
				CreateAgentContext(agent.ID, request.Key).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			request := generateValidCreateAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentContextTestingSuite, request *requests.CreateAgentContextRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithOwner(owner).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithIsInternal(true).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.agentContextRequester.EXPECT().
				CreateAgentContext(agent.ID, request.Key).
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentContextBeforeEach(t)
			request := generateValidCreateAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Create(request)

			assert.NoError(t, err)
		})
	})
}
