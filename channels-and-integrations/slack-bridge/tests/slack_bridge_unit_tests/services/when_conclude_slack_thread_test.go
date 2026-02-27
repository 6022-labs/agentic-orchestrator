package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenConcludeSlackThreadTestingSuite struct {
	sut *services.SlackThreadConclude

	mockSlackThreadRepository *repositories_mocks.MockSlackThreadRepositoryInterface
}

func WhenConcludeSlackThreadBeforeEach(t *testing.T) *WhenConcludeSlackThreadTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackThreadRepository := repositories_mocks.NewMockSlackThreadRepositoryInterface(mockController)

	sut := services.NewSlackThreadConclude(
		zap.NewNop(),
		mockSlackThreadRepository,
	)

	return &WhenConcludeSlackThreadTestingSuite{
		sut:                       sut,
		mockSlackThreadRepository: mockSlackThreadRepository,
	}
}

func TestWhenConcludeSlackThread(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while fetching slack thread from database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConcludeSlackThreadTestingSuite) {
			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId("threadId", "channelId").
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludeSlackThreadBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Conclude("threadId", "channelId")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to fetch slack thread from database")
		})
	})

	t.Run("Given slack thread is not found in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConcludeSlackThreadTestingSuite) {
			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId("threadId", "channelId").
				Return(nil, nil)
		}

		t.Run("Should return an error indicating thread not found", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludeSlackThreadBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Conclude("threadId", "channelId")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "slack thread not found in database")
		})
	})

	t.Run("Given error occurs while saving slack thread state", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConcludeSlackThreadTestingSuite) {
			slackThread := model_builders.NewSlackThreadBuilder().Build()
			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId("threadId", "channelId").
				Return(slackThread, nil)

			suite.mockSlackThreadRepository.
				EXPECT().
				Save(slackThread).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludeSlackThreadBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Conclude("threadId", "channelId")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to save slack thread state")
		})
	})

	t.Run("Given thread is already concluded", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConcludeSlackThreadTestingSuite) {
			slackThread := model_builders.NewSlackThreadBuilder().
				WithIsConcluded(true).
				Build()

			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId("threadId", "channelId").
				Return(slackThread, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludeSlackThreadBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Conclude("threadId", "channelId")

			assert.NoError(t, err)
		})
	})

	t.Run("Given no error occurs", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenConcludeSlackThreadTestingSuite) {
			slackThread := model_builders.NewSlackThreadBuilder().
				WithIsConcluded(false).
				Build()
			suite.mockSlackThreadRepository.
				EXPECT().
				GetByIdAndChannelId("threadId", "channelId").
				Return(slackThread, nil)

			suite.mockSlackThreadRepository.
				EXPECT().
				Save(slackThread).
				Return(slackThread, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludeSlackThreadBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Conclude("threadId", "channelId")

			assert.NoError(t, err)
		})
	})
}
