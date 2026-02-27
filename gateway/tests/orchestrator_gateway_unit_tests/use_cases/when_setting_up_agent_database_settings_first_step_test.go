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
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite struct {
	sut *use_cases.SetupAgentDatabaseSettingsFirstStep

	agentRepositoryMock             *repositories_mocks.MockAgentRepositoryInterface
	databaseConnectionRequesterMock *interfaces_mocks.MockDatabaseConnectionRequesterInterface
}

func WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t *testing.T) *WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite {
	mockController := gomock.NewController(t)

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	databaseConnectionRequesterMock := interfaces_mocks.NewMockDatabaseConnectionRequesterInterface(mockController)

	sut := use_cases.NewSetupAgentDatabaseSettingsFirstStep(
		agentRepositoryMock,
		databaseConnectionRequesterMock,
	)

	return &WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite{
		sut: sut,

		agentRepositoryMock:             agentRepositoryMock,
		databaseConnectionRequesterMock: databaseConnectionRequesterMock,
	}
}

func generateValidSetupAgentDatabaseSettingsFirstStepRequest() *requests.SetupAgentDatabaseSettingsFirstStepRequest {
	authenticatedWallet := common.HexToAddress("0x1234567890123456789012345678901234567890")
	return &requests.SetupAgentDatabaseSettingsFirstStepRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionAddress: "0xabc",
			AgentCollectionTokenId: 1,
			AuthenticatedWallet:    &authenticatedWallet,
		},
		Host:   "localhost",
		Port:   5432,
		Driver: enums.DatabaseConnectionDriverPostgres,
	}
}

func TestWhenSettingUpAgentDatabaseSettingsFirstStep(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because driver is not valid", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t)

			request := generateValidSetupAgentDatabaseSettingsFirstStepRequest()
			request.Driver = "invalid-driver"

			_, err := suite.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent repository returns error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite, request *requests.SetupAgentDatabaseSettingsFirstStepRequest) {
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t)
			request := generateValidSetupAgentDatabaseSettingsFirstStepRequest()
			initSuite(suite, request)

			resp, err := suite.sut.Setup(request)

			assert.Nil(t, resp)
			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite, request *requests.SetupAgentDatabaseSettingsFirstStepRequest) {
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t)
			request := generateValidSetupAgentDatabaseSettingsFirstStepRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Setup(request)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not found")
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite, request *requests.SetupAgentDatabaseSettingsFirstStepRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(false).
				Build()
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t)
			request := generateValidSetupAgentDatabaseSettingsFirstStepRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Setup(request)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not internal")
		})
	})

	t.Run("Given caller is not the owner", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite, request *requests.SetupAgentDatabaseSettingsFirstStepRequest) {
			otherOwner := common_tests.RandomEthAddress().Hex()
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(otherOwner).
				WithIsInternal(true).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t)
			request := generateValidSetupAgentDatabaseSettingsFirstStepRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Setup(request)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not the owner")
		})
	})

	t.Run("Given database connection requester returns error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite, request *requests.SetupAgentDatabaseSettingsFirstStepRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			suite.databaseConnectionRequesterMock.EXPECT().Create(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t)
			request := generateValidSetupAgentDatabaseSettingsFirstStepRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Setup(request)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to create database connection")
		})
	})

	t.Run("Given all is successful", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentDatabaseSettingsFirstStepTestSuite, request *requests.SetupAgentDatabaseSettingsFirstStepRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithIsInternal(true).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			result := &interfaces.CreateDatabaseConnectionResult{
				DatabaseName:     "db",
				DatabaseUsername: "user",
				DatabasePassword: "pass",
			}
			suite.databaseConnectionRequesterMock.EXPECT().Create(gomock.Any()).Return(result, nil)
		}

		t.Run("Should return database connection details", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentDatabaseSettingsFirstStepTestSuiteBeforeEach(t)
			request := generateValidSetupAgentDatabaseSettingsFirstStepRequest()

			initSuite(suite, request)

			resp, err := suite.sut.Setup(request)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "db", resp.DatabaseName)
			assert.Equal(t, "user", resp.DatabaseUsername)
			assert.Equal(t, "pass", resp.DatabasePassword)
		})
	})
}
