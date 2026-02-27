package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
	sc_sync_services "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenImportingAgentThroughPreviousSignerTestSuite struct {
	sut *sc_sync_services.AgentPreviousSignerImporter

	agentRepository                 *repositories_mocks.MockAgentRepositoryInterface
	agentPrivateKeyStorage          *interfaces_mocks.MockAgentPrivateKeyStorageInterface
	externalAgentInformationFetcher *interfaces_mocks.MockExternalAgentInformationFetcherInterface
}

func WhenImportingAgentThroughPreviousSignerBeforeEach(t *testing.T) *WhenImportingAgentThroughPreviousSignerTestSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRepository := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentPrivateKeyStorage := interfaces_mocks.NewMockAgentPrivateKeyStorageInterface(mockController)
	externalAgentInformationFetcher := interfaces_mocks.NewMockExternalAgentInformationFetcherInterface(mockController)

	sut := sc_sync_services.NewPreviousAgentSignerImporter(
		zap.NewNop(),
		agentRepository,
		agentPrivateKeyStorage,
		externalAgentInformationFetcher,
	)

	return &WhenImportingAgentThroughPreviousSignerTestSuite{
		sut: sut,

		agentRepository:                 agentRepository,
		agentPrivateKeyStorage:          agentPrivateKeyStorage,
		externalAgentInformationFetcher: externalAgentInformationFetcher,
	}
}

func TestWhenImportingAgentThroughPreviousSigner(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while fetching agent is known by previous signer", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentThroughPreviousSignerTestSuite, domain string, collectionTokenId uint64) {
			suite.externalAgentInformationFetcher.EXPECT().IsInternalAgent(domain, collectionTokenId).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			domain := "test.eth"
			agent := model_builders.NewAgentBuilder().Build()

			suite := WhenImportingAgentThroughPreviousSignerBeforeEach(t)
			initSuite(suite, domain, *agent.CollectionTokenId)

			err := suite.sut.Import(domain, agent)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given agent is unknown by previous signer", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentThroughPreviousSignerTestSuite, domain string, collectionTokenId uint64) {
			isKnown := false
			suite.externalAgentInformationFetcher.EXPECT().IsInternalAgent(domain, collectionTokenId).Return(&isKnown, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			domain := "test.eth"
			agent := model_builders.NewAgentBuilder().Build()

			suite := WhenImportingAgentThroughPreviousSignerBeforeEach(t)
			initSuite(suite, domain, *agent.CollectionTokenId)

			err := suite.sut.Import(domain, agent)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given error occurs while import agent from external orchestrator", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentThroughPreviousSignerTestSuite, domain string, collectionTokenId uint64) {
			isKnown := true
			suite.externalAgentInformationFetcher.EXPECT().IsInternalAgent(domain, collectionTokenId).Return(&isKnown, nil)

			suite.externalAgentInformationFetcher.
				EXPECT().
				ImportAgent(domain, collectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			domain := "test.eth"
			agent := model_builders.NewAgentBuilder().Build()

			suite := WhenImportingAgentThroughPreviousSignerBeforeEach(t)
			initSuite(suite, domain, *agent.CollectionTokenId)

			err := suite.sut.Import(domain, agent)

			assert.NotNil(t, err)
		})
	})

	t.Run("Given information are fetched successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenImportingAgentThroughPreviousSignerTestSuite, domain string, collectionTokenId uint64) {
			isKnown := true
			suite.externalAgentInformationFetcher.EXPECT().IsInternalAgent(domain, collectionTokenId).Return(&isKnown, nil)

			importAgentInformation := &interfaces.ImportAgentInformation{PrivateKey: "private-key"}
			suite.externalAgentInformationFetcher.
				EXPECT().
				ImportAgent(domain, collectionTokenId).
				Return(importAgentInformation, nil)

			suite.agentPrivateKeyStorage.EXPECT().Store(gomock.Any(), "private-key").Return(nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			domain := "test.eth"
			agent := model_builders.NewAgentBuilder().Build()

			suite := WhenImportingAgentThroughPreviousSignerBeforeEach(t)
			initSuite(suite, domain, *agent.CollectionTokenId)

			err := suite.sut.Import(domain, agent)

			assert.Nil(t, err)
		})

		t.Run("Should import agent private key", func(t *testing.T) {
			t.Parallel()

			domain := "test.eth"
			agent := model_builders.NewAgentBuilder().Build()

			suite := WhenImportingAgentThroughPreviousSignerBeforeEach(t)
			initSuite(suite, domain, *agent.CollectionTokenId)

			suite.agentPrivateKeyStorage.EXPECT().Store(agent.ID, gomock.Any()).DoAndReturn(func(id uuid.UUID, privateKey string) error {
				assert.Equal(t, "private-key", privateKey)

				return nil
			})

			err := suite.sut.Import(domain, agent)

			assert.Nil(t, err)
		})
	})
}
