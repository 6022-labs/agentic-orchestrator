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
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
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
	mockSlackThreadRepository          *repositories_mocks.MockSlackThreadRepositoryInterface
	mockSlackThreadMessagesAdder       *external_interfaces_mocks.MockSlackThreadMessagesAdderInterface
	mockSlackUserSynchronizerCached    *external_interfaces_mocks.MockSlackUserSynchronizerCachedInterface
	mockSlackChannelSynchronizerCached *external_interfaces_mocks.MockSlackChannelSynchronizerCachedInterface
}

func WhenHandlingMessageCreatedBeforeEach(t *testing.T) *WhenHandlingMessageCreatedTestingSuite {
	mockController := gomock.NewController(t)

	slackSettings := &settings.SlackSettings{
		AppId: "test_app_id",
	}

	mockAddAgentToChannel := command_handlers_mocks.NewMockAddAgentToChannelInterface(mockController)
	mockSlackThreadRequester := interfaces_mocks.NewMockSlackThreadRequesterInterface(mockController)
	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)
	mockSlackThreadRepository := repositories_mocks.NewMockSlackThreadRepositoryInterface(mockController)
	mockSlackThreadMessagesAdder := external_interfaces_mocks.NewMockSlackThreadMessagesAdderInterface(mockController)
	mockSlackUserSynchronizerCached := external_interfaces_mocks.NewMockSlackUserSynchronizerCachedInterface(mockController)
	mockSlackChannelSynchronizerCached := external_interfaces_mocks.NewMockSlackChannelSynchronizerCachedInterface(mockController)

	sut := event_handlers.NewMessageCreatedHandler(
		zap.NewNop(),
		slackSettings,
		mockAddAgentToChannel,
		mockSlackThreadRequester,
		mockSlackAgentRepository,
		mockSlackThreadRepository,
		mockSlackThreadMessagesAdder,
		mockSlackUserSynchronizerCached,
		mockSlackChannelSynchronizerCached,
	)

	return &WhenHandlingMessageCreatedTestingSuite{
		sut: sut,

		slackSettings: slackSettings,

		mockAddAgentToChannel:              mockAddAgentToChannel,
		mockSlackThreadRequester:           mockSlackThreadRequester,
		mockSlackAgentRepository:           mockSlackAgentRepository,
		mockSlackThreadRepository:          mockSlackThreadRepository,
		mockSlackThreadMessagesAdder:       mockSlackThreadMessagesAdder,
		mockSlackUserSynchronizerCached:    mockSlackUserSynchronizerCached,
		mockSlackChannelSynchronizerCached: mockSlackChannelSynchronizerCached,
	}
}

func TestWhenHandlingMessageCreated(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while syncing slack channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
			suite.mockSlackChannelSynchronizerCached.
				EXPECT().
				Sync(request.Event.Channel).
				Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
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

	t.Run("Given channel has no facilitator", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{*model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleExpert).
					Build(),
				}).
				Build()
			suite.mockSlackChannelSynchronizerCached.
				EXPECT().
				Sync(request.Event.Channel).
				Return(slackChannel, nil)
		}

		t.Run("And error message is already sent", func(t *testing.T) {
			t.Parallel()

			t.Run("Should not call PushCompletionToThread", func(t *testing.T) {
				t.Parallel()
				suite := WhenHandlingMessageCreatedBeforeEach(t)

				request := &requests.MessageCreatedRequest{
					ApiAppId: suite.slackSettings.AppId,
					Event: requests.MessageCreatedEvent{
						AppId:   &suite.slackSettings.AppId,
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

		t.Run("And error message is not already sent", func(t *testing.T) {
			t.Parallel()

			t.Run("Should call PushCompletionToThread", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)

				request := &requests.MessageCreatedRequest{
					ApiAppId: suite.slackSettings.AppId,
					Event: requests.MessageCreatedEvent{
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Some other message",
					},
				}

				initSuite(suite, request)

				suite.mockSlackThreadRequester.
					EXPECT().
					PushCompletionToThread(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := suite.sut.Handle(request)
				assert.NoError(t, err)
			})
		})
	})

	t.Run("Given message sent from an agent", func(t *testing.T) {
		t.Parallel()

		t.Run("And message does not contain metadata", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And unsupported metadata event type", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()

				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)

				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &suite.slackSettings.AppId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hi",
						Metadata: &requests.MessageCreatedMetadata{
							EventType:    "some_other_event",
							EventPayload: map[string]any{},
						},
					},
				}

				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})

		t.Run("And invalid metadata payload (collection_token_id type)", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)

				appId := suite.slackSettings.AppId

				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &appId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hi",
						Metadata: &requests.MessageCreatedMetadata{
							EventType: event_handlers.AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE,
							EventPayload: map[string]any{
								"collection_token_address": "0x0000000000000000000000000000000000000001",
								// Wrong type on purpose: should be uint64
								"collection_token_id": "1",
							},
						},
					},
				}

				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})

		t.Run("And invalid metadata payload (collection_token_id negative)", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)

				appId := suite.slackSettings.AppId

				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &appId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hi",
						Metadata: &requests.MessageCreatedMetadata{
							EventType: event_handlers.AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE,
							EventPayload: map[string]any{
								"collection_token_address": "0x0000000000000000000000000000000000000001",
								// Wrong type on purpose: should be uint64
								"collection_token_id": float64(-1),
							},
						},
					},
				}

				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})

		t.Run("And invalid metadata payload (collection_token_id non-integer)", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)

				appId := suite.slackSettings.AppId

				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &appId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hi",
						Metadata: &requests.MessageCreatedMetadata{
							EventType: event_handlers.AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE,
							EventPayload: map[string]any{
								"collection_token_address": "0x0000000000000000000000000000000000000001",
								// Wrong type on purpose: should be uint64
								"collection_token_id": float64(5.5),
							},
						},
					},
				}

				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})

		t.Run("And invalid metadata payload (collection_token_address format)", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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
								// Invalid hex address
								"collection_token_address": "not_an_address",
								"collection_token_id":      float64(1),
							},
						},
					},
				}

				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while fetching the agent", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0x0000000000000000000000000000000000000001"
				collectionTokenId := uint64(1)
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)
				appId := suite.slackSettings.AppId

				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &appId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Agent message",
						Metadata: &requests.MessageCreatedMetadata{
							EventType: event_handlers.AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE,
							EventPayload: map[string]any{
								"collection_token_address": "0x0000000000000000000000000000000000000001",
								"collection_token_id":      float64(1),
							},
						},
					},
				}

				initSuite(suite, request)

				err := suite.sut.Handle(request)
				assert.Error(t, err)
			})
		})

		t.Run("And agent is not found", func(t *testing.T) {
			t.Parallel()

			collectionAddr := "0x0000000000000000000000000000000000000001"
			collectionTokenId := uint64(42)

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().Sync(request.Event.Channel).Return(slackChannel, nil)

				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(nil, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)

				appId := suite.slackSettings.AppId

				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &appId,
						Ts:      "1234567890.123456",
						Channel: "C1234567890",
						Text:    "Hi",
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
				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while adding message to the thread", func(t *testing.T) {
			t.Parallel()

			collectionAddr := "0x0000000000000000000000000000000000000001"
			collectionTokenId := uint64(7)

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				agent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleExpert).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(agent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						nil,
						&agent.CollectionAddress,
						&agent.CollectionTokenId,
					).
					Return(assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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
				assert.Error(t, err)
			})
		})

		t.Run("And message is successfully stored", func(t *testing.T) {
			t.Parallel()

			collectionAddr := "0x0000000000000000000000000000000000000001"
			collectionTokenId := uint64(9)

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync("C1234567890").
					Return(slackChannel, nil)

				agent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleExpert).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(agent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						nil,
						&agent.CollectionAddress,
						&agent.CollectionTokenId,
					).
					Return(nil)
			}

			t.Run("Should not return a error", func(t *testing.T) {
				t.Parallel()

				suite := WhenHandlingMessageCreatedBeforeEach(t)

				appId := suite.slackSettings.AppId

				request := &requests.MessageCreatedRequest{
					Event: requests.MessageCreatedEvent{
						AppId:   &appId,
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
	})

	t.Run("Given message sent from a user", func(t *testing.T) {
		t.Parallel()

		t.Run("And message does not contains the user ID", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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
				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while fetching the user", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And user is not found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(nil, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And error occurs while adding message to the thread", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(nil).
					WithPairedHumanAgentCollectionTokenId(nil).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(slackUser, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						request.Event.User,
						nil,
						nil,
					).
					Return(assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And error occurs while getting slack thread", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(nil).
					WithPairedHumanAgentCollectionTokenId(nil).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(slackUser, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						request.Event.User,
						nil,
						nil,
					).
					Return(nil)

				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And slack thread is not found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(nil).
					WithPairedHumanAgentCollectionTokenId(nil).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(slackUser, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						request.Event.User,
						nil,
						nil,
					).
					Return(nil)

				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(nil, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And thread is concluded but error occurs while unconcluding it", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(nil).
					WithPairedHumanAgentCollectionTokenId(nil).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(slackUser, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						request.Event.User,
						nil,
						nil,
					).
					Return(nil)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(request.Event.Ts).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(true).
					Build()
				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(slackThread, nil)

				suite.mockSlackThreadRepository.EXPECT().
					Save(slackThread).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And message is successfully stored", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(nil).
					WithPairedHumanAgentCollectionTokenId(nil).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(slackUser, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						request.Event.User,
						nil,
						nil,
					).
					Return(nil)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(request.Event.Ts).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(true).
					Build()
				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(slackThread, nil)

				suite.mockSlackThreadRepository.EXPECT().
					Save(slackThread).
					Return(nil, nil)
			}

			t.Run("Should not return a error", func(t *testing.T) {
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
	})

	t.Run("Given message sent from a user that is paired to a human agent", func(t *testing.T) {
		t.Parallel()

		t.Run("And message does not contains the user ID", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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
				assert.Error(t, err)
			})
		})

		t.Run("And error occurs while fetching the user", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And user is not found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(*request.Event.User).
					Return(nil, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And error occurs while fetching the paired human agent", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(nil, assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And paired human agent is not found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()

				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(nil, nil)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And error occurs while adding message to the thread", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				pairedAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(pairedAgent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						&slackUser.Id,
						&pairedAgent.CollectionAddress,
						&pairedAgent.CollectionTokenId,
					).
					Return(assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And a error occurs while adding the human agent to the channel", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				pairedAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(pairedAgent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						&slackUser.Id,
						&pairedAgent.CollectionAddress,
						&pairedAgent.CollectionTokenId,
					).
					Return(nil)

				suite.mockAddAgentToChannel.EXPECT().
					Add(gomock.Any()).
					Return(assert.AnError)
			}

			t.Run("Should return error", func(t *testing.T) {
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

		t.Run("And error occurs while getting slack thread", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				pairedAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(pairedAgent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						&slackUser.Id,
						&pairedAgent.CollectionAddress,
						&pairedAgent.CollectionTokenId,
					).
					Return(nil)

				suite.mockAddAgentToChannel.EXPECT().
					Add(gomock.Any()).
					Return(nil)

				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(nil, assert.AnError)
			}

			t.Run("Should return a error", func(t *testing.T) {
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

		t.Run("And slack thread is not found", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				pairedAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(pairedAgent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						&slackUser.Id,
						&pairedAgent.CollectionAddress,
						&pairedAgent.CollectionTokenId,
					).
					Return(nil)

				suite.mockAddAgentToChannel.EXPECT().
					Add(gomock.Any()).
					Return(nil)

				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(nil, nil)
			}

			t.Run("Should return a error", func(t *testing.T) {
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

		t.Run("And slack thread is concluded but error occurs while unconcluding it", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				pairedAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(pairedAgent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						&slackUser.Id,
						&pairedAgent.CollectionAddress,
						&pairedAgent.CollectionTokenId,
					).
					Return(nil)

				suite.mockAddAgentToChannel.EXPECT().
					Add(gomock.Any()).
					Return(nil)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(request.Event.Ts).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(true).
					Build()
				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(slackThread, nil)

				suite.mockSlackThreadRepository.EXPECT().
					Save(slackThread).
					Return(nil, assert.AnError)
			}

			t.Run("Should return a error", func(t *testing.T) {
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

		t.Run("And message is successfully stored", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenHandlingMessageCreatedTestingSuite, request *requests.MessageCreatedRequest) {
				facilitator := *model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleFacilitator).
					Build()
				slackChannel := model_builders.NewSlackChannelBuilder().
					WithAgents([]models.SlackAgent{facilitator}).
					Build()
				suite.mockSlackChannelSynchronizerCached.EXPECT().
					Sync(request.Event.Channel).
					Return(slackChannel, nil)

				collectionAddr := "0xABC"
				collectionTokenId := uint64(1)

				slackUser := model_builders.NewSlackUserBuilder().
					WithId(*request.Event.User).
					WithPairedHumanAgentCollectionAddress(&collectionAddr).
					WithPairedHumanAgentCollectionTokenId(&collectionTokenId).
					Build()
				suite.mockSlackUserSynchronizerCached.EXPECT().
					Sync(slackUser.Id).
					Return(slackUser, nil)

				pairedAgent := model_builders.NewSlackAgentBuilder().
					WithRole(enums.AgentRoleHuman).
					WithCollectionAddress(collectionAddr).
					WithCollectionTokenId(collectionTokenId).
					Build()
				suite.mockSlackAgentRepository.EXPECT().
					GetByCollectionAddressAndCollectionTokenId(collectionAddr, collectionTokenId).
					Return(pairedAgent, nil)

				suite.mockSlackThreadMessagesAdder.EXPECT().
					Add(
						slackChannel.Id,
						request.Event.Ts,
						request.Event.Ts,
						request.Event.Text,
						&slackUser.Id,
						&pairedAgent.CollectionAddress,
						&pairedAgent.CollectionTokenId,
					).
					Return(nil)

				suite.mockAddAgentToChannel.EXPECT().
					Add(gomock.Any()).
					Return(nil)

				slackThread := model_builders.NewSlackThreadBuilder().
					WithId(request.Event.Ts).
					WithChannelId(slackChannel.Id).
					WithIsConcluded(true).
					Build()
				suite.mockSlackThreadRepository.EXPECT().
					GetByIdAndChannelId(request.Event.Ts, slackChannel.Id).
					Return(slackThread, nil)

				suite.mockSlackThreadRepository.EXPECT().
					Save(slackThread).
					Return(nil, nil)
			}

			t.Run("Should not return a error", func(t *testing.T) {
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
	})
}
