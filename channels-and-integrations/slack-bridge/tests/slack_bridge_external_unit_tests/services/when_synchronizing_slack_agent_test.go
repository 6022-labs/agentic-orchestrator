package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSynchronizingSlackAgentTestingSuite struct {
	sut *services.SlackAgentSynchronizer

	mockAgentRequester       *interfaces_mocks.MockAgentRequesterInterface
	mockSlackAgentRepository *repositories_mocks.MockSlackAgentRepositoryInterface
}

func WhenSynchronizingSlackAgentBeforeEach(t *testing.T) *WhenSynchronizingSlackAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRequester := interfaces_mocks.NewMockAgentRequesterInterface(mockController)
	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)

	sut := services.NewSlackAgentSynchronizer(
		mockAgentRequester,
		mockSlackAgentRepository,
	)

	return &WhenSynchronizingSlackAgentTestingSuite{
		sut: sut,

		mockAgentRequester:       mockAgentRequester,
		mockSlackAgentRepository: mockSlackAgentRepository,
	}
}

func TestWhenSynchronizingSlackAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting agent from orchestrator", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite) {
			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackAgentBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(common_tests.RandomEthAddress().Hex(), 123)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite) {
			agentFromOrchestrator := model_builders.NewSlackAgentBuilder().Build()
			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				Return(agentFromOrchestrator, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackAgentBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(common_tests.RandomEthAddress().Hex(), 123)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is found in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite, collectionAddress string, collectionTokenId uint64) {
			updatedCid := "initial-cid"
			agentFromOrchestrator := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				WithIconImageCid(&updatedCid).
				WithOwnerWallet(common_tests.RandomEthAddress().Hex()).
				Build()

			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				Return(agentFromOrchestrator, nil)

			agentInDatabase := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				WithIconImageCid(nil).
				WithOwnerWallet(common_tests.RandomEthAddress().Hex()).
				Build()
			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(agentInDatabase, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				UpdateOwnerWallet(collectionAddress, collectionTokenId, agentFromOrchestrator.OwnerWallet).
				Return(nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				UpdateIconImageCid(collectionAddress, collectionTokenId, agentFromOrchestrator.IconImageCid).
				Return(nil)
		}

		t.Run("Should return the updated slack agent from the database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackAgentBeforeEach(t)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(123)

			initSuite(suite, collectionAddress, collectionTokenId)

			agent, err := suite.sut.Sync(collectionAddress, collectionTokenId)

			assert.Nil(t, err)
			assert.NotNil(t, agent)
		})
	})

	t.Run("Given agent is not found in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite) {
			agentFromOrchestrator := model_builders.NewSlackAgentBuilder().Build()
			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				Return(agentFromOrchestrator, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any()).
				Return(nil, nil)

			newAgent := model_builders.NewSlackAgentBuilder().Build()
			suite.mockSlackAgentRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(newAgent, nil)
		}

		t.Run("Should return the new slack agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackAgentBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(common_tests.RandomEthAddress().Hex(), 123)

			assert.Nil(t, err)
		})
	})
}
