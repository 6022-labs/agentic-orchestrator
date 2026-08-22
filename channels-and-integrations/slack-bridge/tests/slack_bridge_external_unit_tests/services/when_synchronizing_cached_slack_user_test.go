package services_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	external_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSynchronizingCachedSlackUserTestingSuite struct {
	sut *services.SlackUserSynchronizerCached

	cacheSettings             *settings.CacheSettings
	mockCacheStorage          *interfaces_mocks.MockCacheStorageInterface
	mockSlackUserRepository   *repositories_mocks.MockSlackUserRepositoryInterface
	mockSlackUserSynchronizer *external_interfaces_mocks.MockSlackUserSynchronizerInterface
}

func WhenSynchronizingCachedSlackUserBeforeEach(t *testing.T) *WhenSynchronizingCachedSlackUserTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockSlackUserRepository := repositories_mocks.NewMockSlackUserRepositoryInterface(mockController)
	mockSlackUserSynchronizer := external_interfaces_mocks.NewMockSlackUserSynchronizerInterface(mockController)

	cacheSettings := &settings.CacheSettings{
		SlackUserSynchronizerTTLMinutes: 10,
	}

	sut := services.NewSlackUserSynchronizerCached(
		zap.NewNop(),
		cacheSettings,
		mockCacheStorage,
		mockSlackUserRepository,
		mockSlackUserSynchronizer,
	)

	return &WhenSynchronizingCachedSlackUserTestingSuite{
		sut: sut,

		cacheSettings:             cacheSettings,
		mockCacheStorage:          mockCacheStorage,
		mockSlackUserRepository:   mockSlackUserRepository,
		mockSlackUserSynchronizer: mockSlackUserSynchronizer,
	}
}

func TestWhenSynchronizingCachedSlackUser(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting user from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackUserTestingSuite) {
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackUserBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), "some-user-id")
			assert.Error(t, err)
		})
	})

	t.Run("Given user to sync is not in cache and user does not exist in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackUserTestingSuite) {
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any(), gomock.Any()).
				Return(nil, nil)

			suite.mockCacheStorage.
				EXPECT().
				ReadBool(gomock.Any(), gomock.Any()).
				Return(nil, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserSynchronizer.
				EXPECT().
				Sync(gomock.Any(), gomock.Any()).
				Return(slackUser, nil)

			suite.mockCacheStorage.
				EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		}

		t.Run("Should sync the user and return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackUserBeforeEach(t)
			initSuite(suite)

			user, err := suite.sut.Sync(context.Background(), "some-user-id")
			assert.NoError(t, err)
			assert.Equal(t, "some-user-id", user.Id)
			assert.Equal(t, "some-name", user.Name)
			assert.Equal(t, "Some Real Name", user.RealName)
		})
	})

	t.Run("Given user to sync is not in cache and user exists in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackUserTestingSuite) {
			slackUserInDb := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("old-name").
				WithRealName("Old Real Name").
				Build()
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any(), gomock.Any()).
				Return(slackUserInDb, nil)

			suite.mockCacheStorage.
				EXPECT().
				ReadBool(gomock.Any(), gomock.Any()).
				Return(nil, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserSynchronizer.
				EXPECT().
				Sync(gomock.Any(), gomock.Any()).
				Return(slackUser, nil)

			suite.mockCacheStorage.
				EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		}

		t.Run("Should sync the user and return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackUserBeforeEach(t)
			initSuite(suite)

			user, err := suite.sut.Sync(context.Background(), "some-user-id")
			assert.NoError(t, err)
			assert.Equal(t, "some-user-id", user.Id)
			assert.Equal(t, "some-name", user.Name)
			assert.Equal(t, "Some Real Name", user.RealName)
		})
	})

	t.Run("Given user to sync is in cache", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingCachedSlackUserTestingSuite) {
			cacheReturn := true
			suite.mockCacheStorage.
				EXPECT().
				ReadBool(gomock.Any(), gomock.Any()).
				Return(&cacheReturn, nil)
		}

		t.Run("Should not sync the user and return the user from the database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingCachedSlackUserBeforeEach(t)
			initSuite(suite)

			slackUserInDb := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any(), gomock.Any()).
				Return(slackUserInDb, nil)

			user, err := suite.sut.Sync(context.Background(), "some-user-id")
			assert.NoError(t, err)
			assert.Equal(t, "some-user-id", user.Id)
			assert.Equal(t, "some-name", user.Name)
			assert.Equal(t, "Some Real Name", user.RealName)
		})
	})
}
