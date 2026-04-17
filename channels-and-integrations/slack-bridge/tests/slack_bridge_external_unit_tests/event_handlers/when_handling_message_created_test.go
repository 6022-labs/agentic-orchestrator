package event_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/command_handlers_mocks"
	external_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenHandlingMessageCreatedTestingSuite struct {
	sut *event_handlers.MessageCreatedHandler

	slackSettings *settings.SlackSettings

	mockAddAgentToChannel              *command_handlers_mocks.MockAddAgentToChannelInterface
	mockSlackThreadRequester           *interfaces_mocks.MockSlackThreadRequesterInterface
	mockSlackAgentRepository           *repositories_mocks.MockSlackAgentRepositoryInterface
	mockBrokerMessageForwarder         *external_interfaces_mocks.MockBrokerMessageForwarderInterface
	mockSlackUserSynchronizerCached    *external_interfaces_mocks.MockSlackUserSynchronizerCachedInterface
	mockSlackChannelSynchronizerCached *external_interfaces_mocks.MockSlackChannelSynchronizerCachedInterface
}

func WhenHandlingMessageCreatedBeforeEach(t *testing.T) *WhenHandlingMessageCreatedTestingSuite {
	mockController := gomock.NewController(t)

	slackSettings := &settings.SlackSettings{AppId: "test_app_id"}

	mockAddAgentToChannel := command_handlers_mocks.NewMockAddAgentToChannelInterface(mockController)
	mockSlackThreadRequester := interfaces_mocks.NewMockSlackThreadRequesterInterface(mockController)
	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)
	mockBrokerMessageForwarder := external_interfaces_mocks.NewMockBrokerMessageForwarderInterface(mockController)
	mockSlackUserSynchronizerCached := external_interfaces_mocks.NewMockSlackUserSynchronizerCachedInterface(mockController)
	mockSlackChannelSynchronizerCached := external_interfaces_mocks.NewMockSlackChannelSynchronizerCachedInterface(mockController)

	sut := event_handlers.NewMessageCreatedHandler(
		zap.NewNop(),
		slackSettings,
		mockAddAgentToChannel,
		mockSlackThreadRequester,
		mockSlackAgentRepository,
		mockBrokerMessageForwarder,
		mockSlackUserSynchronizerCached,
		mockSlackChannelSynchronizerCached,
	)

	return &WhenHandlingMessageCreatedTestingSuite{
		sut:                                sut,
		slackSettings:                      slackSettings,
		mockAddAgentToChannel:              mockAddAgentToChannel,
		mockSlackThreadRequester:           mockSlackThreadRequester,
		mockSlackAgentRepository:           mockSlackAgentRepository,
		mockBrokerMessageForwarder:         mockBrokerMessageForwarder,
		mockSlackUserSynchronizerCached:    mockSlackUserSynchronizerCached,
		mockSlackChannelSynchronizerCached: mockSlackChannelSynchronizerCached,
	}
}

func TestWhenHandlingMessageCreated(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while syncing slack channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
			suite.mockSlackChannelSynchronizerCached.EXPECT().
				Sync(request.Event.Channel).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageCreatedBeforeEach(t)
			request := &requests.MessageCreatedRequest{
				Event: requests.MessageCreatedEvent{
					Ts:      "1234567890.123456",
					Channel: "C1234567890",
				},
			}
			initSuite(suite, request)

			err := suite.sut.Handle(request)
			assert.Error(t, err)
		})
	})

	t.Run("Given channel has no facilitator", func(t *testing.T) {
		t.Parallel()

		t.Run("And error message is not already sent", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{
						*model_builders.NewSlackAgentBuilder().WithRole(enums.AgentRoleExpert).Build(),
					}).Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
				suite.mockSlackThreadRequester.EXPECT().
					PushCompletionToThread(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			}

			t.Run("Should push the no-facilitator error to the thread", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)
				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hello",
					},
				}
				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.NoError(t, err)
			})
		})

		t.Run("And error message is already sent", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{
						*model_builders.NewSlackAgentBuilder().WithRole(enums.AgentRoleExpert).Build(),
					}).Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
			}

			t.Run("Should skip without posting again", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)
				request := &requests.MessageCreatedRequest{
					ApiAppId: suite.slackSettings.AppId,
					Event: requests.MessageCreatedEvent{
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    event_handlers.NO_FACILITATOR_FOUND_ERROR,
					},
				}
				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.NoError(t, err)
			})
		})
	})

	t.Run("Given message sent from an agent", func(t *testing.T) {
		t.Parallel()

		collectionAddr := "0x0000000000000000000000000000000000000001"
		collectionTokenId := uint64(7)

		t.Run("And agent is found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().WithRole(enums.AgentRoleFacilitator).Build()
				slackChannel := model_builders.NewSlackChannelBuilder().WithAgents([]models.SlackAgent{facilitator}).Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)

				agent := model_builders.NewSlackAgentBuilder().
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(agent, nil)

				suite.mockBrokerMessageForwarder.EXPECT().
					Forward(slackChannel, request.Event.Ts, request.Event.Ts, request.Event.Text, nil, agent).
					Return(nil)
			}

			t.Run("Should forward the agent message to the broker", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)
				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &suite.slackSettings.AppId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Agent message",
						Metadata: &requests.MessageCreatedMetadata{
							EventType: event_handlers.AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE,
							EventPayload: map[string]any{
								"collection_token_address": collectionAddr,
								"collection_token_id":      float64(collectionTokenId),
							},
						},
					},
				}
				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.NoError(t, err)
			})
		})

		t.Run("And metadata is missing", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().WithRole(enums.AgentRoleFacilitator).Build()
				slackChannel := model_builders.NewSlackChannelBuilder().WithAgents([]models.SlackAgent{facilitator}).Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
			}

			t.Run("Should return an error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)
				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &suite.slackSettings.AppId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
					},
				}
				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})
	})

	t.Run("Given message sent from a slack user", func(t *testing.T) {
		t.Parallel()

		t.Run("And user is found and not paired", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().WithRole(enums.AgentRoleFacilitator).Build()
				slackChannel := model_builders.NewSlackChannelBuilder().WithAgents([]models.SlackAgent{facilitator}).Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(nil).
					WithPairedHumanAgentCollectionTokenId(nil).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().Sync(*request.Event.User).Return(slackUser, nil)

				suite.mockBrokerMessageForwarder.EXPECT().
					Forward(slackChannel, request.Event.Ts, request.Event.Ts, request.Event.Text, slackUser, nil).
					Return(nil)
			}

			t.Run("Should forward the user message to the broker", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)
				userId := "U123"
				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						User:    &userId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hello",
					},
				}
				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.NoError(t, err)
			})
		})

		t.Run("And user sync fails", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().WithRole(enums.AgentRoleFacilitator).Build()
				slackChannel := model_builders.NewSlackChannelBuilder().WithAgents([]models.SlackAgent{facilitator}).Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
				suite.mockSlackUserSynchronizerCached.EXPECT().Sync(*request.Event.User).Return(nil, assert.AnError)
			}

			t.Run("Should return the error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)
				userId := "U123"
				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						User:    &userId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hello",
					},
				}
				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})
	})

	t.Run("Given message has no user and no app_id", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
			facilitator := *model_builders.NewSlackAgentBuilder().WithRole(enums.AgentRoleFacilitator).Build()
			slackChannel := model_builders.NewSlackChannelBuilder().WithAgents([]models.SlackAgent{facilitator}).Build()
			suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingMessageCreatedBeforeEach(t)
			request := &requests.MessageCreatedRequest{
				Event: requests.MessageCreatedEvent{
					Ts:      "1234567890.123456",
					Channel: "C1234567890",
					Text:    "???",
				},
			}
			initSuite(suite, request)

			err := suite.sut.Handle(request)
			assert.Error(t, err)
		})
	})
}
