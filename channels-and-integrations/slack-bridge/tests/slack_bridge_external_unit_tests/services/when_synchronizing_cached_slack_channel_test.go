package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	external_interfaces_mocks "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSynchronizingCachedSlackChannelTestingSuite struct {
	sut *services.SlackChannelSynchronizerCached

	cacheSettings                *settings.CacheSettings
	mockCacheStorage             *interfaces_mocks.MockCacheStorageInterface
	mockSlackChannelSynchronizer *external_interfaces_mocks.MockSlackChannelSynchronizerInterface
	mockSlackChannelRepository   *repositories_mocks.MockSlackChannelRepositoryInterface
}

func WhenSynchronizingCachedSlackChannelBeforeEach(t *testing.T) *WhenSynchronizingCachedSlackChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)
	mockSlackChannelSynchronizer := external_interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)

	cacheSettings := &settings.CacheSettings{
		SlackChannelSynchronizerTTLMinutes: 10,
	}

	sut := services.NewSlackChannelSynchronizerCached(
		zap.NewNop(),
		cacheSettings,
		mockCacheStorage,
		mockSlackChannelRepository,
		mockSlackChannelSynchronizer,
	)

	return &WhenSynchronizingCachedSlackChannelTestingSuite{
		sut: sut,

		cacheSettings:                cacheSettings,
		mockCacheStorage:             mockCacheStorage,
		mockSlackChannelRepository:   mockSlackChannelRepository,
		mockSlackChannelSynchronizer: mockSlackChannelSynchronizer,
	}
}

func TestWhenSynchronizingCachedSlackChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting channel information from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackChannelTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackChannelBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("C1234567890")
			assert.Error(t, err)
		})
	})

	t.Run("Given channel to sync is not in cache and channel does not exist in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackChannelTestingSuite) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil)

			suite.mockCacheStorage.
				EXPECT().
				ReadBool(gomock.Any()).
				Return(nil, nil)

			channel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any()).
				Return(channel, nil)

			suite.mockCacheStorage.
				EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any())
		}

		t.Run("Should sync the channel and return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackChannelBeforeEach(t)
			initSuite(suite)

			channel, err := suite.sut.Sync("C1234567890")

			assert.NoError(t, err)
			assert.NotNil(t, channel)
		})
	})

	t.Run("Given channel to sync is not in cache and channel exists in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackChannelTestingSuite) {
			existingChannel := model_builders.NewSlackChannelBuilder().
				WithId("C1234567890").
				WithName("existing-channel").
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(existingChannel, nil)

			suite.mockCacheStorage.
				EXPECT().
				ReadBool(gomock.Any()).
				Return(nil, nil)

			updatedChannel := model_builders.NewSlackChannelBuilder().
				WithId("C1234567890").
				WithName("updated-channel").
				Build()
			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any()).
				Return(updatedChannel, nil)

			suite.mockCacheStorage.
				EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any())
		}

		t.Run("Should sync the channel and return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackChannelBeforeEach(t)
			initSuite(suite)

			channel, err := suite.sut.Sync("C1234567890")

			assert.NoError(t, err)
			assert.NotNil(t, channel)
			assert.Equal(t, "updated-channel", channel.Name)
		})
	})

	t.Run("Given channel to sync is in cache", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackChannelTestingSuite) {
			existingChannel := model_builders.NewSlackChannelBuilder().
				WithId("C1234567890").
				WithName("existing-channel").
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(existingChannel, nil)

			cacheResult := true
			suite.mockCacheStorage.
				EXPECT().
				ReadBool(gomock.Any()).
				Return(&cacheResult, nil)
		}

		t.Run("Should return the channel from the database without syncing", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackChannelBeforeEach(t)
			initSuite(suite)

			channel, err := suite.sut.Sync("C1234567890")

			assert.NoError(t, err)
			assert.NotNil(t, channel)
			assert.Equal(t, "existing-channel", channel.Name)
		})
	})
}
