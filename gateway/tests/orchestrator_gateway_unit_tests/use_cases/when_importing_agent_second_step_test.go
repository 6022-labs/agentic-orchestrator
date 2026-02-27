package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	gateway_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenImportingAgentSecondStepTestSuite struct {
	sut *use_cases.ImportExternalAgentSecondStep

	agentRepository                 *repositories_mocks.MockAgentRepositoryInterface
	llmConnectionRequester          *gateway_interfaces_mocks.MockLLMConnectionRequesterInterface
	externalAgentInformationFetcher *interfaces_mocks.MockExternalAgentInformationFetcherInterface
	databaseConnectionRequester     *gateway_interfaces_mocks.MockDatabaseConnectionRequesterInterface
}

func WhenImportingAgentSecondStepBeforeEach(t *testing.T) *WhenImportingAgentSecondStepTestSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	llmConnectionRequester := gateway_interfaces_mocks.NewMockLLMConnectionRequesterInterface(mockController)
	databaseConnectionRequester := gateway_interfaces_mocks.NewMockDatabaseConnectionRequesterInterface(mockController)
	externalAgentInformationFetcher := interfaces_mocks.NewMockExternalAgentInformationFetcherInterface(mockController)

	sut := use_cases.NewImportExternalAgentSecondStep(
		agentRepository,
		llmConnectionRequester,
		databaseConnectionRequester,
		externalAgentInformationFetcher,
	)

	return &WhenImportingAgentSecondStepTestSuite{
		sut: sut,

		agentRepository:                 agentRepository,
		llmConnectionRequester:          llmConnectionRequester,
		externalAgentInformationFetcher: externalAgentInformationFetcher,
		databaseConnectionRequester:     databaseConnectionRequester,
	}
}

func generateValidImportExternalAgentSecondStepRequest() *requests.ImportExternalAgentSecondStepRequest {
	llmApiKey := "AZERTYUIOP"

	authenticatedWallet := common_tests.RandomEthAddress()
	collectionAddress := common_tests.RandomEthAddress().Hex()

	return &requests.ImportExternalAgentSecondStepRequest{
		AgentByCollectionTokenIdRequest: requests.AgentByCollectionTokenIdRequest{
			AgentCollectionTokenId: 1,
			AgentCollectionAddress: collectionAddress,
			AuthenticatedWallet:    &authenticatedWallet,
		},
		DatabasePort:   5432,
		LLMModel:       "gpt-3",
		LLMApiKey:      &llmApiKey,
		DatabaseHost:   "localhost",
		LLMPlatform:    enums.LLMConnectionPlatformOpenAI,
		DatabaseDriver: enums.DatabaseConnectionDriverPostgres,
	}
}

func TestWhenImportingAgentSecondStep(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because LLM model is too short", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()
			request.LLMModel = "A"

			suite := WhenImportingAgentSecondStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because LLM platform is invalid", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()
			request.LLMPlatform = "invalid"

			suite := WhenImportingAgentSecondStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because DatabaseDriver is invalid", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()
			request.DatabaseDriver = "invalid"

			suite := WhenImportingAgentSecondStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because LLM platform is local and base url is missing", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()
			request.LLMPlatform = enums.GetAllLocalLLMConnectionPlatforms()[0]
			request.LLMBaseUrl = nil

			suite := WhenImportingAgentSecondStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because LLM platform is cloud and API key is missing", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()
			request.LLMPlatform = enums.GetAllCloudLLMConnectionPlatforms()[0]
			request.LLMApiKey = nil

			suite := WhenImportingAgentSecondStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given agent is not found in persistence", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentSecondStepTestSuite) {
			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenImportingAgentSecondStepBeforeEach(t)
			initSuite(suite)

			request := generateValidImportExternalAgentSecondStepRequest()
			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentSecondStepTestSuite, request *requests.ImportExternalAgentSecondStepRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(true).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()

			suite := WhenImportingAgentSecondStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent importation status is not first step", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentSecondStepTestSuite, request *requests.ImportExternalAgentSecondStepRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithOwner(owner).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusSecondStep.ToPointer()).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()

			suite := WhenImportingAgentSecondStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentSecondStepTestSuite, request *requests.ImportExternalAgentSecondStepRequest) {
			otherOwner := common_tests.RandomEthAddress().Hex()
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithOwner(otherOwner).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusFirstStep.ToPointer()).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()

			suite := WhenImportingAgentSecondStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "caller is not the owner")
		})
	})

	t.Run("Given agent is external and is first step importation status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentSecondStepTestSuite, request *requests.ImportExternalAgentSecondStepRequest) {
			owner := request.AuthenticatedWallet.Hex()
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithOwner(owner).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusFirstStep.ToPointer()).
				Build()

			suite.agentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).
				Return(agent, nil)

			suite.agentRepository.EXPECT().Save(gomock.Any()).Return(nil, nil)

			suite.llmConnectionRequester.EXPECT().Save(gomock.Any()).Return(nil)
			suite.databaseConnectionRequester.EXPECT().Create(gomock.Any()).Return(&interfaces.CreateDatabaseConnectionResult{
				DatabaseName:     "azertiop",
				DatabaseUsername: "azertiop",
				DatabasePassword: "azertiop",
			}, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()

			suite := WhenImportingAgentSecondStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.NoError(t, err)
		})

		t.Run("Should save agent importation status", func(t *testing.T) {
			t.Parallel()

			request := generateValidImportExternalAgentSecondStepRequest()

			suite := WhenImportingAgentSecondStepBeforeEach(t)
			initSuite(suite, request)

			suite.agentRepository.EXPECT().Save(gomock.Any()).DoAndReturn(func(agent *models.Agent) (*models.Agent, error) {
				assert.Equal(t, enums.AgentImportationStatusSecondStep, *agent.ImportationStatus)

				return agent, nil
			})

			suite.sut.Execute(request)
		})
	})
}
