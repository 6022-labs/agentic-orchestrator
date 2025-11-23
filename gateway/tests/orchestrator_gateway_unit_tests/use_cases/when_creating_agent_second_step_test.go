package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	publicKey           = "0x76B86F1deaa7241F409b8cA98098b287888e750C"
	authenticatedWallet = "0xd7Df1101B5933522c6469e57110FD4901F4ac7Ca"
)

type WhenCreatingAgentSecondStepTestSuite struct {
	orchestratorSettings *settings.OrchestratorSettings
	mintSignerMock       *interfaces_mocks.MockMintSignerInterface
	agentRepositoryMock  *repositories_mocks.MockAgentRepositoryInterface

	sut *use_cases.CreateAgentSecondStep
}

func WhenCreatingAgentSecondStepBeforeEach(t *testing.T) *WhenCreatingAgentSecondStepTestSuite {
	// Create mocks
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common.HexToAddress(publicKey),
	}

	mintSignerMock := interfaces_mocks.NewMockMintSignerInterface(mockController)
	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)

	sut := use_cases.NewCreateAgentSecondStep(
		mintSignerMock,
		orchestratorSettings,
		agentRepositoryMock,
	)

	return &WhenCreatingAgentSecondStepTestSuite{
		sut: sut,

		mintSignerMock:       mintSignerMock,
		agentRepositoryMock:  agentRepositoryMock,
		orchestratorSettings: orchestratorSettings,
	}
}

func generateValidAgentByIdRequest() *requests.AgentByIdRequest {
	return &requests.AgentByIdRequest{
		AgentId:             uuid.New(),
		AuthenticatedWallet: common.HexToAddress("0xd7Df1101B5933522c6469e57110FD4901F4ac7Ca"),
	}
}

func TestWhenCreatingAgentSecondStep(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Run("Because agent id is empty", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByIdRequest()
			request.AgentId = uuid.Nil

			basics := WhenCreatingAgentSecondStepBeforeEach(t)

			_, err := basics.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error happens while fetching agent from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentSecondStepTestSuite) {
			suite.agentRepositoryMock.EXPECT().GetById(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			initSuite(basics)

			_, err := basics.sut.Execute(generateValidAgentByIdRequest())

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentSecondStepTestSuite) {
			suite.agentRepositoryMock.EXPECT().GetById(gomock.Any()).Return(nil, nil)
		}

		t.Run("Should return an not found error", func(t *testing.T) {
			t.Parallel()

			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			initSuite(basics)

			_, err := basics.sut.Execute(generateValidAgentByIdRequest())

			assert.NotNil(t, err)
			assert.IsType(t, &errors.NotFoundError{}, err)
		})
	})

	t.Run("Given agent is already minted", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentSecondStepTestSuite) {
			agentCollectionTokenId := uint64(1)
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(&agentCollectionTokenId).
				Build()

			suite.agentRepositoryMock.EXPECT().GetById(gomock.Any()).Return(agent, nil)
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			initSuite(basics)

			_, err := basics.sut.Execute(generateValidAgentByIdRequest())

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given authenticated wallet is not the creator of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentSecondStepTestSuite) {
			agent := model_builders.NewAgentBuilder().
				WithCreator(common_tests.RandomEthAddress().Hex()).
				Build()

			suite.agentRepositoryMock.EXPECT().GetById(gomock.Any()).Return(agent, nil)
		}

		t.Run("Should return a validation error", func(t *testing.T) {
			t.Parallel()

			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			initSuite(basics)

			_, err := basics.sut.Execute(generateValidAgentByIdRequest())

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error happens while signing mint request", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingAgentSecondStepTestSuite, request *requests.AgentByIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(nil).
				WithCreator(request.AuthenticatedWallet.Hex()).
				Build()
			suite.agentRepositoryMock.EXPECT().GetById(gomock.Any()).Return(agent, nil)
			suite.mintSignerMock.EXPECT().
				Sign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			request := generateValidAgentByIdRequest()
			initSuite(basics, request)

			_, err := basics.sut.Execute(request)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given all steps are successful", func(t *testing.T) {
		t.Parallel()

		expectedSignature := "signature"

		initSuite := func(suite *WhenCreatingAgentSecondStepTestSuite, request *requests.AgentByIdRequest) *models.Agent {
			agent := model_builders.NewAgentBuilder().
				WithCollectionTokenId(nil).
				WithCreator(request.AuthenticatedWallet.Hex()).
				Build()
			suite.agentRepositoryMock.EXPECT().GetById(gomock.Any()).Return(agent, nil)
			suite.mintSignerMock.EXPECT().
				Sign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&expectedSignature, nil)

			return agent
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			request := generateValidAgentByIdRequest()
			initSuite(basics, request)

			_, err := basics.sut.Execute(request)

			assert.Nil(t, err)
		})

		t.Run("Should return the expected response", func(t *testing.T) {
			t.Parallel()

			basics := WhenCreatingAgentSecondStepBeforeEach(t)
			request := generateValidAgentByIdRequest()
			initSuite(basics, request)

			response, _ := basics.sut.Execute(request)

			assert.NotNil(t, response)
			assert.Equal(t, expectedSignature, response.Signature)
			assert.Equal(t, publicKey, response.Signer)
			assert.NotEmpty(t, response.Nonce)
		})
	})
}
