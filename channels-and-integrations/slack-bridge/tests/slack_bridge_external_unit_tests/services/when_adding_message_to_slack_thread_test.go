package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenAddingMessageToSlackThreadTestingSuite struct {
	sut *services.SlackThreadMessagesAdder

	mockSlackThreadRepository  *repositories_mocks.MockSlackThreadRepositoryInterface
	mockSlackMessageRepository *repositories_mocks.MockSlackMessageRepositoryInterface
	mockSlackChannelRepository *repositories_mocks.MockSlackChannelRepositoryInterface
}

func WhenAddingMessageToSlackThreadBeforeEach(t *testing.T) *WhenAddingMessageToSlackThreadTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackThreadRepository := repositories_mocks.NewMockSlackThreadRepositoryInterface(mockController)
	mockSlackMessageRepository := repositories_mocks.NewMockSlackMessageRepositoryInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)

	sut := services.NewSlackThreadMessagesAdder(
		zap.NewNop(),
		mockSlackThreadRepository,
		mockSlackMessageRepository,
		mockSlackChannelRepository,
	)

	return &WhenAddingMessageToSlackThreadTestingSuite{
		sut: sut,

		mockSlackThreadRepository:  mockSlackThreadRepository,
		mockSlackMessageRepository: mockSlackMessageRepository,
		mockSlackChannelRepository: mockSlackChannelRepository,
	}
}

func TestWhenAddingMessageToSlackThread(t *testing.T) {
	t.Parallel()

	t.Run("Given neither authorUserId nor both authorCollectionAddress and authorCollectionTokenId are provided", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToSlackThreadTestingSuite) {
			// No initialization needed for this test case
		}

		t.Run("Should return an error indicating that author information is missing", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToSlackThreadBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Add("channelId", "parentMessageTs", "messageTs", "text", nil, nil, nil)

			assert.Error(t, err)
			assert.Equal(t, "either authorUserId or both authorCollectionAddress and authorCollectionTokenId must be provided", err.Error())
		})
	})

	t.Run("Given error occurs while getting slack channel from database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToSlackThreadTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToSlackThreadBeforeEach(t)
			initSuite(suite)

			authorUserId := "authorUserId"
			err := suite.sut.Add("channelId", "parentMessageTs", "messageTs", "text", &authorUserId, nil, nil)

			assert.Error(t, err)
		})
	})

	t.Run("Given channel is not found in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToSlackThreadTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return an error indicating channel not found", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToSlackThreadBeforeEach(t)
			initSuite(suite)

			authorUserId := "authorUserId"
			err := suite.sut.Add("channelId", "parentMessageTs", "messageTs", "text", &authorUserId, nil, nil)

			assert.Error(t, err)
			assert.Equal(t, "slack channel with ID channelId not found", err.Error())
		})
	})

	t.Run("Given error occurs while getting thread information from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToSlackThreadTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackChannel, nil)

			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToSlackThreadBeforeEach(t)
			initSuite(suite)

			authorUserId := "authorUserId"
			err := suite.sut.Add("channelId", "parentMessageTs", "messageTs", "text", &authorUserId, nil, nil)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving thread information to the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToSlackThreadTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackChannel, nil)

			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId(gomock.Any(), gomock.Any()).
				Return(nil, nil)

			suite.mockSlackThreadRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToSlackThreadBeforeEach(t)

			initSuite(suite)

			authorUserId := "authorUserId"
			err := suite.sut.Add("channelId", "parentMessageTs", "messageTs", "text", &authorUserId, nil, nil)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving messages to the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToSlackThreadTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackChannel, nil)

			slackThread := model_builders.NewSlackThreadBuilder().
				Build()
			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId(gomock.Any(), gomock.Any()).
				Return(slackThread, nil)

			suite.mockSlackMessageRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToSlackThreadBeforeEach(t)

			initSuite(suite)

			authorUserId := "authorUserId"
			err := suite.sut.Add("channelId", "parentMessageTs", "messageTs", "text", &authorUserId, nil, nil)

			assert.Error(t, err)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingMessageToSlackThreadTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackChannel, nil)

			slackThread := model_builders.NewSlackThreadBuilder().
				Build()
			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId(gomock.Any(), gomock.Any()).
				Return(slackThread, nil)

			suite.mockSlackMessageRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return slack thread", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingMessageToSlackThreadBeforeEach(t)

			initSuite(suite)

			authorUserId := "authorUserId"
			err := suite.sut.Add("channelId", "parentMessageTs", "messageTs", "text", &authorUserId, nil, nil)

			assert.NoError(t, err)
		})
	})
}
