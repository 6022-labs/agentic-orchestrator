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

type WhenSynchronizingAgentCollectionTestingSuite struct {
	sut *services.AgentCollectionSynchronizer

	agentCollectionRepositoryMock *repositories_mocks.MockAgentCollectionRepositoryInterface
}

func WhenSynchronizingAgentCollectionBeforeEach(t *testing.T) *WhenSynchronizingAgentCollectionTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentCollectionRequester := interfaces_mocks.NewMockAgentCollectionRequesterInterface(mockController)
	agentCollectionRepositoryMock := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)

	sut := services.NewAgentCollectionSynchronizer(
		mockAgentCollectionRequester,
		agentCollectionRepositoryMock,
	)

	return &WhenSynchronizingAgentCollectionTestingSuite{
		sut: sut,

		agentCollectionRepositoryMock: agentCollectionRepositoryMock,
	}
}

func TestWhenSynchronizingAgentCollection(t *testing.T) {
	t.Parallel()

	t.Run("Given repository returns error", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentCollectionBeforeEach(t)

			suite.agentCollectionRepositoryMock.EXPECT().Save(gomock.Any()).Return(nil, assert.AnError)

			address := common_tests.RandomEthAddress().Hex()
			err := suite.sut.Sync(address, "test", "TST")
			assert.Error(t, err)
		})
	})

	t.Run("Given no errors", func(t *testing.T) {
		t.Parallel()

		t.Run("Should save collection", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingAgentCollectionBeforeEach(t)

			suite.agentCollectionRepositoryMock.EXPECT().Save(gomock.Any()).Return(&models.AgentCollection{}, nil)

			address := common_tests.RandomEthAddress().Hex()
			err := suite.sut.Sync(address, "test", "TST")
			assert.NoError(t, err)
		})
	})
}
