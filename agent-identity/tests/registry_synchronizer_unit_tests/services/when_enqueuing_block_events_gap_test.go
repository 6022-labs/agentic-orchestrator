package services_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	synchronizer_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	sync_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/metrics_mocks/interfaces_mocks"
	synchronizer_repositories_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/repositories_mocks"
	synchronizer_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenEnqueuingBlockEventsGapTestingSuite struct {
	mockSyncMetrics *sync_metrics_mocks.MockSyncMetricsInterface

	sut *services.BlockEventsGapEnqueuer

	syncSettings *synchronizer_settings.SyncSettings

	mockBlockEventsQueue                           *synchronizer_interfaces_mocks.MockBlockEventsQueueInterface
	mockBlockchainRequester                        *synchronizer_interfaces_mocks.MockBlockchainRequesterInterface
	mockAgentCollectionBlockchainRequester         *synchronizer_interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface
	mockAgentCollectionsManagerBlockchainRequester *synchronizer_interfaces_mocks.MockAgentCollectionsManagerBlockchainRequesterInterface
	mockAgentCollectionRepository                  *repositories_mocks.MockAgentCollectionRepositoryInterface
	mockSyncCheckpointRepository                   *synchronizer_repositories_mocks.MockSyncCheckpointRepositoryInterface
}

func WhenEnqueuingBlockEventsGapBeforeEach(t *testing.T) *WhenEnqueuingBlockEventsGapTestingSuite {
	mockController := gomock.NewController(t)

	syncSettings := &synchronizer_settings.SyncSettings{StartBlock: 100, GapWindowBlocks: 10}

	mockBlockEventsQueue := synchronizer_interfaces_mocks.NewMockBlockEventsQueueInterface(mockController)
	mockBlockchainRequester := synchronizer_interfaces_mocks.NewMockBlockchainRequesterInterface(mockController)
	mockAgentCollectionBlockchainRequester := synchronizer_interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(
		mockController,
	)
	mockAgentCollectionsManagerBlockchainRequester := synchronizer_interfaces_mocks.NewMockAgentCollectionsManagerBlockchainRequesterInterface(
		mockController,
	)
	mockAgentCollectionRepository := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	mockSyncCheckpointRepository := synchronizer_repositories_mocks.NewMockSyncCheckpointRepositoryInterface(
		mockController,
	)

	mockSyncMetrics := sync_metrics_mocks.NewMockSyncMetricsInterface(mockController)

	sut := services.NewBlockEventsGapEnqueuer(
		zap.NewNop(),
		mockSyncMetrics,
		syncSettings,
		mockBlockEventsQueue,
		mockBlockchainRequester,
		mockAgentCollectionBlockchainRequester,
		mockAgentCollectionsManagerBlockchainRequester,
		mockAgentCollectionRepository,
		mockSyncCheckpointRepository,
	)

	return &WhenEnqueuingBlockEventsGapTestingSuite{
		sut:          sut,
		syncSettings: syncSettings,

		mockSyncMetrics: mockSyncMetrics,

		mockBlockEventsQueue:                           mockBlockEventsQueue,
		mockBlockchainRequester:                        mockBlockchainRequester,
		mockAgentCollectionBlockchainRequester:         mockAgentCollectionBlockchainRequester,
		mockAgentCollectionsManagerBlockchainRequester: mockAgentCollectionsManagerBlockchainRequester,
		mockAgentCollectionRepository:                  mockAgentCollectionRepository,
		mockSyncCheckpointRepository:                   mockSyncCheckpointRepository,
	}
}

func (suite *WhenEnqueuingBlockEventsGapTestingSuite) expectHeadAndCheckpoint(toBlock uint64) {
	suite.mockBlockchainRequester.EXPECT().
		GetCurrentBlockNumber(gomock.Any()).
		Return(&toBlock, nil).
		Times(1)

	suite.mockSyncCheckpointRepository.EXPECT().
		Get(gomock.Any()).
		Return(nil, nil).
		Times(1)
}

func (suite *WhenEnqueuingBlockEventsGapTestingSuite) expectCollectionQueries(times int) {
	suite.mockAgentCollectionBlockchainRequester.EXPECT().
		GetAgentCollectionV1MintedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	suite.expectNonMintedCollectionQueries(times)
}

func (suite *WhenEnqueuingBlockEventsGapTestingSuite) expectNonMintedCollectionQueries(times int) {
	requester := suite.mockAgentCollectionBlockchainRequester.EXPECT()
	requester.GetAgentCollectionV1TransferEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	requester.GetAgentCollectionV1AgentImageUpdatedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	requester.GetAgentCollectionV1AgentImageRemovedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	requester.GetAgentCollectionV1AgentAttributeUpdatedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	requester.GetAgentCollectionV1AgentAddressAddedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	requester.GetAgentCollectionV1AgentAddressRemovedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	requester.GetAgentCollectionV1MintProposalCreatedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
	requester.GetAgentCollectionV1MintProposalRemovedEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).Times(times)
}

func TestWhenEnqueuingBlockEventsGap(t *testing.T) {
	t.Parallel()

	t.Run("Given getting current block number fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			suite.mockBlockchainRequester.EXPECT().
				GetCurrentBlockNumber(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap(context.Background())

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "error getting current block number")
		})
	})

	t.Run("Given getting sync checkpoint fails with unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			toBlock := uint64(200)
			suite.mockBlockchainRequester.EXPECT().
				GetCurrentBlockNumber(gomock.Any()).
				Return(&toBlock, nil).
				Times(1)

			suite.mockSyncCheckpointRepository.EXPECT().
				Get(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap(context.Background())

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "failed to get sync checkpoint")
		})
	})

	t.Run("Given the sync checkpoint is past the current block", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			toBlock := uint64(149)
			suite.mockBlockchainRequester.EXPECT().
				GetCurrentBlockNumber(gomock.Any()).
				Return(&toBlock, nil).
				Times(1)

			suite.mockSyncCheckpointRepository.EXPECT().
				Get(gomock.Any()).
				Return(&models.SyncCheckpoint{Block: 150}, nil).
				Times(1)

			suite.mockSyncMetrics.EXPECT().RecordBackfilledBlocks(gomock.Any(), int64(0)).Times(1)
		}

		t.Run("Should succeed without querying any events", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap(context.Background())

			assert.Nil(t, err)
		})
	})

	t.Run("Given getting agent collections from the database fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			suite.expectHeadAndCheckpoint(105)

			suite.mockAgentCollectionRepository.EXPECT().
				GetAll(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap(context.Background())

			assert.ErrorIs(t, err, assert.AnError)
			assert.ErrorContains(t, err, "failed to get agent collections from database")
		})
	})

	t.Run("Given getting agent collections manager events fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			suite.expectHeadAndCheckpoint(105)

			suite.mockAgentCollectionRepository.EXPECT().
				GetAll(gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionsManagerBlockchainRequester.EXPECT().
				GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap(context.Background())

			assert.ErrorIs(t, err, assert.AnError)
			assert.ErrorContains(t, err, "failed to get agent collections manager collection created events")
		})
	})

	t.Run("Given no events in range", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			suite.expectHeadAndCheckpoint(100)

			suite.mockAgentCollectionRepository.EXPECT().
				GetAll(gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionsManagerBlockchainRequester.EXPECT().
				GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*abi.AgentCollectionsManagerCollectionCreated{}, nil).
				Times(1)

			suite.expectCollectionQueries(1)

			suite.mockSyncMetrics.EXPECT().RecordBackfilledBlocks(gomock.Any(), int64(0)).Times(1)
		}

		t.Run("Should succeed without enqueuing any events", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap(context.Background())

			assert.Nil(t, err)
		})
	})

	t.Run("Given a gap wider than one window", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			suite.expectHeadAndCheckpoint(125)

			suite.mockAgentCollectionRepository.EXPECT().
				GetAll(gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.expectCollectionQueries(3)
			suite.mockSyncMetrics.EXPECT().RecordBackfilledBlocks(gomock.Any(), int64(0)).Times(3)
		}

		t.Run("Should query each window of at most the configured size", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			manager := basics.mockAgentCollectionsManagerBlockchainRequester.EXPECT()
			gomock.InOrder(
				manager.GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), uint64(100), uint64(109)).
					Return(nil, nil).Times(1),
				manager.GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), uint64(110), uint64(119)).
					Return(nil, nil).Times(1),
				manager.GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), uint64(120), uint64(125)).
					Return(nil, nil).Times(1),
			)

			err := basics.sut.EnqueueGap(context.Background())

			assert.Nil(t, err)
		})
	})

	t.Run("Given a collection created in an earlier window", func(t *testing.T) {
		t.Parallel()

		collectionAddress := common_tests.RandomEthAddress()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			suite.expectHeadAndCheckpoint(119)

			suite.mockAgentCollectionRepository.EXPECT().
				GetAll(gomock.Any()).
				Return(nil, nil).
				Times(1)

			manager := suite.mockAgentCollectionsManagerBlockchainRequester.EXPECT()
			manager.GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), uint64(100), uint64(109)).
				Return([]*abi.AgentCollectionsManagerCollectionCreated{
					{CollectionAddress: collectionAddress, Raw: types.Log{BlockNumber: 105}},
				}, nil).
				Times(1)
			manager.GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), uint64(110), uint64(119)).
				Return(nil, nil).
				Times(1)

			suite.expectNonMintedCollectionQueries(2)

			queue := suite.mockBlockEventsQueue.EXPECT()
			queue.AddAgentCollectionsManagerCollectionCreatedEvents(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			queue.AddAgentCollectionV1MintedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1TransferEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1AgentImageUpdatedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1AgentImageRemovedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1AgentAttributeUpdatedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1AgentAddressAddedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1AgentAddressRemovedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1MintProposalCreatedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			queue.AddAgentCollectionV1MintProposalRemovedEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			suite.mockSyncMetrics.EXPECT().RecordBackfilledBlocks(gomock.Any(), int64(1)).Times(1)
			suite.mockSyncMetrics.EXPECT().RecordBackfilledBlocks(gomock.Any(), int64(0)).Times(1)
		}

		t.Run("Should query its events in later windows", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			minted := basics.mockAgentCollectionBlockchainRequester.EXPECT()
			minted.GetAgentCollectionV1MintedEvent(gomock.Any(), uint64(100), uint64(109), []common.Address{collectionAddress}).
				Return(nil, nil).
				Times(1)
			minted.GetAgentCollectionV1MintedEvent(gomock.Any(), uint64(110), uint64(119), []common.Address{collectionAddress}).
				Return(nil, nil).
				Times(1)

			err := basics.sut.EnqueueGap(context.Background())

			assert.Nil(t, err)
		})
	})
}
