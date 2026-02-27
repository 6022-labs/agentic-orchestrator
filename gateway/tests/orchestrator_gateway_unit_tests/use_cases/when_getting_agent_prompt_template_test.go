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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.GetAgentPromptTemplate

	mockAgentRepository              *repositories_mocks.MockAgentRepositoryInterface
	mockAgentPromptTemplateRequester *interfaces_mocks.MockAgentPromptTemplateRequesterInterface
}

func WhenGettingAgentPromptTemplateBeforeEach(t *testing.T) *WhenGettingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mockAgentPromptTemplateRequester := interfaces_mocks.NewMockAgentPromptTemplateRequesterInterface(mockController)

	sut := use_cases.NewGetAgentPromptTemplate(
		mockAgentRepository,
		mockAgentPromptTemplateRequester,
	)

	return &WhenGettingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentRepository:              mockAgentRepository,
		mockAgentPromptTemplateRequester: mockAgentPromptTemplateRequester,
	}
}

func generateValidGetAgentPromptTemplateRequest() *requests.GetAgentPromptTemplateRequest {
	authenticatedWallet := common_tests.RandomEthAddress()
	collectionAddress := common_tests.RandomEthAddress().Hex()

	return &requests.GetAgentPromptTemplateRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionAddress: collectionAddress,
			AgentCollectionTokenId: 1,
			AuthenticatedWallet:    &authenticatedWallet,
		},
		Key: enums.AgentPromptTemplateKeyCloneDefault,
	}
}

func TestWhenGettingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because key is invalid", func(t *testing.T) {
			t.Parallel()

			sut := WhenGettingAgentPromptTemplateBeforeEach(t)
			request := generateValidGetAgentPromptTemplateRequest()
			request.Key = enums.AgentPromptTemplateKey("invalid_key")

			_, err := sut.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent repository returns error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(owner).
				WithIsInternal(false).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			otherOwner := common_tests.RandomEthAddress().Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(otherOwner).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(owner).
				WithIsInternal(false).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent prompt template requester returns error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(owner).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockAgentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agent.ID, request.Key).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent prompt template requester returns empty template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(owner).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockAgentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agent.ID, request.Key).
				Return(nil, nil)
		}

		t.Run("Should return empty template", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			template, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.Nil(t, template)
		})
	})

	t.Run("Given agent prompt template requester returns template", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentPromptTemplateTestingSuite, request *requests.GetAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(owner).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			template := "template"
			suite.mockAgentPromptTemplateRequester.EXPECT().
				GetAgentPromptTemplateByAgentIdAndKey(agent.ID, request.Key).
				Return(&template, nil)
		}

		t.Run("Should return template", func(t *testing.T) {
			t.Parallel()

			request := generateValidGetAgentPromptTemplateRequest()
			suite := WhenGettingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			template, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.Equal(t, "template", *template)
		})
	})
}
