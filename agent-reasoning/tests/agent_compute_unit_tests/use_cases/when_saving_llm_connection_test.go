package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/agent_compute_mocks/services_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSavingLLMConnectionTestingSuite struct {
	sut *use_cases.SaveLLMConnection

	mockAgentLLMSecretStorage   *services_mocks.MockAgentLLMSecretStorageInterface
	mockLLMConnectionRepository *repositories_mocks.MockLLMConnectionRepositoryInterface
}

func WhenSavingLLMConnectionBeforeEach(t *testing.T) *WhenSavingLLMConnectionTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentLLMSecretStorage := services_mocks.NewMockAgentLLMSecretStorageInterface(mockController)
	mockLLMConnectionRepository := repositories_mocks.NewMockLLMConnectionRepositoryInterface(mockController)

	sut := use_cases.NewSaveLLMConnection(mockAgentLLMSecretStorage, mockLLMConnectionRepository)

	return &WhenSavingLLMConnectionTestingSuite{
		sut: sut,

		mockAgentLLMSecretStorage:   mockAgentLLMSecretStorage,
		mockLLMConnectionRepository: mockLLMConnectionRepository,
	}
}

func generateValidSaveLLMConnectionRequest() *requests.SaveLLMConnectionRequest {
	baseUrl := "https://127.0.0.1"

	return &requests.SaveLLMConnectionRequest{
		BaseUrl:  &baseUrl,
		AgentId:  uuid.New(),
		Model:    "gpt-3.5-turbo",
		Platform: enums.LLMConnectionPlatformOllama,
	}
}

func TestWhenSavingLLMConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			request := generateValidSaveLLMConnectionRequest()
			request.AgentId = uuid.Nil

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "agentId")
		})

		t.Run("Because model is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			request := generateValidSaveLLMConnectionRequest()
			request.Model = ""

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "model")
		})

		t.Run("Because model is too short", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			request := generateValidSaveLLMConnectionRequest()
			request.Model = "a"

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "model")
		})

		t.Run("Because platform is invalid", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			request := generateValidSaveLLMConnectionRequest()
			request.Platform = enums.LLMConnectionPlatform("invalid")

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "platform")
		})

		t.Run("Because platform is local and baseUrl is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			request := generateValidSaveLLMConnectionRequest()
			request.Platform = enums.LLMConnectionPlatformOllama
			request.BaseUrl = nil

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "baseUrl")
		})

		t.Run("Because platform is local and baseUrl is malformed", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			request := generateValidSaveLLMConnectionRequest()
			malformedBaseUrl := "htp:/invalid-url"
			request.Platform = enums.LLMConnectionPlatformOllama
			request.BaseUrl = &malformedBaseUrl

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "baseUrl")
		})

		t.Run("Because platform is cloud and apiKey is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			request := generateValidSaveLLMConnectionRequest()
			request.Platform = enums.LLMConnectionPlatformOpenAI
			request.ApiKey = nil

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "apiKey")
		})

		t.Run("Because platform is amazon bedrock and amazon platform region is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			apiKey := "api-key"

			request := generateValidSaveLLMConnectionRequest()
			request.Platform = enums.LLMConnectionPlatformAmazonBedrock
			request.AmazonPlatformRegion = nil
			request.ApiKey = &apiKey
			request.BaseUrl = nil

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "amazonPlatformRegion")
		})

		t.Run("Because amazon platform region is invalid", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)

			apiKey := "api-key"
			amazonPlatformRegion := "invalid-region"

			request := generateValidSaveLLMConnectionRequest()
			request.Platform = enums.LLMConnectionPlatformAmazonBedrock
			request.AmazonPlatformRegion = &amazonPlatformRegion
			request.ApiKey = &apiKey
			request.BaseUrl = nil

			err := suite.sut.Save(request)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "amazonPlatformRegion")
		})
	})

	t.Run("Given error occurs while requesting existing connection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingLLMConnectionTestingSuite, request *requests.SaveLLMConnectionRequest) {
			suite.mockLLMConnectionRepository.EXPECT().
				GetByAgentId(request.AgentId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)
			request := generateValidSaveLLMConnectionRequest()
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving connection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingLLMConnectionTestingSuite, request *requests.SaveLLMConnectionRequest) {
			suite.mockLLMConnectionRepository.EXPECT().
				GetByAgentId(request.AgentId).
				Return(nil, nil)

			suite.mockLLMConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)
			request := generateValidSaveLLMConnectionRequest()
			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while storing API key", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingLLMConnectionTestingSuite, request *requests.SaveLLMConnectionRequest) {
			suite.mockLLMConnectionRepository.EXPECT().
				GetByAgentId(request.AgentId).
				Return(nil, nil)

			suite.mockLLMConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Return(model_builders.NewLLMConnectionBuilder().Build(), nil)

			suite.mockAgentLLMSecretStorage.EXPECT().
				Store(request.AgentId, *request.ApiKey).
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)
			request := generateValidSaveLLMConnectionRequest()
			apiKey := "api-key"
			request.ApiKey = &apiKey
			request.Platform = enums.LLMConnectionPlatformOpenAI

			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given connection already exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSavingLLMConnectionTestingSuite, request *requests.SaveLLMConnectionRequest) {
			agent := model_builders.NewLLMConnectionBuilder().
				WithModel(request.Model).
				WithAgentId(request.AgentId).
				WithBaseUrl(*request.BaseUrl).
				WithPlatform(request.Platform).
				Build()

			suite.mockLLMConnectionRepository.EXPECT().
				GetByAgentId(request.AgentId).
				Return(agent, nil)

			suite.mockLLMConnectionRepository.EXPECT().
				Save(gomock.Any()).
				Return(nil, nil)

			suite.mockAgentLLMSecretStorage.EXPECT().
				Store(request.AgentId, *request.ApiKey).
				Return(nil)
		}

		t.Run("Should store connection", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingLLMConnectionBeforeEach(t)
			request := generateValidSaveLLMConnectionRequest()
			apiKey := "api-key"
			request.ApiKey = &apiKey
			request.Platform = enums.LLMConnectionPlatformOpenAI

			initSuite(suite, request)

			err := suite.sut.Save(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given connection does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("And platform is amazon bedrock", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSavingLLMConnectionTestingSuite, request *requests.SaveLLMConnectionRequest) {
				suite.mockLLMConnectionRepository.EXPECT().
					GetByAgentId(request.AgentId).
					Return(nil, nil)

				suite.mockLLMConnectionRepository.EXPECT().
					Save(gomock.Any()).
					DoAndReturn(func(llmConnection *models.LLMConnection) (*models.LLMConnection, error) {
						assert.NotNil(t, llmConnection)
						assert.NotNil(t, llmConnection.BaseUrl)
						assert.Equal(t, request.Model, llmConnection.Model)
						assert.Equal(t, request.AgentId, llmConnection.AgentID)
						assert.Equal(t, request.Platform, llmConnection.Platform)

						return llmConnection, nil
					})

				suite.mockAgentLLMSecretStorage.EXPECT().
					Store(request.AgentId, *request.ApiKey).
					Return(nil)
			}

			t.Run("Should save connection with api key & host", func(t *testing.T) {
				t.Parallel()

				suite := WhenSavingLLMConnectionBeforeEach(t)
				request := generateValidSaveLLMConnectionRequest()
				apiKey := "api-key"
				amazonPlatformRegion := "us-west-2"
				request.BaseUrl = nil
				request.ApiKey = &apiKey
				request.AmazonPlatformRegion = &amazonPlatformRegion
				request.Platform = enums.LLMConnectionPlatformAmazonBedrock

				initSuite(suite, request)

				err := suite.sut.Save(request)

				assert.NoError(t, err)
			})
		})

		t.Run("And platform is cloud based", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSavingLLMConnectionTestingSuite, request *requests.SaveLLMConnectionRequest) {
				suite.mockLLMConnectionRepository.EXPECT().
					GetByAgentId(request.AgentId).
					Return(nil, nil)

				suite.mockLLMConnectionRepository.EXPECT().
					Save(gomock.Any()).
					DoAndReturn(func(llmConnection *models.LLMConnection) (*models.LLMConnection, error) {
						assert.NotNil(t, llmConnection)
						assert.Nil(t, llmConnection.BaseUrl)
						assert.Equal(t, request.Model, llmConnection.Model)
						assert.Equal(t, request.AgentId, llmConnection.AgentID)
						assert.Equal(t, request.Platform, llmConnection.Platform)

						return llmConnection, nil
					})

				suite.mockAgentLLMSecretStorage.EXPECT().
					Store(request.AgentId, *request.ApiKey).
					Return(nil)
			}

			t.Run("Should save connection with api key", func(t *testing.T) {
				t.Parallel()

				suite := WhenSavingLLMConnectionBeforeEach(t)
				request := generateValidSaveLLMConnectionRequest()
				apiKey := "api-key"
				request.BaseUrl = nil
				request.ApiKey = &apiKey
				request.Platform = enums.LLMConnectionPlatformOpenAI

				initSuite(suite, request)

				err := suite.sut.Save(request)

				assert.NoError(t, err)
			})
		})

		t.Run("And platform is local based", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenSavingLLMConnectionTestingSuite, request *requests.SaveLLMConnectionRequest) {
				suite.mockLLMConnectionRepository.EXPECT().
					GetByAgentId(request.AgentId).
					Return(nil, nil)

				suite.mockLLMConnectionRepository.EXPECT().
					Save(gomock.Any()).
					DoAndReturn(func(llmConnection *models.LLMConnection) (*models.LLMConnection, error) {
						assert.NotNil(t, llmConnection)
						assert.Equal(t, request.Model, llmConnection.Model)
						assert.Equal(t, request.AgentId, llmConnection.AgentID)
						assert.Equal(t, *request.BaseUrl, *llmConnection.BaseUrl)
						assert.Equal(t, request.Platform, llmConnection.Platform)

						return llmConnection, nil
					})

				suite.mockAgentLLMSecretStorage.EXPECT().
					Delete(request.AgentId).
					Return(nil)
			}

			t.Run("Should save connection with host & port", func(t *testing.T) {
				t.Parallel()

				suite := WhenSavingLLMConnectionBeforeEach(t)
				request := generateValidSaveLLMConnectionRequest()
				baseUrl := "https://custom.host"
				request.BaseUrl = &baseUrl
				request.Platform = enums.LLMConnectionPlatformOllama

				initSuite(suite, request)

				err := suite.sut.Save(request)

				assert.NoError(t, err)
			})
		})
	})
}
