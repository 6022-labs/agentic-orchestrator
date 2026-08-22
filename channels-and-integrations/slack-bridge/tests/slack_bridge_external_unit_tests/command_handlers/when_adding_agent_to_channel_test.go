package command_handlers_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenAddingAgentToChannelTestingSuite struct {
	sut *command_handlers.AddAgentToChannel

	mockSlackAgentSynchronizer      *interfaces_mocks.MockSlackAgentSynchronizerInterface
	mockSlackChannelSynchronizer    *interfaces_mocks.MockSlackChannelSynchronizerInterface
	mockSlackChannelRepository      *repositories_mocks.MockSlackChannelRepositoryInterface
	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester             *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenAddingAgentToChannelBeforeEach(t *testing.T) *WhenAddingAgentToChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackAgentSynchronizer := interfaces_mocks.NewMockSlackAgentSynchronizerInterface(mockController)
	mockSlackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)
	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := command_handlers.NewAddAgentToChannel(
		mockSlackAgentSynchronizer,
		mockSlackChannelRepository,
		mockSlackChannelSynchronizer,
		mockSlackConversationRepository,
		mockBrokerRequester,
	)

	return &WhenAddingAgentToChannelTestingSuite{
		sut: sut,

		mockSlackAgentSynchronizer:      mockSlackAgentSynchronizer,
		mockSlackChannelRepository:      mockSlackChannelRepository,
		mockSlackChannelSynchronizer:    mockSlackChannelSynchronizer,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:             mockBrokerRequester,
	}
}

func TestWhenAddingAgentToChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting the channel from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the channel already contains the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*model_builders.NewSlackAgentBuilder().
						WithCollectionAddress(request.CollectionAddress).
						WithCollectionTokenId(request.CollectionTokenId).
						Build(),
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)
		}

		t.Run("Should not attempt to add the agent again and return no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an error occurs while syncing the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while syncing agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return sync agent error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while adding agent to channel in database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			brokerParticipantId := uuid.New()
			slackAgent := model_builders.NewSlackAgentBuilder().Build()
			slackAgent.BrokerParticipantId = &brokerParticipantId
			suite.mockSlackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				AddAgent(gomock.Any(), slackChannel, slackAgent).
				Return(nil, assert.AnError)
		}

		t.Run("Should return save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the channel is not yet in the database", func(t *testing.T) {
		t.Parallel()

		brokerAgentId := uuid.New()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), request.ChannelID).Return(nil, nil)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelSynchronizer.EXPECT().Sync(gomock.Any(), request.ChannelID).Return(slackChannel, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackAgent.BrokerParticipantId = &brokerAgentId
			suite.mockSlackAgentSynchronizer.EXPECT().
				Sync(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.EXPECT().
				AddAgent(gomock.Any(), slackChannel, slackAgent).
				Return(slackChannel, nil)

			suite.mockSlackConversationRepository.EXPECT().
				ListByChannelId(gomock.Any(), request.ChannelID).
				Return([]*models.SlackConversation{}, nil)
		}

		t.Run("Should sync the channel and add the agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given all steps succeed and the channel has existing broker conversations", func(t *testing.T) {
		t.Parallel()

		brokerAgentId := uuid.New()
		convId1 := uuid.New()
		convId2 := uuid.New()

		initSuite := func(suite *WhenAddingAgentToChannelTestingSuite, request *requests.AddAgentToChannelRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), request.ChannelID).Return(slackChannel, nil)
			suite.mockSlackChannelSynchronizer.EXPECT().Sync(gomock.Any(), request.ChannelID).Return(slackChannel, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackAgent.BrokerParticipantId = &brokerAgentId
			suite.mockSlackAgentSynchronizer.EXPECT().
				Sync(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.EXPECT().
				AddAgent(gomock.Any(), slackChannel, slackAgent).
				Return(slackChannel, nil)

			suite.mockSlackConversationRepository.EXPECT().
				ListByChannelId(gomock.Any(), request.ChannelID).
				Return([]*models.SlackConversation{
					{SlackChannelId: request.ChannelID, SlackThreadTs: "1.0", BrokerConversationID: convId1},
					{SlackChannelId: request.ChannelID, SlackThreadTs: "2.0", BrokerConversationID: convId2},
				}, nil)

			suite.mockBrokerRequester.EXPECT().AttachParticipant(gomock.Any(), convId1, brokerAgentId).Return(nil)
			suite.mockBrokerRequester.EXPECT().AttachParticipant(gomock.Any(), convId2, brokerAgentId).Return(nil)
		}

		t.Run("Should attach the agent to every existing broker conversation in the channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenAddingAgentToChannelBeforeEach(t)

			request := &requests.AddAgentToChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Add(context.Background(), request)

			assert.NoError(t, err)
		})
	})
}
