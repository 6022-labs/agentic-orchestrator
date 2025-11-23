package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/orchestrator_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSynchronizingAgentCollectionByAddressTestingSuite struct {
	sut *services.AgentCollectionSynchronizer

	agentCollectionRequesterMock  *interfaces_mocks.MockAgentCollectionRequesterInterface
	agentCollectionRepositoryMock *repositories_mocks.MockAgentCollectionRepositoryInterface
}

func WhenSynchronizingAgentCollectionByAddressBeforeEach(t *testing.T) *WhenSynchronizingAgentCollectionByAddressTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentCollectionRequester := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)

	sut := services.NewAgentCollectionSynchronizer(
		mockAgentCollectionRequester,
		agentCollectionRepositoryMock,
	)

	return &WhenSynchronizingAgentCollectionByAddressTestingSuite{
		sut: sut,

		agentCollectionRequesterMock:  mockAgentCollectionRequester,
		agentCollectionRepositoryMock: agentCollectionRepositoryMock,
	}
}

func TestWhenSynchronizingAgentCollectionByAddress(t *testing.T) {
	t.Parallel()

	t.Run("Given invalid collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentCollectionByAddressBeforeEach(t)

			invalidAddress := "invalid_address"
			err := suite.sut.SyncByCollectionAddress(invalidAddress)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting collection name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentCollectionByAddressBeforeEach(t)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			suite.agentCollectionRequesterMock.EXPECT().GetName(gomock.Any()).Return(nil, assert.AnError)

			err := suite.sut.SyncByCollectionAddress(collectionAddress)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting collection symbol", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentCollectionByAddressBeforeEach(t)

			collectionAddress := common_tests.RandomEthAddress().Hex()

			name := "Test Collection"
			suite.agentCollectionRequesterMock.EXPECT().GetName(gomock.Any()).Return(&name, nil)
			suite.agentCollectionRequesterMock.EXPECT().GetSymbol(gomock.Any()).Return(nil, assert.AnError)

			err := suite.sut.SyncByCollectionAddress(collectionAddress)
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving collection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentCollectionByAddressBeforeEach(t)

			collectionAddress := common_tests.RandomEthAddress().Hex()

			name := "Test Collection"
			symbol := "TST"
			suite.agentCollectionRequesterMock.EXPECT().GetName(gomock.Any()).Return(&name, nil)
			suite.agentCollectionRequesterMock.EXPECT().GetSymbol(gomock.Any()).Return(&symbol, nil)
			suite.agentCollectionRepositoryMock.EXPECT().Save(gomock.Any()).Return(nil, assert.AnError)

			err := suite.sut.SyncByCollectionAddress(collectionAddress)
			assert.Error(t, err)
		})
	})

	t.Run("Given no errors", func(t *testing.T) {
		t.Parallel()

		t.Run("Should save collection", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentCollectionByAddressBeforeEach(t)

			collectionAddress := common_tests.RandomEthAddress().Hex()

			name := "Test Collection"
			symbol := "TST"
			suite.agentCollectionRequesterMock.EXPECT().GetName(gomock.Any()).Return(&name, nil)
			suite.agentCollectionRequesterMock.EXPECT().GetSymbol(gomock.Any()).Return(&symbol, nil)
			suite.agentCollectionRepositoryMock.EXPECT().Save(gomock.Any()).Return(&models.AgentCollection{}, nil)

			err := suite.sut.SyncByCollectionAddress(collectionAddress)
			assert.NoError(t, err)
		})
	})
}
