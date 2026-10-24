package services_test

import (
	"context"
	"strings"

	"math/big"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingCollectionCreatedEventTestingSuite struct {
	agentCollectionRepositoryMock *repositories_mocks.MockAgentCollectionRepositoryInterface

	sut *services.CollectionCreatedEventHandler
}

func WhenHandlingCollectionCreatedEventBeforeEach(t *testing.T) *WhenHandlingCollectionCreatedEventTestingSuite {
	mockController := gomock.NewController(t)

	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)

	params := services.NewCollectionCreatedEventHandlerParams{
		Logger: zap.NewNop(),

		AgentCollectionRepository: agentCollectionRepositoryMock,
	}

	sut := services.NewCollectionCreatedEventHandler(
		params,
	)

	return &WhenHandlingCollectionCreatedEventTestingSuite{
		agentCollectionRepositoryMock: agentCollectionRepositoryMock,

		sut: sut,
	}
}

func TestWhenHandlingCollectionCreatedEvent(t *testing.T) {
	t.Parallel()

	newEvent := func(collectionID uint64, name, symbol string) *abi.AgentCollectionsManagerCollectionCreated {
		return &abi.AgentCollectionsManagerCollectionCreated{
			CollectionId:      new(big.Int).SetUint64(collectionID),
			CollectionAddress: common_tests.RandomEthAddress(),
			Name:              name,
			Symbol:            symbol,
			Raw:               types.Log{},
		}
	}

	t.Run("Given saving collection fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCollectionCreatedEventTestingSuite) *abi.AgentCollectionsManagerCollectionCreated {
			event := newEvent(1, "Collection", "COL")
			suite.agentCollectionRepositoryMock.
				EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)

			return event
		}

		t.Run("Should return wrapped save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCollectionCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "failed to save new collection")
		})
	})

	t.Run("Given event is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCollectionCreatedEventTestingSuite) *abi.AgentCollectionsManagerCollectionCreated {
			event := newEvent(2, "My Collection", "MC")
			suite.agentCollectionRepositoryMock.
				EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, collection *models.AgentCollection) (*models.AgentCollection, error) {
					assert.Equal(t, strings.ToLower(event.CollectionAddress.Hex()), collection.Address)
					assert.Equal(t, event.Name, collection.Name)
					assert.Equal(t, event.Symbol, collection.Symbol)
					return collection, nil
				})

			return event
		}

		t.Run("Should map event fields and save", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCollectionCreatedEventBeforeEach(t)
			event := initSuite(suite)

			err := suite.sut.Handle(context.Background(), event)

			assert.NoError(t, err)
		})
	})
}
