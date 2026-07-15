package services_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenProcessingBlockEventsTestingSuite struct {
	syncCheckpointRepositoryMock          *repositories_mocks.MockSyncCheckpointRepositoryInterface
	mintedEventHandlerMock                *interfaces_mocks.MockMintedEventHandlerInterface
	transferEventHandlerMock              *interfaces_mocks.MockTransferEventHandlerInterface
	agentImageUpdatedEventHandlerMock     *interfaces_mocks.MockAgentImageUpdatedEventHandlerInterface
	mintProposalCreatedEventHandlerMock   *interfaces_mocks.MockMintProposalCreatedEventHandlerInterface
	mintProposalRemovedEventHandlerMock   *interfaces_mocks.MockMintProposalRemovedEventHandlerInterface
	agentAttributeUpdatedEventHandlerMock *interfaces_mocks.MockAgentAttributeUpdatedEventHandlerInterface
	collectionCreatedEventHandlerMock     *interfaces_mocks.MockCollectionCreatedEventHandlerInterface

	sut *services.BlockEventsProcessor
}

func WhenProcessingBlockEventsBeforeEach(t *testing.T) *WhenProcessingBlockEventsTestingSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	syncCheckpointRepositoryMock := repositories_mocks.NewMockSyncCheckpointRepositoryInterface(mockController)
	mintedEventHandlerMock := interfaces_mocks.NewMockMintedEventHandlerInterface(mockController)
	transferEventHandlerMock := interfaces_mocks.NewMockTransferEventHandlerInterface(mockController)
	agentImageUpdatedEventHandlerMock := interfaces_mocks.NewMockAgentImageUpdatedEventHandlerInterface(mockController)
	mintProposalCreatedEventHandlerMock := interfaces_mocks.NewMockMintProposalCreatedEventHandlerInterface(mockController)
	mintProposalRemovedEventHandlerMock := interfaces_mocks.NewMockMintProposalRemovedEventHandlerInterface(mockController)
	agentAttributeUpdatedEventHandlerMock := interfaces_mocks.NewMockAgentAttributeUpdatedEventHandlerInterface(mockController)
	collectionCreatedEventHandlerMock := interfaces_mocks.NewMockCollectionCreatedEventHandlerInterface(mockController)

	sut := services.NewBlockEventsProcessor(
		services.NewBlockEventsProcessorParams{
			Logger:                                         zap.NewNop(),
			SyncCheckpointRepository:                       syncCheckpointRepositoryMock,
			AgentCollectionV1MintedHandler:                 mintedEventHandlerMock,
			AgentCollectionV1TransferHandler:               transferEventHandlerMock,
			AgentCollectionV1AgentImageUpdatedHandler:      agentImageUpdatedEventHandlerMock,
			AgentCollectionV1MintProposalCreatedHandler:    mintProposalCreatedEventHandlerMock,
			AgentCollectionV1MintProposalRemovedHandler:    mintProposalRemovedEventHandlerMock,
			AgentCollectionV1AgentAttributeUpdatedHandler:  agentAttributeUpdatedEventHandlerMock,
			AgentCollectionManagerCollectionCreatedHandler: collectionCreatedEventHandlerMock,
		},
	)

	return &WhenProcessingBlockEventsTestingSuite{
		syncCheckpointRepositoryMock:          syncCheckpointRepositoryMock,
		mintedEventHandlerMock:                mintedEventHandlerMock,
		transferEventHandlerMock:              transferEventHandlerMock,
		agentImageUpdatedEventHandlerMock:     agentImageUpdatedEventHandlerMock,
		mintProposalCreatedEventHandlerMock:   mintProposalCreatedEventHandlerMock,
		mintProposalRemovedEventHandlerMock:   mintProposalRemovedEventHandlerMock,
		agentAttributeUpdatedEventHandlerMock: agentAttributeUpdatedEventHandlerMock,
		collectionCreatedEventHandlerMock:     collectionCreatedEventHandlerMock,

		sut: sut,
	}
}

func TestWhenProcessingBlockEvents(t *testing.T) {
	t.Parallel()

	newQueueItem := func() *service_interfaces.QueueItem {
		return &service_interfaces.QueueItem{
			AgentCollectionsManagerCollectionCreatedEvents: []*abi.AgentCollectionsManagerCollectionCreated{
				{
					CollectionId:      big.NewInt(1),
					CollectionAddress: common_tests.RandomEthAddress(),
					Name:              "Collection",
					Symbol:            "COL",
					Raw:               types.Log{},
				},
			},
		}
	}

	t.Run("Given processing block events fails once then succeeds on retry", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenProcessingBlockEventsTestingSuite, item *service_interfaces.QueueItem) {
			suite.collectionCreatedEventHandlerMock.
				EXPECT().
				Handle(item.AgentCollectionsManagerCollectionCreatedEvents[0]).
				Return(assert.AnError).
				Times(1)
			suite.collectionCreatedEventHandlerMock.
				EXPECT().
				Handle(item.AgentCollectionsManagerCollectionCreatedEvents[0]).
				Return(nil).
				Times(1)
			suite.syncCheckpointRepositoryMock.
				EXPECT().
				Get().
				Return(nil, nil)
			suite.syncCheckpointRepositoryMock.
				EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(entity *models.SyncCheckpoint) (*models.SyncCheckpoint, error) {
					assert.Equal(t, uint64(42), entity.Block)
					return entity, nil
				})
		}

		t.Run("Should retry then save checkpoint", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingBlockEventsBeforeEach(t)
			item := newQueueItem()
			initSuite(suite, item)

			err := suite.sut.ProcessBlock(context.Background(), 42, item)

			assert.NoError(t, err)
		})
	})

	t.Run("Given processing collection-created events keeps failing", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenProcessingBlockEventsTestingSuite, item *service_interfaces.QueueItem) {
			suite.collectionCreatedEventHandlerMock.
				EXPECT().
				Handle(item.AgentCollectionsManagerCollectionCreatedEvents[0]).
				Return(assert.AnError).
				Times(3)
		}

		t.Run("Should fail after max retry attempts", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingBlockEventsBeforeEach(t)
			item := newQueueItem()
			initSuite(suite, item)

			err := suite.sut.ProcessBlock(context.Background(), 77, item)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "block 77: failed to process block events")
		})
	})
}
