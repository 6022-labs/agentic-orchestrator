package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	synchronizer_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	synchronizer_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	synchronizer_repositories_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenEnqueuingBlockEventsGapTestingSuite struct {
	sut *services.BlockEventsGapEnqueuer

	syncSettings *synchronizer_settings.SyncSettings

	mockBlockEventsQueue                          *synchronizer_interfaces_mocks.MockBlockEventsQueueInterface
	mockBlockchainRequester                       *synchronizer_interfaces_mocks.MockBlockchainRequesterInterface
	mockAgentCollectionBlockchainRequester        *synchronizer_interfaces_mocks.MockAgentCollectionBlockchainRequesterInterface
	mockAgentCollectionsManagerBlockchainRequester *synchronizer_interfaces_mocks.MockAgentCollectionsManagerBlockchainRequesterInterface
	mockAgentCollectionRepository                 *repositories_mocks.MockAgentCollectionRepositoryInterface
	mockSyncCheckpointRepository                  *synchronizer_repositories_mocks.MockSyncCheckpointRepositoryInterface
}

func WhenEnqueuingBlockEventsGapBeforeEach(t *testing.T) *WhenEnqueuingBlockEventsGapTestingSuite {
	mockController := gomock.NewController(t)

	syncSettings := &synchronizer_settings.SyncSettings{StartBlock: 100}

	mockBlockEventsQueue := synchronizer_interfaces_mocks.NewMockBlockEventsQueueInterface(mockController)
	mockBlockchainRequester := synchronizer_interfaces_mocks.NewMockBlockchainRequesterInterface(mockController)
	mockAgentCollectionBlockchainRequester := synchronizer_interfaces_mocks.NewMockAgentCollectionBlockchainRequesterInterface(mockController)
	mockAgentCollectionsManagerBlockchainRequester := synchronizer_interfaces_mocks.NewMockAgentCollectionsManagerBlockchainRequesterInterface(mockController)
	mockAgentCollectionRepository := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	mockSyncCheckpointRepository := synchronizer_repositories_mocks.NewMockSyncCheckpointRepositoryInterface(mockController)

	sut := services.NewBlockEventsGapEnqueuer(
		zap.NewNop(),
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

		mockBlockEventsQueue:                          mockBlockEventsQueue,
		mockBlockchainRequester:                       mockBlockchainRequester,
		mockAgentCollectionBlockchainRequester:        mockAgentCollectionBlockchainRequester,
		mockAgentCollectionsManagerBlockchainRequester: mockAgentCollectionsManagerBlockchainRequester,
		mockAgentCollectionRepository:                 mockAgentCollectionRepository,
		mockSyncCheckpointRepository:                  mockSyncCheckpointRepository,
	}
}

func TestWhenEnqueuingBlockEventsGap(t *testing.T) {
	t.Parallel()

	t.Run("Given getting current block number fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			suite.mockBlockchainRequester.EXPECT().
				GetCurrentBlockNumber().
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap()

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "error getting current block number")
		})
	})

	t.Run("Given getting sync checkpoint fails with unexpected error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			toBlock := uint64(200)
			suite.mockBlockchainRequester.EXPECT().
				GetCurrentBlockNumber().
				Return(&toBlock, nil).
				Times(1)

			suite.mockSyncCheckpointRepository.EXPECT().
				Get().
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap()

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "failed to get sync checkpoint")
		})
	})

	t.Run("Given getting agent collections manager events fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			toBlock := uint64(200)
			suite.mockBlockchainRequester.EXPECT().
				GetCurrentBlockNumber().
				Return(&toBlock, nil).
				Times(1)

			suite.mockSyncCheckpointRepository.EXPECT().
				Get().
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)

			suite.mockAgentCollectionsManagerBlockchainRequester.EXPECT().
				GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap()

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "failed to get agent collections manager collection created events")
		})
	})

	t.Run("Given no events in range", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnqueuingBlockEventsGapTestingSuite) {
			toBlock := uint64(100)
			suite.mockBlockchainRequester.EXPECT().
				GetCurrentBlockNumber().
				Return(&toBlock, nil).
				Times(1)

			suite.mockSyncCheckpointRepository.EXPECT().
				Get().
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)

			suite.mockAgentCollectionsManagerBlockchainRequester.EXPECT().
				GetAgentCollectionsManagerCollectionCreatedEvent(gomock.Any(), gomock.Any()).
				Return([]*abi.AgentCollectionsManagerCollectionCreated{}, nil).
				Times(1)

			suite.mockAgentCollectionRepository.EXPECT().
				GetAll().
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionBlockchainRequester.EXPECT().
				GetAgentCollectionV1MintedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionBlockchainRequester.EXPECT().
				GetAgentCollectionV1TransferEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionBlockchainRequester.EXPECT().
				GetAgentCollectionV1AgentImageUpdatedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionBlockchainRequester.EXPECT().
				GetAgentCollectionV1AgentSignerUpdatedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionBlockchainRequester.EXPECT().
				GetAgentCollectionV1AgentAttributeUpdatedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionBlockchainRequester.EXPECT().
				GetAgentCollectionV1MintProposalCreatedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)

			suite.mockAgentCollectionBlockchainRequester.EXPECT().
				GetAgentCollectionV1MintProposalRemovedEvent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should succeed without enqueuing any events", func(t *testing.T) {
			t.Parallel()
			basics := WhenEnqueuingBlockEventsGapBeforeEach(t)
			initSuite(basics)

			err := basics.sut.EnqueueGap()

			assert.Nil(t, err)
		})
	})
}
