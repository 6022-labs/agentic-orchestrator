package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSynchronizingUserTestingSuite struct {
	sut *services.SlackUserSynchronizer

	mockSlackUserRequester  *interfaces_mocks.MockSlackUserRequesterInterface
	mockSlackUserRepository *repositories_mocks.MockSlackUserRepositoryInterface
	mockBrokerRequester     *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenSynchronizingUserBeforeEach(t *testing.T) *WhenSynchronizingUserTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackUserRequester := interfaces_mocks.NewMockSlackUserRequesterInterface(mockController)
	mockSlackUserRepository := repositories_mocks.NewMockSlackUserRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := services.NewSlackUserSynchronizer(
		zap.NewNop(),
		mockSlackUserRequester,
		mockSlackUserRepository,
		mockBrokerRequester,
	)

	return &WhenSynchronizingUserTestingSuite{
		sut: sut,

		mockSlackUserRequester:  mockSlackUserRequester,
		mockSlackUserRepository: mockSlackUserRepository,
		mockBrokerRequester:     mockBrokerRequester,
	}
}

func TestWhenSynchronizingUser(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting user from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.EXPECT().GetById(gomock.Any()).Return(nil, assert.AnError)
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
			suite.mockSlackUserRepository.EXPECT().GetById(gomock.Any()).Return(nil, nil)
			suite.mockSlackUserRequester.EXPECT().GetUserInformation(gomock.Any()).Return(nil, assert.AnError)
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
			suite.mockSlackUserRepository.EXPECT().GetById(gomock.Any()).Return(nil, nil)
			suite.mockSlackUserRequester.EXPECT().GetUserInformation(gomock.Any()).Return(nil, nil)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("some-user-id")
			assert.Error(t, err)
		})
	})

	t.Run("Given user is new to the bridge", func(t *testing.T) {
		t.Parallel()

		brokerUserId := uuid.New()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.EXPECT().GetById(gomock.Any()).Return(nil, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserRequester.EXPECT().GetUserInformation(gomock.Any()).Return(slackUser, nil)

			suite.mockBrokerRequester.EXPECT().
				RegisterUser("Some Real Name").
				Return(brokerUserId, nil)

			suite.mockSlackUserRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(user *models.SlackUser) (*models.SlackUser, error) {
					assert.NotNil(t, user.BrokerParticipantId)
					assert.Equal(t, brokerUserId, *user.BrokerParticipantId)
					return user, nil
				})
		}

		t.Run(
			"Should register the user with the broker, store the broker participant id, and return the user",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenSynchronizingUserBeforeEach(t)
				initSuite(suite)

				user, err := suite.sut.Sync("some-user-id")
				assert.NoError(t, err)
				assert.Equal(t, "some-user-id", user.Id)
				assert.NotNil(t, user.BrokerParticipantId)
				assert.Equal(t, brokerUserId, *user.BrokerParticipantId)
			},
		)
	})

	t.Run("Given user already exists with a broker participant id", func(t *testing.T) {
		t.Parallel()

		existingBrokerId := uuid.New()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			existing := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			existing.BrokerParticipantId = &existingBrokerId
			suite.mockSlackUserRepository.EXPECT().GetById(gomock.Any()).Return(existing, nil)

			fresh := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name-updated").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserRequester.EXPECT().GetUserInformation(gomock.Any()).Return(fresh, nil)

			// No broker call — participant id already present.
			suite.mockSlackUserRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(user *models.SlackUser) (*models.SlackUser, error) { return user, nil })
		}

		t.Run("Should skip the broker call and just refresh the profile", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			user, err := suite.sut.Sync("some-user-id")
			assert.NoError(t, err)
			assert.Equal(t, existingBrokerId, *user.BrokerParticipantId)
		})
	})

	t.Run("Given error occurs while saving user to the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingUserTestingSuite) {
			suite.mockSlackUserRepository.EXPECT().GetById(gomock.Any()).Return(nil, nil)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId("some-user-id").
				WithName("some-name").
				WithRealName("Some Real Name").
				Build()
			suite.mockSlackUserRequester.EXPECT().GetUserInformation(gomock.Any()).Return(slackUser, nil)

			suite.mockBrokerRequester.EXPECT().
				RegisterUser(gomock.Any()).
				Return(uuid.New(), nil)

			suite.mockSlackUserRepository.EXPECT().Save(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingUserBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync("some-user-id")
			assert.Error(t, err)
		})
	})
}
