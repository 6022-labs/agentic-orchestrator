package use_cases_test

import (
	"testing"

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

type WhenGettingAgentLLMSettingsTestingSuite struct {
	sut *use_cases.GetAgentLLMSettings

	mockAgentRepository        *repositories_mocks.MockAgentRepositoryInterface
	mockLLMConnectionRequester *interfaces_mocks.MockLLMConnectionRequesterInterface
}

func WhenGettingAgentLLMSettingsBeforeEach(t *testing.T) *WhenGettingAgentLLMSettingsTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mockLLMConnectionRequester := interfaces_mocks.NewMockLLMConnectionRequesterInterface(mockController)

	sut := use_cases.NewGetAgentLLMSettings(
		mockAgentRepository,
		mockLLMConnectionRequester,
	)

	return &WhenGettingAgentLLMSettingsTestingSuite{
		sut: sut,

		mockAgentRepository:        mockAgentRepository,
		mockLLMConnectionRequester: mockLLMConnectionRequester,
	}
}

func TestWhenGettingAgentLLMSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentLLMSettingsTestingSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentLLMSettingsBeforeEach(t)

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentLLMSettingsTestingSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentLLMSettingsBeforeEach(t)

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentLLMSettingsTestingSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(false).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentLLMSettingsBeforeEach(t)

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentLLMSettingsTestingSuite, request *requests.AgentByCollectionTokenIdRequest) {
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

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentLLMSettingsBeforeEach(t)

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentLLMSettingsTestingSuite, request *requests.AgentByCollectionTokenIdRequest) {
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

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentLLMSettingsBeforeEach(t)

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting LLM settings", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentLLMSettingsTestingSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockLLMConnectionRequester.EXPECT().
				GetSettings(agent.ID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentLLMSettingsBeforeEach(t)

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingAgentLLMSettingsTestingSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockLLMConnectionRequester.EXPECT().
				GetSettings(agent.ID).
				Return(&interfaces.LLMConnectionSettings{}, nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentLLMSettingsBeforeEach(t)

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()
			initSuite(suite, request)

			settings, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.NotNil(t, settings)
		})
	})
}
