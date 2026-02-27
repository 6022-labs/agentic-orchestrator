package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenSynchronizingSlackChannelTestingSuite struct {
	sut *services.SlackChannelSynchronizer

	mockSlackChannelRequester  *interfaces_mocks.MockSlackChannelRequesterInterface
	mockSlackChannelRepository *repositories_mocks.MockSlackChannelRepositoryInterface
}

func WhenSynchronizingSlackChannelBeforeEach(t *testing.T) *WhenSynchronizingSlackChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackChannelRequester := interfaces_mocks.NewMockSlackChannelRequesterInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)

	sut := services.NewSlackChannelSynchronizer(
		zap.NewNop(),
		mockSlackChannelRequester,
		mockSlackChannelRepository,
	)

	return &WhenSynchronizingSlackChannelTestingSuite{
		sut: sut,

		mockSlackChannelRequester:  mockSlackChannelRequester,
		mockSlackChannelRepository: mockSlackChannelRepository,
	}
}

func TestWhenSynchronizingSlackChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting channel information from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)

			initSuite(suite)

			_, err := suite.sut.Sync("channelId")

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting channel information from slack", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackChannel, nil)

			suite.mockSlackChannelRequester.
				EXPECT().
				GetChannelInformation(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)

			initSuite(suite)

			_, err := suite.sut.Sync("channelId")

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while joining the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, gorm.ErrRecordNotFound)

			suite.mockSlackChannelRequester.
				EXPECT().
				GetChannelInformation(gomock.Any()).
				Return(&interfaces.ChannelInformation{
					IsMember: false,
				}, nil)

			suite.mockSlackChannelRequester.
				EXPECT().
				JoinChannel(gomock.Any()).
				Return(assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)

			initSuite(suite)

			_, err := suite.sut.Sync("channelId")

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving channel information to the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackChannel, nil)

			suite.mockSlackChannelRequester.
				EXPECT().
				GetChannelInformation(gomock.Any()).
				Return(&interfaces.ChannelInformation{
					IsMember:    true,
					Id:          slackChannel.Id,
					Name:        slackChannel.Name,
					Members:     slackChannel.Members,
					Description: slackChannel.Description,
				}, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("channelId")

			assert.Error(t, err)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackChannelTestingSuite) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackChannel, nil)

			suite.mockSlackChannelRequester.
				EXPECT().
				GetChannelInformation(gomock.Any()).
				Return(&interfaces.ChannelInformation{
					IsMember:    true,
					Id:          slackChannel.Id,
					Name:        slackChannel.Name,
					Members:     slackChannel.Members,
					Description: slackChannel.Description,
				}, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(slackChannel, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackChannelBeforeEach(t)

			initSuite(suite)

			_, err := suite.sut.Sync("channelId")

			assert.NoError(t, err)
		})
	})
}
