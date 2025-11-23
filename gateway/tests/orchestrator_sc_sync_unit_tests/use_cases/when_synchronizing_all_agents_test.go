package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"

	"math/big"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_sc_sync_mocks/services_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSynchronizingAllAgentsTestSuite struct {
	sut *use_cases.SynchronizeAllAgents

	agentSynchronizerMock         *services_mocks.MockAgentSynchronizerInterface
	agentRepositoryMock           *repositories_mocks.MockAgentRepositoryInterface
	agentCollectionRequesterMock  *interfaces_mocks.MockAgentCollectionRequesterInterface
	agentCollectionRepositoryMock *repositories_mocks.MockAgentCollectionRepositoryInterface
}

func WhenSynchronizingAllAgentsBeforeEach(t *testing.T) *WhenSynchronizingAllAgentsTestSuite {
	mockController := gomock.NewController(t)

	agentSynchronizerMock := services_mocks.NewMockAgentSynchronizerInterface(mockController)
	agentRepositoryMock := repositories_mocks.NewMockAgentRepositoryInterface(mockController)
	agentCollectionRequesterMock := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)

	sut := use_cases.NewSynchronizeAllAgents(
		agentSynchronizerMock,
		agentRepositoryMock,
		agentCollectionRequesterMock,
		agentCollectionRepositoryMock,
	)

	return &WhenSynchronizingAllAgentsTestSuite{
		sut: sut,

		agentRepositoryMock:           agentRepositoryMock,
		agentSynchronizerMock:         agentSynchronizerMock,
		agentCollectionRequesterMock:  agentCollectionRequesterMock,
		agentCollectionRepositoryMock: agentCollectionRepositoryMock,
	}
}

func TestWhenSynchronizingAllAgents(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while fetching all agent collections", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenSynchronizingAllAgentsBeforeEach(t)
			suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(nil, assert.AnError)

			err := suite.sut.Synchronize()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to get all agent collections")
		})
	})

	t.Run("Given agent collections are fetched", func(t *testing.T) {
		t.Parallel()

		collections := []*models.AgentCollection{
			model_builders.NewAgentCollectionBuilder().WithAddress(common_tests.RandomEthAddress().Hex()).Build(),
			model_builders.NewAgentCollectionBuilder().WithAddress(common_tests.RandomEthAddress().Hex()).Build(),
		}

		t.Run("And error occurs while getting next token ID", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentsBeforeEach(t)
				suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(collections, nil)
				suite.agentCollectionRequesterMock.EXPECT().GetNextTokenId(common.HexToAddress(collections[0].Address)).Return(nil, assert.AnError)

				err := suite.sut.Synchronize()

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get next token ID")
			})
		})

		t.Run("And error occurs while checking agent existence", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentsBeforeEach(t)
				suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(collections, nil)
				nextTokenId := big.NewInt(2) // Setting to 2 so we have at least one token to check
				suite.agentCollectionRequesterMock.EXPECT().GetNextTokenId(common.HexToAddress(collections[0].Address)).Return(nextTokenId, nil)
				suite.agentRepositoryMock.EXPECT().ExistsByCollectionAddressAndCollectionTokenId(collections[0].Address, uint64(1)).Return(false, assert.AnError)

				err := suite.sut.Synchronize()

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to check existence of agent")
			})
		})

		t.Run("And error occurs while synchronizing agent", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentsBeforeEach(t)
				suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(collections, nil)
				nextTokenId := big.NewInt(2) // Setting to 2 so we have at least one token to sync
				suite.agentCollectionRequesterMock.EXPECT().GetNextTokenId(common.HexToAddress(collections[0].Address)).Return(nextTokenId, nil)
				suite.agentRepositoryMock.EXPECT().ExistsByCollectionAddressAndCollectionTokenId(collections[0].Address, uint64(1)).Return(false, nil)
				suite.agentSynchronizerMock.EXPECT().Sync(collections[0].Address, uint64(1)).Return(assert.AnError)

				err := suite.sut.Synchronize()

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to synchronize agent")
			})
		})

		t.Run("And no error occurs", func(t *testing.T) {
			t.Parallel()

			t.Run("Should skip existing agents and synchronize only non-existing agents for all collections", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentsBeforeEach(t)
				suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(collections, nil)
				nextTokenId := big.NewInt(3) // Setting to 3 so we have tokens 1 and 2 to process

				for _, collection := range collections {
					suite.agentCollectionRequesterMock.EXPECT().GetNextTokenId(common.HexToAddress(collection.Address)).Return(nextTokenId, nil)
					// Token 1 exists, should be skipped
					suite.agentRepositoryMock.EXPECT().ExistsByCollectionAddressAndCollectionTokenId(collection.Address, uint64(1)).Return(true, nil)
					// Token 2 doesn't exist, should be synchronized
					suite.agentRepositoryMock.EXPECT().ExistsByCollectionAddressAndCollectionTokenId(collection.Address, uint64(2)).Return(false, nil)
					suite.agentSynchronizerMock.EXPECT().Sync(collection.Address, uint64(2)).Return(nil)
				}

				err := suite.sut.Synchronize()

				assert.NoError(t, err)
			})

			t.Run("Should handle collections with no tokens to synchronize", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentsBeforeEach(t)
				suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(collections, nil)
				nextTokenId := big.NewInt(1) // nextTokenId of 1 means no tokens to process (loop from 1 to < 1)

				for _, collection := range collections {
					suite.agentCollectionRequesterMock.EXPECT().GetNextTokenId(common.HexToAddress(collection.Address)).Return(nextTokenId, nil)
					// No expectations for agentRepository or agentSynchronizer since no tokens should be processed
				}

				err := suite.sut.Synchronize()

				assert.NoError(t, err)
			})

			t.Run("Should synchronize all non-existing agents when none exist", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentsBeforeEach(t)
				suite.agentCollectionRepositoryMock.EXPECT().GetAll().Return(collections, nil)
				nextTokenId := big.NewInt(3) // tokens 1 and 2 to process

				for _, collection := range collections {
					suite.agentCollectionRequesterMock.EXPECT().GetNextTokenId(common.HexToAddress(collection.Address)).Return(nextTokenId, nil)
					// Neither token exists, both should be synchronized
					suite.agentRepositoryMock.EXPECT().ExistsByCollectionAddressAndCollectionTokenId(collection.Address, uint64(1)).Return(false, nil)
					suite.agentSynchronizerMock.EXPECT().Sync(collection.Address, uint64(1)).Return(nil)
					suite.agentRepositoryMock.EXPECT().ExistsByCollectionAddressAndCollectionTokenId(collection.Address, uint64(2)).Return(false, nil)
					suite.agentSynchronizerMock.EXPECT().Sync(collection.Address, uint64(2)).Return(nil)
				}

				err := suite.sut.Synchronize()

				assert.NoError(t, err)
			})
		})
	})
}
