package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite struct {
	sut *use_cases.SetupAgentDatabaseSettingsSecondStep

	agentRepositoryMock             *repositories_mocks.MockAgentRepositoryInterface
	databaseConnectionRequesterMock *interfaces_mocks.MockDatabaseConnectionRequesterInterface
}

func WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t *testing.T) *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	databaseConnectionRequesterMock := interfaces_mocks.NewMockDatabaseConnectionRequesterInterface(mockController)

	sut := use_cases.NewSetupAgentDatabaseSettingsSecondStep(
		agentRepositoryMock,
		databaseConnectionRequesterMock,
	)

	return &WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite{
		sut: sut,

		agentRepositoryMock:             agentRepositoryMock,
		databaseConnectionRequesterMock: databaseConnectionRequesterMock,
	}
}

func TestWhenSettingUpAgentDatabaseSettingsSecondStep(t *testing.T) {
	t.Parallel()

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(false).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given caller is not the owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(common_tests.RandomEthAddress().Hex()).
				WithIsInternal(true).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given agent is a clone and base agent fetch fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			// configure clone fields on agent returned by repo
			baseTokenId := uint64(999)
			agent.CloneOf = &baseTokenId
			agent.CollectionAddress = request.AgentCollectionAddress

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
			// fetching base agent returns an unexpected error
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(agent.CollectionAddress, baseTokenId).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given agent is a clone and base agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			baseTokenId := uint64(1000)
			agent.CloneOf = &baseTokenId
			agent.CollectionAddress = request.AgentCollectionAddress

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
			// base agent not found (nil, nil) to simulate ErrRecordNotFound being ignored and nil result
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(agent.CollectionAddress, baseTokenId).Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given agent is a clone and base agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			baseTokenId := uint64(1001)
			agent.CloneOf = &baseTokenId
			agent.CollectionAddress = request.AgentCollectionAddress

			agentToCopyFrom := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(agent.CollectionAddress, baseTokenId).Return(agentToCopyFrom, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given agent is a clone and base agent is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			baseTokenId := uint64(1002)
			agent.CloneOf = &baseTokenId
			agent.CollectionAddress = request.AgentCollectionAddress

			agentToCopyFrom := model_builders.NewAgentBuilder().
				WithIsInternal(true).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(agent.CollectionAddress, baseTokenId).Return(agentToCopyFrom, nil)
			suite.databaseConnectionRequesterMock.EXPECT().SetupAgentDatabase(agent.ID, &agentToCopyFrom.ID).Return(nil)
		}

		t.Run("Should not return error and pass copyFromAgentId", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given error running migrations", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
			suite.databaseConnectionRequesterMock.EXPECT().SetupAgentDatabase(agent.ID, nil).Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given all is successful", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsSecondStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
			suite.databaseConnectionRequesterMock.EXPECT().SetupAgentDatabase(agent.ID, nil).Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsSecondStepTestSuiteBeforeEach(t)
			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NoError(t, err)
		})
	})
}
