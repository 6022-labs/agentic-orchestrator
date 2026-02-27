package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_sc_sync_mocks/services_mocks"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/orchestrator_sc_sync_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSynchronizingAllAgentCollectionsTestSuite struct {
	sut *use_cases.SynchronizeAllAgentCollections

	agentCollectionSynchronizerMock      *services_mocks.MockAgentCollectionSynchronizerInterface
	agentCollectionRepositoryMock        *repositories_mocks.MockAgentCollectionRepositoryInterface
	agentCollectionsManagerRequesterMock *interfaces_mocks.MockAgentCollectionsManagerRequesterInterface
}

func WhenSynchronizingAllAgentCollectionsBeforeEach(t *testing.T) *WhenSynchronizingAllAgentCollectionsTestSuite {
	mockController := gomock.NewController(t)

	agentCollectionSynchronizerMock := services_mocks.NewMockAgentCollectionSynchronizerInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	agentCollectionsManagerRequesterMock := interfaces_mocks.NewMockAgentCollectionsManagerRequesterInterface(mockController)

	sut := use_cases.NewSynchronizeAllAgentCollections(
		agentCollectionSynchronizerMock,
		agentCollectionRepositoryMock,
		agentCollectionsManagerRequesterMock,
	)

	return &WhenSynchronizingAllAgentCollectionsTestSuite{
		sut:                                  sut,
		agentCollectionSynchronizerMock:      agentCollectionSynchronizerMock,
		agentCollectionRepositoryMock:        agentCollectionRepositoryMock,
		agentCollectionsManagerRequesterMock: agentCollectionsManagerRequesterMock,
	}
}

func TestWhenSynchronizingAllAgentCollections(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while fetching collection addresses", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()
			suite := WhenSynchronizingAllAgentCollectionsBeforeEach(t)
			suite.agentCollectionsManagerRequesterMock.EXPECT().GetCollectionAddresses().Return(nil, assert.AnError)

			err := suite.sut.Synchronize()
			assert.Error(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given collection addresses are fetched", func(t *testing.T) {
		t.Parallel()

		addresses := []common.Address{
			common_tests.RandomEthAddress(),
			common_tests.RandomEthAddress(),
		}

		t.Run("And error occurs while checking existence in persistence", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentCollectionsBeforeEach(t)
				suite.agentCollectionsManagerRequesterMock.EXPECT().GetCollectionAddresses().Return(addresses, nil)

				suite.agentCollectionRepositoryMock.EXPECT().ExistsByCollectionAddress(addresses[0].Hex()).Return(false, assert.AnError)

				err := suite.sut.Synchronize()
				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while synchronizing missing collections", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentCollectionsBeforeEach(t)
				suite.agentCollectionsManagerRequesterMock.EXPECT().GetCollectionAddresses().Return(addresses, nil)

				suite.agentCollectionRepositoryMock.EXPECT().ExistsByCollectionAddress(addresses[0].Hex()).Return(false, nil)
				suite.agentCollectionSynchronizerMock.EXPECT().SyncByCollectionAddress(addresses[0].Hex()).Return(assert.AnError)

				err := suite.sut.Synchronize()
				assert.Error(t, err)
			})
		})

		t.Run("And no error occurs", func(t *testing.T) {
			t.Parallel()

			t.Run("Should check existence and synchronize missing collections", func(t *testing.T) {
				t.Parallel()
				suite := WhenSynchronizingAllAgentCollectionsBeforeEach(t)
				suite.agentCollectionsManagerRequesterMock.EXPECT().GetCollectionAddresses().Return(addresses, nil)

				// First address does not exist, should synchronize
				suite.agentCollectionRepositoryMock.EXPECT().ExistsByCollectionAddress(addresses[0].Hex()).Return(false, nil)
				suite.agentCollectionSynchronizerMock.EXPECT().SyncByCollectionAddress(addresses[0].Hex()).Return(nil)

				// Second address exists, should not synchronize
				suite.agentCollectionRepositoryMock.EXPECT().ExistsByCollectionAddress(addresses[1].Hex()).Return(true, nil)

				err := suite.sut.Synchronize()
				assert.NoError(t, err)
			})
		})
	})

	t.Run("Given no addresses to synchronize", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAllAgentCollectionsBeforeEach(t)
			suite.agentCollectionsManagerRequesterMock.EXPECT().GetCollectionAddresses().Return([]common.Address{}, nil)

			err := suite.sut.Synchronize()
			assert.NoError(t, err)
		})
	})
}
