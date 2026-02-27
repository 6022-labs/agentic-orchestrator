package command_handlers_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenRemovingAgentFromChannelTestingSuite struct {
	sut *command_handlers.RemoveAgentFromChannel

	mockSlackAgentRepository     *repositories_mocks.MockSlackAgentRepositoryInterface
	mockSlackChannelSynchronizer *interfaces_mocks.MockSlackChannelSynchronizerInterface
	mockSlackChannelRepository   *repositories_mocks.MockSlackChannelRepositoryInterface
}

func WhenRemovingAgentFromChannelBeforeEach(t *testing.T) *WhenRemovingAgentFromChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)
	mockSlackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)

	sut := command_handlers.NewRemoveAgentFromChannel(
		zap.NewNop(),
		mockSlackAgentRepository,
		mockSlackChannelRepository,
		mockSlackChannelSynchronizer,
	)

	return &WhenRemovingAgentFromChannelTestingSuite{
		sut:                          sut,
		mockSlackAgentRepository:     mockSlackAgentRepository,
		mockSlackChannelSynchronizer: mockSlackChannelSynchronizer,
		mockSlackChannelRepository:   mockSlackChannelRepository,
	}
}

func TestWhenRemovingAgentFromChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting the channel from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the channel did not contain the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)
		}

		t.Run("Should not attempt to remove the agent again and return no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an error occurs while syncing the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(nil, errors.New("sync error"))
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent from repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, request.CollectionTokenId).
				Return(nil, errors.New("repo error"))
		}

		t.Run("Should return repo error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent not found in repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, request.CollectionTokenId).
				Return(nil, gorm.ErrRecordNotFound)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given error occurs while saving channel after removal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				RemoveAgent(slackChannel, slackAgent).
				Return(nil, errors.New("save error"))
		}

		t.Run("Should return save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given all steps succeed", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				RemoveAgent(slackChannel, slackAgent).
				Return(slackChannel, nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(request)

			assert.NoError(t, err)
		})
	})
}
