package use_cases_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenImportingAgentThirdStepTestSuite struct {
	sut *use_cases.ImportExternalAgentThirdStep

	orchestratorSettings        *settings.OrchestratorSettings
	agentRepositoryMock         *repositories_mocks.MockAgentRepositoryInterface
	updateAgentSignerSignerMock *interfaces_mocks.MockUpdateAgentSignerSignerInterface
}

func WhenImportingAgentThirdStepBeforeEach(t *testing.T) *WhenImportingAgentThirdStepTestSuite {
	mockController := gomock.NewController(t)

	orchestratorSettings := &settings.OrchestratorSettings{
		PublicKey: common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
	}

	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	updateAgentSignerSignerMock := interfaces_mocks.NewMockUpdateAgentSignerSignerInterface(mockController)

	sut := use_cases.NewImportExternalAgentThirdStep(
		orchestratorSettings,
		agentRepositoryMock,
		updateAgentSignerSignerMock,
	)

	return &WhenImportingAgentThirdStepTestSuite{
		sut: sut,

		agentRepositoryMock:         agentRepositoryMock,
		orchestratorSettings:        orchestratorSettings,
		updateAgentSignerSignerMock: updateAgentSignerSignerMock,
	}
}

func TestWhenImportingAgentThirdStep(t *testing.T) {
	t.Parallel()

	t.Run("Given agent is not found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentThirdStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, nil)
		}

		t.Run("Should return not found error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentThirdStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Agent not found")
		})
	})

	t.Run("Given error occurs while fetching agent in persistence", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentThirdStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentThirdStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given caller is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentThirdStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithOwner("0xDifferentOwner").
				WithImportationStatus(enums.AgentImportationStatusSecondStep.ToPointer()).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentThirdStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "caller is not the owner")
		})
	})

	t.Run("Given all steps are successful", func(t *testing.T) {
		t.Parallel()

		expectedSignature := "signature"

		initSuite := func(suite *WhenImportingAgentThirdStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithImportationStatus(enums.AgentImportationStatusSecondStep.ToPointer()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				Build()

			suite.agentRepositoryMock.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionTokenId := big.NewInt(0).SetUint64(*agent.CollectionTokenId)
			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			suite.updateAgentSignerSignerMock.EXPECT().Sign(
				*collectionTokenId,
				suite.orchestratorSettings.PublicKey.Hex(),
				gomock.Any(), // nonce (uuid.UUID)
				collectionAddress,
			).Return(&expectedSignature, nil)

			suite.agentRepositoryMock.EXPECT().Save(agent).Return(nil, nil)
		}

		t.Run("Should return the expected response", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentThirdStepBeforeEach(t)
			initSuite(suite, request)

			response, err := suite.sut.Execute(request)

			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, expectedSignature, response.Signature)
		})
	})
}
