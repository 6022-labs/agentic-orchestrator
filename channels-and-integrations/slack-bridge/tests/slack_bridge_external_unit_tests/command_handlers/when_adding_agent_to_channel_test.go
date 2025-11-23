package command_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenAddingAgentToChannelTestingSuite struct {
	sut *command_handlers.AddAgentToChannel

	mockSlackAgentSynchronizer   *interfaces_mocks.MockSlackAgentSynchronizerInterface
	mockSlackChannelSynchronizer *interfaces_mocks.MockSlackChannelSynchronizerInterface
	mockSlackChannelRepository   *repositories_mocks.MockSlackChannelRepositoryInterface
}

func WhenAddingAgentToChannelBeforeEach(t *testing.T) *WhenAddingAgentToChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackAgentSynchronizer := interfaces_mocks.NewMockSlackAgentSynchronizerInterface(mockController)
	mockSlackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)

	sut := command_handlers.NewAddAgentToChannel(
		mockSlackAgentSynchronizer,
		mockSlackChannelRepository,
		mockSlackChannelSynchronizer,
	)

	return &WhenAddingAgentToChannelTestingSuite{
		sut: sut,

		mockSlackAgentSynchronizer:   mockSlackAgentSynchronizer,
		mockSlackChannelRepository:   mockSlackChannelRepository,
		mockSlackChannelSynchronizer: mockSlackChannelSynchronizer,
	}
}

func TestWhenAddingAgentToChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting the channel from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the channel already contains the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*model_builders.NewSlackAgentBuilder().
						WithCollectionAddress(request.CollectionAddress).
						WithCollectionTokenId(request.CollectionTokenId).
						Build(),
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)
		}

		t.Run("Should not attempt to add the agent again and return no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an error occurs while syncing the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while syncing agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentSynchronizer.
				EXPECT().
				Sync(request.CollectionAddress, request.CollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return sync agent error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while adding agent to channel in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(slackChannel, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().Build()
			suite.mockSlackAgentSynchronizer.
				EXPECT().
				Sync(request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				AddAgent(slackChannel, slackAgent).
				Return(nil, assert.AnError)
		}

		t.Run("Should return save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given all steps succeed", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(slackChannel, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().Build()
			suite.mockSlackAgentSynchronizer.
				EXPECT().
				Sync(request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				AddAgent(slackChannel, slackAgent).
				Return(slackChannel, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(request)

			assert.NoError(t, err)
		})
	})
}
