package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenAddingMessageToAgentContextTestingSuite struct {
	sut *use_cases.AddMessageToAgentContext

	agentRepository              *repositories_mocks.MockAgentRepositoryInterface
	agentContextMessageRequester *interfaces_mocks.MockAgentContextMessageRequesterInterface
}

func WhenAddingMessageToAgentContextBeforeEach(t *testing.T) *WhenAddingMessageToAgentContextTestingSuite {
	mockController := gomock.NewController(t)

	agentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentContextMessageRequester := interfaces_mocks.NewMockAgentContextMessageRequesterInterface(mockController)

	sut := use_cases.NewAddMessageToAgentContext(
		agentRepository,
		agentContextMessageRequester,
	)

	return &WhenAddingMessageToAgentContextTestingSuite{
		sut: sut,

		agentRepository:              agentRepository,
		agentContextMessageRequester: agentContextMessageRequester,
	}
}

func generateValidAddMessageToAgentContextRequest() *requests.AddMessageToAgentContextRequest {
	authenticatedWallet := common.HexToAddress("0x1234567890123456789012345678901234567890")
	return &requests.AddMessageToAgentContextRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
			AgentCollectionTokenId: 1,
			AuthenticatedWallet:    &authenticatedWallet,
		},
		ContextKey: "test-context-key",
		Message:    "test-message",
		Role:       enums.AgentContextMessageRoleUser,
	}
}

func TestWhenAddingMessageToAgentContext(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because contextKey is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()
			request.ContextKey = ""

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})

		t.Run("Because message is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()
			request.Message = ""

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})

		t.Run("Because role is invalid", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()
			request.Role = enums.AgentContextMessageRole("invalid-role")

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while requesting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToAgentContextTestingSuite, request *requests.AddMessageToAgentContextRequest) {
			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToAgentContextTestingSuite, request *requests.AddMessageToAgentContextRequest) {
			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToAgentContextTestingSuite, request *requests.AddMessageToAgentContextRequest) {
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(false).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToAgentContextTestingSuite, request *requests.AddMessageToAgentContextRequest) {
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				WithOwner("0xDifferentOwner").
				WithIsInternal(true).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while adding message to agent context", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToAgentContextTestingSuite, request *requests.AddMessageToAgentContextRequest) {
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.agentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToAgentContextTestingSuite, request *requests.AddMessageToAgentContextRequest) {
			tokenId := request.AgentCollectionTokenId
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&tokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.agentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.agentContextMessageRequester.EXPECT().
				AddMessagesToAgentContext(gomock.Any()).
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToAgentContextBeforeEach(t)
			request := generateValidAddMessageToAgentContextRequest()

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.NoError(t, err)
		})
	})
}
