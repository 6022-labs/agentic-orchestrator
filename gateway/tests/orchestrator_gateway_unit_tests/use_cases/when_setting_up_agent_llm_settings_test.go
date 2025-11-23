package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSettingUpAgentLLMSettingsTestingSuite struct {
	sut *use_cases.SetupAgentLLMSettings

	mockAgentRepository        *repositories_mocks.MockAgentRepositoryInterface
	mockLLMConnectionRequester *interfaces_mocks.MockLLMConnectionRequesterInterface
}

func WhenSettingUpAgentLLMSettingsBeforeEach(t *testing.T) *WhenSettingUpAgentLLMSettingsTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	mockLLMConnectionRequester := interfaces_mocks.NewMockLLMConnectionRequesterInterface(mockController)

	sut := use_cases.NewSetupAgentLLMSettings(
		mockAgentRepository,
		mockLLMConnectionRequester,
	)

	return &WhenSettingUpAgentLLMSettingsTestingSuite{
		sut: sut,

		mockAgentRepository:        mockAgentRepository,
		mockLLMConnectionRequester: mockLLMConnectionRequester,
	}
}

func generateValidSetupAgentLLMSettingsRequest() *requests.SetupAgentLLMSettingsRequest {
	apiKey := "sk-1234567890abcdef1234567890abcdef1234567890"
	authenticatedWallet := common_tests.RandomEthAddress()
	collectionAddress := common_tests.RandomEthAddress().Hex()

	return &requests.SetupAgentLLMSettingsRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionTokenId: 1,
			AgentCollectionAddress: collectionAddress,
			AuthenticatedWallet:    &authenticatedWallet,
		},
		Model:    "gpt-3.5-turbo",
		ApiKey:   &apiKey,
		Platform: enums.LLMConnectionPlatformOpenAI,
	}
}

func TestWhenSettingUpAgentLLMSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because model is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentLLMSettingsBeforeEach(t)

			request := generateValidSetupAgentLLMSettingsRequest()
			request.Model = ""

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})

		t.Run("Because platform is invalid", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentLLMSettingsBeforeEach(t)

			request := generateValidSetupAgentLLMSettingsRequest()
			request.Platform = "invalid_platform"

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})

		t.Run("Because platform is cloud and api key is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentLLMSettingsBeforeEach(t)

			request := generateValidSetupAgentLLMSettingsRequest()
			request.Platform = enums.LLMConnectionPlatformOpenAI
			request.ApiKey = nil

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})

		t.Run("Because platform is self-hosted and host url is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentLLMSettingsBeforeEach(t)

			request := generateValidSetupAgentLLMSettingsRequest()
			request.Platform = enums.LLMConnectionPlatformOllama
			request.BaseUrl = nil

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})

		t.Run("Because platform is self-hosted and api key is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenSettingUpAgentLLMSettingsBeforeEach(t)

			request := generateValidSetupAgentLLMSettingsRequest()
			request.Platform = enums.LLMConnectionPlatformOllama
			request.ApiKey = nil

			err := sut.sut.Setup(request)

			assert.Error(t, err)
		})

		t.Run("Because platform is amazon bedrock and amazon platform region is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)

			apiKey := "api-key"

			request := generateValidSetupAgentLLMSettingsRequest()
			request.Platform = enums.LLMConnectionPlatformAmazonBedrock
			request.AmazonPlatformRegion = nil
			request.ApiKey = &apiKey
			request.BaseUrl = nil

			err := suite.sut.Setup(request)

			assert.Error(t, err)
			assert.IsType(t, err, &errors.ValidationError{})
			assert.Contains(t, err.Error(), "amazonPlatformRegion")
		})

		t.Run("Because amazon platform region is invalid", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)

			apiKey := "api-key"
			amazonPlatformRegion := "invalid-region"

			request := generateValidSetupAgentLLMSettingsRequest()
			request.Platform = enums.LLMConnectionPlatformAmazonBedrock
			request.AmazonPlatformRegion = &amazonPlatformRegion
			request.ApiKey = &apiKey
			request.BaseUrl = nil

			err := suite.sut.Setup(request)

			assert.Error(t, err)
			assert.IsType(t, err, &errors.ValidationError{})
			assert.Contains(t, err.Error(), "amazonPlatformRegion")
		})
	})

	t.Run("Given error occurs while requesting agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentLLMSettingsTestingSuite, request *requests.SetupAgentLLMSettingsRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)
			request := generateValidSetupAgentLLMSettingsRequest()
			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentLLMSettingsTestingSuite, request *requests.SetupAgentLLMSettingsRequest) {
			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)
			request := generateValidSetupAgentLLMSettingsRequest()
			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentLLMSettingsTestingSuite, request *requests.SetupAgentLLMSettingsRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithIsInternal(false).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)
			request := generateValidSetupAgentLLMSettingsRequest()
			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentLLMSettingsTestingSuite, request *requests.SetupAgentLLMSettingsRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(common_tests.RandomEthAddress().Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)
			request := generateValidSetupAgentLLMSettingsRequest()
			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is human", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentLLMSettingsTestingSuite, request *requests.SetupAgentLLMSettingsRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithIsInternal(true).
				WithRole(enums.AgentRoleHuman).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)
			request := generateValidSetupAgentLLMSettingsRequest()
			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentLLMSettingsTestingSuite, request *requests.SetupAgentLLMSettingsRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockLLMConnectionRequester.EXPECT().
				Save(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)
			request := generateValidSetupAgentLLMSettingsRequest()
			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSettingUpAgentLLMSettingsTestingSuite, request *requests.SetupAgentLLMSettingsRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithIsInternal(true).
				Build()

			suite.mockAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.mockLLMConnectionRequester.EXPECT().
				Save(gomock.Any()).
				Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSettingUpAgentLLMSettingsBeforeEach(t)
			request := generateValidSetupAgentLLMSettingsRequest()
			initSuite(suite, request)

			err := suite.sut.Setup(request)

			assert.NoError(t, err)
		})
	})
}
