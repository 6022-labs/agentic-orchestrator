package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSynchronizingUserTestingSuite struct {
	sut *services.SlackUserSynchronizer

	mockSlackUserRequester  *interfaces_mocks.MockSlackUserRequesterInterface
	mockSlackUserRepository *repositories_mocks.MockSlackUserRepositoryInterface
}

func WhenSynchronizingUserBeforeEach(t *testing.T) *WhenSynchronizingUserTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackUserRequester := interfaces_mocks.NewMockSlackUserRequesterInterface(mockController)
	mockSlackUserRepository := repositories_mocks.NewMockSlackUserRepositoryInterface(mockController)

	sut := services.NewSlackUserSynchronizer(
		zap.NewNop(),
		mockSlackUserRequester,
		mockSlackUserRepository,
	)

	return &WhenSynchronizingUserTestingSuite{
		sut: sut,

		mockSlackUserRequester:  mockSlackUserRequester,
		mockSlackUserRepository: mockSlackUserRepository,
	}
}

func TestWhenSynchronizingUser(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting user from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("some-user-id")
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting user from Slack", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil)

			suite.mockSlackUserRequester.
				EXPECT().
				GetUserInformation(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("some-user-id")
			assert.Error(t, err)
		})
	})

	t.Run("Given user does not exist in Slack", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil)

			suite.mockSlackUserRequester.
				EXPECT().
				GetUserInformation(gomock.Any()).
				Return(nil, nil)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("some-user-id")
			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving user to the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserRequester.
				EXPECT().
				GetUserInformation(gomock.Any()).
				Return(slackUser, nil)

			suite.mockSlackUserRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("some-user-id")
			assert.Error(t, err)
		})
	})

	t.Run("Given user exists in Slack and does not exist in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(nil, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserRequester.
				EXPECT().
				GetUserInformation(gomock.Any()).
				Return(slackUser, nil)

			suite.mockSlackUserRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(slackUser, nil)
		}

		t.Run("Should create the user in the database and return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			user, err := suite.sut.Sync("some-user-id")
			assert.NoError(t, err)
			assert.Equal(t, "some-user-id", user.Id)
			assert.Equal(t, "some-name", user.Name)
			assert.Equal(t, "Some Real Name", user.RealName)
		})
	})

	t.Run("Given user exists in Slack and in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			slackUser := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserRepository.
				EXPECT().
				GetById(gomock.Any()).
				Return(slackUser, nil)

			suite.mockSlackUserRequester.
				EXPECT().
				GetUserInformation(gomock.Any()).
				Return(slackUser, nil)

			suite.mockSlackUserRepository.
				EXPECT().
				Save(gomock.Any()).
				Return(slackUser, nil)
		}

		t.Run("Should update the user in the database and return it", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			user, err := suite.sut.Sync("some-user-id")
			assert.NoError(t, err)
			assert.Equal(t, "some-user-id", user.Id)
			assert.Equal(t, "some-name", user.Name)
			assert.Equal(t, "Some Real Name", user.RealName)
		})
	})
}
