package services_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSynchronizingSlackChannelTestingSuite struct {
	sut *services.SlackChannelSynchronizer

	mockSlackChannelRequester  *interfaces_mocks.MockSlackChannelRequesterInterface
	mockSlackChannelRepository *repositories_mocks.MockSlackChannelRepositoryInterface
	mockSlackUserRepository    *repositories_mocks.MockSlackUserRepositoryInterface
	mockSlackUserSynchronizer  *interfaces_mocks.MockSlackUserSynchronizerInterface
}

func WhenSynchronizingSlackChannelBeforeEach(t *testing.T) *WhenSynchronizingSlackChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackChannelRequester := interfaces_mocks.NewMockSlackChannelRequesterInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)
	mockSlackUserRepository := repositories_mocks.NewMockSlackUserRepositoryInterface(mockController)
	mockSlackUserSynchronizer := interfaces_mocks.NewMockSlackUserSynchronizerInterface(mockController)

	sut := services.NewSlackChannelSynchronizer(
		zap.NewNop(),
		mockSlackChannelRequester,
		mockSlackChannelRepository,
		mockSlackUserRepository,
		mockSlackUserSynchronizer,
	)

	return &WhenSynchronizingSlackChannelTestingSuite{
		sut: sut,

		mockSlackChannelRequester:  mockSlackChannelRequester,
		mockSlackChannelRepository: mockSlackChannelRepository,
		mockSlackUserRepository:    mockSlackUserRepository,
		mockSlackUserSynchronizer:  mockSlackUserSynchronizer,
	}
}

func TestWhenSynchronizingSlackChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting channel information from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), "channelId")
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting channel information from slack", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(slackChannel, nil)
			suite.mockSlackChannelRequester.EXPECT().
				GetChannelInformation(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), "channelId")
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while joining the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, nil)
			suite.mockSlackChannelRequester.EXPECT().
				GetChannelInformation(gomock.Any(), gomock.Any()).
				Return(&interfaces.ChannelInformation{IsMember: false}, nil)
			suite.mockSlackChannelRequester.EXPECT().JoinChannel(gomock.Any(), gomock.Any()).Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), "channelId")
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving channel information to the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(slackChannel, nil)
			suite.mockSlackChannelRequester.EXPECT().
				GetChannelInformation(gomock.Any(), gomock.Any()).
				Return(&interfaces.ChannelInformation{
					IsMember:    true,
					Id:          slackChannel.Id,
					Name:        slackChannel.Name,
					Description: slackChannel.Description,
				}, nil)
			suite.mockSlackChannelRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), "channelId")
			assert.Error(t, err)
		})
	})

	t.Run("Given the channel has no members", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(slackChannel, nil)
			suite.mockSlackChannelRequester.EXPECT().
				GetChannelInformation(gomock.Any(), gomock.Any()).
				Return(&interfaces.ChannelInformation{
					IsMember:    true,
					Id:          slackChannel.Id,
					Name:        slackChannel.Name,
					Description: slackChannel.Description,
					UserIds:     nil,
				}, nil)
			suite.mockSlackChannelRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(slackChannel, nil)
		}

		t.Run("Should not touch the user repository or the user synchronizer", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), "channelId")
			assert.NoError(t, err)
		})
	})

	t.Run("Given all channel members already exist in the database", func(t *testing.T) {
		t.Parallel()

		channelId := "C-target"
		userIds := []string{"U1", "U2"}

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().WithId(channelId).Build()
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), channelId).Return(slackChannel, nil)
			suite.mockSlackChannelRequester.EXPECT().
				GetChannelInformation(gomock.Any(), channelId).
				Return(&interfaces.ChannelInformation{
					IsMember:    true,
					Id:          channelId,
					Name:        slackChannel.Name,
					Description: slackChannel.Description,
					UserIds:     userIds,
				}, nil)
			suite.mockSlackChannelRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(slackChannel, nil)

			u1 := model_builders.NewSlackUserBuilder().WithId("U1").Build()
			u2 := model_builders.NewSlackUserBuilder().WithId("U2").Build()
			suite.mockSlackUserRepository.EXPECT().
				ListByIds(gomock.Any(), userIds).
				Return([]*models.SlackUser{u1, u2}, nil)

			suite.mockSlackChannelRepository.EXPECT().AddUser(gomock.Any(), channelId, "U1").Return(nil)
			suite.mockSlackChannelRepository.EXPECT().AddUser(gomock.Any(), channelId, "U2").Return(nil)
		}

		t.Run("Should attach each member synchronously without syncing them from Slack", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), channelId)
			assert.NoError(t, err)
		})
	})
}
