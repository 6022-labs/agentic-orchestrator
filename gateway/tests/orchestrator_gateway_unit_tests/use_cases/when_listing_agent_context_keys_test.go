package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenListingAgentContextKeysTestingSuite struct {
	sut *use_cases.ListAgentContextKeys

	agentRepository       *repositories_mocks.MockAgentRepositoryInterface
	agentContextRequester *interfaces_mocks.MockAgentContextRequesterInterface
}

func WhenListingAgentContextKeysBeforeEach(t *testing.T) *WhenListingAgentContextKeysTestingSuite {
	mockController := gomock.NewController(t)

	agentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentContextRequester := interfaces_mocks.NewMockAgentContextRequesterInterface(mockController)

	sut := use_cases.NewListAgentContextKeys(
		agentRepository,
		agentContextRequester,
	)

	return &WhenListingAgentContextKeysTestingSuite{
		sut: sut,

		agentRepository:       agentRepository,
		agentContextRequester: agentContextRequester,
	}
}

func generateValidListAgentContextKeysRequest() *requests.ListAgentContextKeysRequest {
	request := &requests.ListAgentContextKeysRequest{
		AgentByCollectionTokenIdRequest: *generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest(),
	}

	request.Sanitize()

	return request
}

func TestWhenListingAgentContextKeys(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentContextKeysTestingSuite, request *requests.ListAgentContextKeysRequest) {
			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error and empty keys", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentContextKeysBeforeEach(t)
			request := generateValidListAgentContextKeysRequest()
			initSuite(suite, request)

			keys, err := suite.sut.List(request)

			assert.Error(t, err)
			assert.Empty(t, keys)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentContextKeysTestingSuite, request *requests.ListAgentContextKeysRequest) {
			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return error and empty keys", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentContextKeysBeforeEach(t)
			request := generateValidListAgentContextKeysRequest()
			initSuite(suite, request)

			keys, err := suite.sut.List(request)

			assert.Error(t, err)
			assert.Empty(t, keys)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentContextKeysTestingSuite, request *requests.ListAgentContextKeysRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(false).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error and empty keys", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentContextKeysBeforeEach(t)
			request := generateValidListAgentContextKeysRequest()
			initSuite(suite, request)

			keys, err := suite.sut.List(request)

			assert.Error(t, err)
			assert.Empty(t, keys)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentContextKeysTestingSuite, request *requests.ListAgentContextKeysRequest) {
			otherOwner := common_tests.RandomEthAddress().Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(otherOwner).
				WithIsInternal(true).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentContextKeysBeforeEach(t)
			request := generateValidListAgentContextKeysRequest()
			initSuite(suite, request)

			keys, err := suite.sut.List(request)

			assert.Error(t, err)
			assert.Empty(t, keys)
		})
	})

	t.Run("Given error occurs while getting agent contexts", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentContextKeysTestingSuite, request *requests.ListAgentContextKeysRequest) {
			agent := model_builders.NewAgentBuilder().
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithIsInternal(true).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.agentContextRequester.
				EXPECT().
				GetAgentContexts(agent.ID, request.Search, (*string)(request.SortBy), request.SortDirection.ToDatabaseDirection()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentContextKeysBeforeEach(t)
			request := generateValidListAgentContextKeysRequest()
			initSuite(suite, request)

			keys, err := suite.sut.List(request)

			assert.Error(t, err)
			assert.Empty(t, keys)
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingAgentContextKeysTestingSuite, request *requests.ListAgentContextKeysRequest) {
			agent := model_builders.NewAgentBuilder().
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithIsInternal(true).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.agentContextRequester.
				EXPECT().
				GetAgentContexts(agent.ID, request.Search, (*string)(request.SortBy), request.SortDirection.ToDatabaseDirection()).
				Return([]*models.AgentContext{
					{
						Id:  uuid.New(),
						Key: "key1",
					},
					{
						Id:  uuid.New(),
						Key: "key2",
					},
				}, nil)
		}

		t.Run("Should return keys", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingAgentContextKeysBeforeEach(t)
			request := generateValidListAgentContextKeysRequest()
			initSuite(suite, request)

			response, err := suite.sut.List(request)

			assert.NoError(t, err)
			assert.Equal(t, uint(1), response.Page)
			assert.Equal(t, uint(2), response.TotalElements)
			assert.Equal(t, uint(10), response.ElementsPerPage)
			assert.Len(t, response.ContextKeys, 2)
			assert.Equal(t, "key1", response.ContextKeys[0])
			assert.Equal(t, "key2", response.ContextKeys[1])
		})
	})
}
