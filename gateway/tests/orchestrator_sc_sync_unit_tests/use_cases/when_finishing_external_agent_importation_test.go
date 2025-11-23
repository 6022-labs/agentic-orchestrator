package use_cases_test

import (
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	sc_sync_services_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_sc_sync_mocks/services_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenFinishingExternalAgentImportationTestSuite struct {
	sut *use_cases.FinishAgentImportation

	agentRepository             *repositories_mocks.MockAgentRepositoryInterface
	orchestratorDomainResolver  *interfaces_mocks.MockOrchestratorDomainResolverInterface
	agentCollectionRequester    *interfaces_mocks.MockAgentCollectionRequesterInterface
	agentPreviousSignerImporter *sc_sync_services_mocks.MockAgentPreviousSignerImporterInterface
}

func WhenFinishingExternalAgentImportationBeforeEach(t *testing.T) *WhenFinishingExternalAgentImportationTestSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	orchestratorDomainResolver := interfaces_mocks.NewMockOrchestratorDomainResolverInterface(mockController)
	agentPreviousSignerImporter := sc_sync_services_mocks.NewMockAgentPreviousSignerImporterInterface(mockController)
	agentCollectionRequester := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)

	sut := use_cases.NewFinishAgentImportation(
		zap.NewNop(),
		agentRepository,
		agentCollectionRequester,
		orchestratorDomainResolver,
		agentPreviousSignerImporter,
	)

	return &WhenFinishingExternalAgentImportationTestSuite{
		sut: sut,

		agentRepository:             agentRepository,
		agentCollectionRequester:    agentCollectionRequester,
		orchestratorDomainResolver:  orchestratorDomainResolver,
		agentPreviousSignerImporter: agentPreviousSignerImporter,
	}
}

func TestWhenFinishingExternalAgentImportation(t *testing.T) {
	t.Parallel()

	t.Run("Given external agent is unknown", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingExternalAgentImportationTestSuite, request *requests.FinishAgentImportationRequest) {
			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}
			initSuite(suite, request)

			err := suite.sut.Finish(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is internal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingExternalAgentImportationTestSuite, request *requests.FinishAgentImportationRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(true).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Finish(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not third step importation status", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingExternalAgentImportationTestSuite, request *requests.FinishAgentImportationRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(true).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusSecondStep.ToPointer()).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Finish(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while querying previous signers of the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingExternalAgentImportationTestSuite, request *requests.FinishAgentImportationRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusThirdStep.ToPointer()).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(request.AgentCollectionTokenId)

			suite.agentCollectionRequester.
				EXPECT().
				GetAgentPreviousSigners(collectionAddress, *collectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Finish(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent has no previous signers", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingExternalAgentImportationTestSuite, request *requests.FinishAgentImportationRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusThirdStep.ToPointer()).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(request.AgentCollectionTokenId)

			previousSigners := []common.Address{}
			suite.agentCollectionRequester.
				EXPECT().
				GetAgentPreviousSigners(collectionAddress, *collectionTokenId).
				Return(previousSigners, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Finish(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given no previous signer can import the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingExternalAgentImportationTestSuite, request *requests.FinishAgentImportationRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusThirdStep.ToPointer()).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(request.AgentCollectionTokenId)

			previousSigners := []common.Address{
				common.HexToAddress("0xeAbC64AbEC6c599BA5B5CCE5Ee18B1C6DE525E33"),
				common.HexToAddress("0x9224778F646Fba91bD94396d798eA50981aAC092"),
			}
			previousSignersDomains := []string{"test.com", "test2.com"}

			suite.agentCollectionRequester.
				EXPECT().
				GetAgentPreviousSigners(collectionAddress, *collectionTokenId).
				Return(previousSigners, nil)

			suite.orchestratorDomainResolver.
				EXPECT().
				ResolveAIDomain(gomock.Any()).
				DoAndReturn(func(previousSigner common.Address) (*string, error) {
					if previousSigner == previousSigners[0] {
						return &(previousSignersDomains[0]), nil
					}
					if previousSigner == previousSigners[1] {
						return &(previousSignersDomains[1]), nil
					}

					return nil, assert.AnError
				}).
				AnyTimes()

			suite.agentPreviousSignerImporter.
				EXPECT().
				Import(gomock.Any(), agent).
				Return(assert.AnError).
				AnyTimes()
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Finish(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given previous signer can import the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenFinishingExternalAgentImportationTestSuite, request *requests.FinishAgentImportationRequest) {
			agent := model_builders.NewAgentBuilder().
				WithIsInternal(false).
				WithCollectionTokenId(&request.AgentCollectionTokenId).
				WithCollectionAddress(request.AgentCollectionAddress).
				WithImportationStatus(enums.AgentImportationStatusThirdStep.ToPointer()).
				Build()

			suite.agentRepository.EXPECT().GetByCollectionAddressAndCollectionTokenId(request.AgentCollectionAddress, request.AgentCollectionTokenId).Return(agent, nil)

			collectionAddress := common.HexToAddress(request.AgentCollectionAddress)
			collectionTokenId := new(big.Int).SetUint64(request.AgentCollectionTokenId)

			previousSigners := []common.Address{
				common.HexToAddress("0xeAbC64AbEC6c599BA5B5CCE5Ee18B1C6DE525E33"),
				common.HexToAddress("0x9224778F646Fba91bD94396d798eA50981aAC092"),
			}
			previousSignersDomains := []string{"test.com", "test2.com"}

			suite.agentRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(nil, nil)

			suite.agentCollectionRequester.
				EXPECT().
				GetAgentPreviousSigners(collectionAddress, *collectionTokenId).
				Return(previousSigners, nil)

			suite.orchestratorDomainResolver.
				EXPECT().
				ResolveAIDomain(gomock.Any()).
				DoAndReturn(func(previousSigner common.Address) (*string, error) {
					if previousSigner == previousSigners[0] {
						return &(previousSignersDomains[0]), nil
					}
					if previousSigner == previousSigners[1] {
						return &(previousSignersDomains[1]), nil
					}

					return nil, assert.AnError
				}).
				AnyTimes()

			suite.agentPreviousSignerImporter.
				EXPECT().
				Import(gomock.Any(), agent).
				DoAndReturn(func(previousSignersDomain string, agent *models.Agent) error {
					if previousSignersDomain == previousSignersDomains[0] {
						return nil
					}

					return assert.AnError
				}).
				AnyTimes()
		}

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			initSuite(suite, request)

			err := suite.sut.Finish(request)

			assert.Nil(t, err)
		})

		t.Run("Should save the agent with the correct status", func(t *testing.T) {
			t.Parallel()

			request := &requests.FinishAgentImportationRequest{
				AgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
				AgentCollectionTokenId: 1,
			}

			suite := WhenFinishingExternalAgentImportationBeforeEach(t)
			initSuite(suite, request)

			suite.agentRepository.
				EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(agent *models.Agent) (*models.Agent, error) {
					assert.Equal(t, enums.AgentImportationStatusFinished, *agent.ImportationStatus)

					return nil, nil
				})

			err := suite.sut.Finish(request)

			assert.Nil(t, err)
		})
	})
}
