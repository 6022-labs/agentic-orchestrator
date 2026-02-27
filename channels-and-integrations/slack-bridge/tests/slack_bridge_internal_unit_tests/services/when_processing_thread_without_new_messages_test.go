package services_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	internal_services "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	services_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenProcessingThreadWithoutNewMessagesTestingSuite struct {
	sut *internal_services.SlackThreadProcessor

	processingSettings *settings.ProcessingSettings

	mockSlackThreadConclude       *services_mocks.MockSlackThreadConcludeInterface
	mockSlackThreadRequester      *services_mocks.MockSlackThreadRequesterInterface
	mockSlackThreadRepository     *repositories_mocks.MockSlackThreadRepositoryInterface
	mockSlackMessageRepository    *repositories_mocks.MockSlackMessageRepositoryInterface
	mockSlackChannelRepository    *repositories_mocks.MockSlackChannelRepositoryInterface
	mockAgentCompletionsRequester *services_mocks.MockAgentCompletionsRequesterInterface
}

func WhenProcessingThreadWithoutNewMessagesBeforeEach(t *testing.T) *WhenProcessingThreadWithoutNewMessagesTestingSuite {
	mockController := gomock.NewController(t)

	processingSettings := &settings.ProcessingSettings{
		AgentProcessingInterval:                     20,
		FacilitatorTickInterval:                     30,
		ConversationForceConcludeByFacilitatorAfter: 300,
		ConversationForceConcludeAfter:              60,
	}

	mockSlackThreadConclude := services_mocks.NewMockSlackThreadConcludeInterface(mockController)
	mockSlackThreadRequester := services_mocks.NewMockSlackThreadRequesterInterface(mockController)
	mockSlackThreadRepository := repositories_mocks.NewMockSlackThreadRepositoryInterface(mockController)
	mockSlackMessageRepository := repositories_mocks.NewMockSlackMessageRepositoryInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)
	mockAgentCompletionsRequester := services_mocks.NewMockAgentCompletionsRequesterInterface(mockController)

	sut := internal_services.NewSlackThreadProcessor(
		zap.NewNop(),
		processingSettings,
		mockSlackThreadConclude,
		mockSlackThreadRequester,
		mockSlackThreadRepository,
		mockSlackMessageRepository,
		mockSlackChannelRepository,
		mockAgentCompletionsRequester,
	)

	return &WhenProcessingThreadWithoutNewMessagesTestingSuite{
		sut:                           sut,
		processingSettings:            processingSettings,
		mockSlackThreadConclude:       mockSlackThreadConclude,
		mockSlackThreadRequester:      mockSlackThreadRequester,
		mockSlackThreadRepository:     mockSlackThreadRepository,
		mockSlackMessageRepository:    mockSlackMessageRepository,
		mockSlackChannelRepository:    mockSlackChannelRepository,
		mockAgentCompletionsRequester: mockAgentCompletionsRequester,
	}
}

func TestWhenProcessingThreadWithoutNewMessages(t *testing.T) {
	t.Parallel()

	t.Run("Given safeguard threshold is not reached", func(t *testing.T) {
		t.Parallel()

		suite := WhenProcessingThreadWithoutNewMessagesBeforeEach(t)

		agent := model_builders.NewSlackAgentBuilder().
			WithRole(enums.AgentRoleFacilitator).
			Build()
		channel := model_builders.NewSlackChannelBuilder().Build()

		thread := model_builders.NewSlackThreadBuilder().
			WithChannelId(channel.Id).
			WithId(model_builders.GenerateRandomSlackMessageID()).
			Build()

		suite.mockSlackChannelRepository.
			EXPECT().
			GetById(thread.ChannelId).
			Return(channel, nil)

		suite.mockSlackMessageRepository.
			EXPECT().
			GetProcessableMessagesForAgentAndThread(agent, thread.Id, thread.ChannelId).
			Return([]*models.SlackMessage{}, nil)

		err := suite.sut.ProcessThread(agent, thread)

		assert.NoError(t, err)
	})

	t.Run("Given safeguard threshold is exceeded", func(t *testing.T) {
		t.Parallel()

		suite := WhenProcessingThreadWithoutNewMessagesBeforeEach(t)

		agent := model_builders.NewSlackAgentBuilder().
			WithRole(enums.AgentRoleFacilitator).
			Build()
		channel := model_builders.NewSlackChannelBuilder().Build()

		inactivityDuration := time.Duration(suite.processingSettings.ConversationForceConcludeAfter+30) * time.Second
		past := time.Now().Add(-inactivityDuration)
		thread := model_builders.NewSlackThreadBuilder().
			WithChannelId(channel.Id).
			WithId(model_builders.GenerateSlackMessageID(past)).
			WithUnconcludedAt(&past).
			Build()

		suite.mockSlackThreadConclude.
			EXPECT().
			Conclude(thread.Id, thread.ChannelId).
			Return(nil)

		err := suite.sut.ProcessThread(agent, thread)

		assert.NoError(t, err)
	})

	t.Run("Given safeguard auto conclude fails", func(t *testing.T) {
		t.Parallel()

		suite := WhenProcessingThreadWithoutNewMessagesBeforeEach(t)

		agent := model_builders.NewSlackAgentBuilder().
			WithRole(enums.AgentRoleFacilitator).
			Build()
		channel := model_builders.NewSlackChannelBuilder().Build()

		inactivityDuration := time.Duration(suite.processingSettings.ConversationForceConcludeAfter+30) * time.Second
		past := time.Now().Add(-inactivityDuration)
		thread := model_builders.NewSlackThreadBuilder().
			WithChannelId(channel.Id).
			WithId(model_builders.GenerateSlackMessageID(past)).
			WithUnconcludedAt(&past).
			Build()

		suite.mockSlackThreadConclude.
			EXPECT().
			Conclude(thread.Id, thread.ChannelId).
			Return(assert.AnError)

		err := suite.sut.ProcessThread(agent, thread)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to force conclude slack thread")
	})
}
