package use_cases_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenImportingAgentFirstStepTestSuite struct {
	sut *use_cases.ImportExternalAgentFirstStep

	agentRepository                 *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRequester        *interfaces_mocks.MockAgentCollectionRequesterInterface
	orchestratorDomainResolver      *interfaces_mocks.MockOrchestratorDomainResolverInterface
	externalAgentInformationFetcher *interfaces_mocks.MockExternalAgentInformationFetcherInterface
}

func WhenImportingAgentFirstStepBeforeEach(t *testing.T) *WhenImportingAgentFirstStepTestSuite {
	mockController := gomock.NewController(t)

	agentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRequester := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)
	orchestratorDomainResolver := interfaces_mocks.NewMockOrchestratorDomainResolverInterface(mockController)
	externalAgentInformationFetcher := interfaces_mocks.NewMockExternalAgentInformationFetcherInterface(mockController)

	sut := use_cases.NewImportExternalAgentFirstStep(
		zap.NewNop(),
		agentRepository,
		agentCollectionRequester,
		orchestratorDomainResolver,
		externalAgentInformationFetcher,
	)

	return &WhenImportingAgentFirstStepTestSuite{
		sut: sut,

		agentRepository:                 agentRepository,
		agentCollectionRequester:        agentCollectionRequester,
		orchestratorDomainResolver:      orchestratorDomainResolver,
		externalAgentInformationFetcher: externalAgentInformationFetcher,
	}
}

func generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest() *requests.AgentByCollectionTokenIdRequest {
	authenticatedWallet := common_tests.RandomEthAddress()
	collectionAddress := common_tests.RandomEthAddress()

	return &requests.AgentByCollectionTokenIdRequest{
		AgentCollectionTokenId: 1,
		AgentCollectionAddress: collectionAddress.Hex(),
		AuthenticatedWallet:    &authenticatedWallet,
	}
}

func TestWhenImportingAgentFirstStep(t *testing.T) {
	t.Parallel()

	t.Run("Given agent is not found in persistence", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentFirstStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentFirstStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentFirstStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(true).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentFirstStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent importation status is higher than first step", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentFirstStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusSecondStep.ToPointer()).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return err", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentFirstStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given authenticated wallet is not the owner of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentFirstStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(*agent.CollectionTokenId)
			otherOwner := common_tests.RandomEthAddress()
			suite.agentCollectionRequester.
				EXPECT().
				GetAgentOwner(collectionAddress, *collectionTokenId).
				Return(&otherOwner, nil)
		}

		t.Run("Should return err", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentFirstStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no previous signer is found", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentFirstStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(*agent.CollectionTokenId)
			suite.agentCollectionRequester.EXPECT().GetAgentOwner(collectionAddress, *collectionTokenId).Return(request.AuthenticatedWallet, nil)

			previousSigners := []common.Address{}
			suite.agentCollectionRequester.EXPECT().GetAgentPreviousSigners(collectionAddress, *collectionTokenId).Return(previousSigners, nil)
		}

		t.Run("Should return err", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentFirstStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no previous signer return metadata", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentFirstStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(*agent.CollectionTokenId)
			suite.agentCollectionRequester.EXPECT().GetAgentOwner(collectionAddress, *collectionTokenId).Return(request.AuthenticatedWallet, nil)

			previousSigners := []common.Address{
				common.HexToAddress("0x88D927EBd1336aAF130e84f99A1E8EeeE0c0148b"),
				common.HexToAddress("0xCB70aF0EE6b6f1C86DCAA6fea7CeB4AF34363D6F"),
			}
			previousSignersDomains := []string{"signer1.com", "signer2.com"}

			suite.agentCollectionRequester.EXPECT().GetAgentPreviousSigners(collectionAddress, *collectionTokenId).Return(previousSigners, nil)

			suite.orchestratorDomainResolver.EXPECT().ResolveAIDomain(previousSigners[0]).Return(&previousSignersDomains[0], nil)
			suite.orchestratorDomainResolver.EXPECT().ResolveAIDomain(previousSigners[1]).Return(&previousSignersDomains[1], nil)

			suite.externalAgentInformationFetcher.
				EXPECT().
				IsInternalAgent(previousSignersDomains[0], *agent.CollectionTokenId).
				Return(nil, assert.AnError)

			suite.externalAgentInformationFetcher.
				EXPECT().
				IsInternalAgent(previousSignersDomains[1], *agent.CollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return err", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentFirstStepBeforeEach(t)
			initSuite(suite, request)

			_, err := suite.sut.Execute(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given previous signer return metadata", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentFirstStepTestSuite, request *requests.AgentByCollectionTokenIdRequest) *interfaces.AgentMetadataInformation {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithOwner(request.AuthenticatedWallet.Hex()).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(*agent.CollectionTokenId)
			suite.agentCollectionRequester.EXPECT().GetAgentOwner(collectionAddress, *collectionTokenId).Return(request.AuthenticatedWallet, nil)

			previousSigners := []common.Address{
				common.HexToAddress("0x88D927EBd1336aAF130e84f99A1E8EeeE0c0148b"),
				common.HexToAddress("0xCB70aF0EE6b6f1C86DCAA6fea7CeB4AF34363D6F"),
			}
			previousSignersDomains := []string{"signer1.com", "signer2.com"}

			suite.agentCollectionRequester.EXPECT().GetAgentPreviousSigners(collectionAddress, *collectionTokenId).Return(previousSigners, nil)

			suite.orchestratorDomainResolver.EXPECT().ResolveAIDomain(previousSigners[0]).Return(&previousSignersDomains[0], nil)
			suite.orchestratorDomainResolver.EXPECT().ResolveAIDomain(previousSigners[1]).Return(&previousSignersDomains[1], nil)

			isKnown := true
			suite.externalAgentInformationFetcher.
				EXPECT().
				IsInternalAgent(previousSignersDomains[0], *agent.CollectionTokenId).
				Return(&isKnown, nil)

			suite.externalAgentInformationFetcher.
				EXPECT().
				IsInternalAgent(previousSignersDomains[1], *agent.CollectionTokenId).
				Return(nil, assert.AnError)

			metadata := &interfaces.AgentMetadataInformation{
				LLMModel:    "gpt-3",
				LLMPlatform: enums.LLMConnectionPlatformOpenAI,
			}

			suite.agentRepository.EXPECT().Save(gomock.Any()).Return(nil, nil)

			suite.externalAgentInformationFetcher.
				EXPECT().
				FetchAgentMetadataInformation(previousSignersDomains[0], *agent.CollectionTokenId).
				Return(metadata, nil)

			return metadata
		}

		t.Run("Should return metadata", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByCollectionTokenIdWithAuthenticatedWalletRequest()

			suite := WhenImportingAgentFirstStepBeforeEach(t)
			metadata := initSuite(suite, request)

			response, err := suite.sut.Execute(request)

			assert.NoError(t, err)
			assert.NotNil(t, response)

			assert.Equal(t, metadata.LLMModel, response.LLMModel)
			assert.Equal(t, metadata.LLMPlatform, response.LLMPlatform)
		})
	})
}
