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

type WhenSavingAgentPromptTemplateTestingSuite struct {
	sut *use_cases.SaveAgentPromptTemplate

	mockAgentRepository              *repositories_mocks.MockAgentRepositoryInterface
	mockAgentPromptTemplateRequester *interfaces_mocks.MockAgentPromptTemplateRequesterInterface
}

func WhenSavingAgentPromptTemplateBeforeEach(t *testing.T) *WhenSavingAgentPromptTemplateTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mockAgentPromptTemplateRequester := interfaces_mocks.NewMockAgentPromptTemplateRequesterInterface(mockController)

	sut := use_cases.NewSaveAgentPromptTemplate(
		mockAgentRepository,
		mockAgentPromptTemplateRequester,
	)

	return &WhenSavingAgentPromptTemplateTestingSuite{
		sut: sut,

		mockAgentRepository:              mockAgentRepository,
		mockAgentPromptTemplateRequester: mockAgentPromptTemplateRequester,
	}
}

func generateValidSaveAgentPromptTemplateRequest() *requests.SaveAgentPromptTemplateRequest {
	authenticatedWallet := common.HexToAddress("0x1234567890123456789012345678901234567890")

	return &requests.SaveAgentPromptTemplateRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionTokenId: 1,
			AuthenticatedWallet:    &authenticatedWallet,
		},
		Key:     enums.AgentPromptTemplateKeyCloneDefault,
		Content: "example content",
	}
}

func TestWhenSavingAgentPromptTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because key is invalid", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSavingAgentPromptTemplateBeforeEach(t)
				request := generateValidSaveAgentPromptTemplateRequest()
				request.Key = "invalid_key"

				err := sut.sut.Save(request)

				assert.Error(t, err)
			})
		})

		t.Run("Because content is empty", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				sut := WhenSavingAgentPromptTemplateBeforeEach(t)
				request := generateValidSaveAgentPromptTemplateRequest()
				request.Content = ""

				err := sut.sut.Save(request)

				assert.Error(t, err)
			})
		})
	})

	t.Run("Given agent repository returns error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidSaveAgentPromptTemplateRequest()
			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			request := generateValidSaveAgentPromptTemplateRequest()
			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not found")
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(owner).
				WithIsInternal(false).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return not internal error", func(t *testing.T) {
			t.Parallel()

			request := generateValidSaveAgentPromptTemplateRequest()
			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not internal")
		})
	})

	t.Run("Given caller is not the owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			otherOwner := common_tests.RandomEthAddress().Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(otherOwner).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return not the owner error", func(t *testing.T) {
			t.Parallel()

			request := generateValidSaveAgentPromptTemplateRequest()
			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not the owner")
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				WithRole(enums.AgentRoleHuman).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return is human error", func(t *testing.T) {
			t.Parallel()

			request := generateValidSaveAgentPromptTemplateRequest()
			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "is human")
		})
	})

	t.Run("Given saving prompt template fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(owner).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockAgentPromptTemplateRequester.EXPECT().
				SaveAgentPromptTemplate(agent.ID, request.Key, request.Content).
				Return(assert.AnError)
		}

		t.Run("Should return save error", func(t *testing.T) {
			t.Parallel()

			request := generateValidSaveAgentPromptTemplateRequest()
			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to save agent prompt template")
		})
	})

	t.Run("Given all is successful", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingAgentPromptTemplateTestingSuite, request *requests.SaveAgentPromptTemplateRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(owner).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockAgentPromptTemplateRequester.EXPECT().
				SaveAgentPromptTemplate(agent.ID, request.Key, request.Content).
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidSaveAgentPromptTemplateRequest()
			suite := WhenSavingAgentPromptTemplateBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.NoError(t, err)
		})
	})
}
