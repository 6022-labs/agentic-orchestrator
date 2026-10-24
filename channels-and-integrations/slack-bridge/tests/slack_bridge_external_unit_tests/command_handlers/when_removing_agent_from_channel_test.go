package command_handlers_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	channel_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	channel_traces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/traces_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRemovingAgentFromChannelTestingSuite struct {
	sut *command_handlers.RemoveAgentFromChannel

	mockChannelMetrics *channel_metrics_mocks.MockChannelMetricsInterface

	mockSlackAgentRepository        *repositories_mocks.MockSlackAgentRepositoryInterface
	mockSlackChannelSynchronizer    *interfaces_mocks.MockSlackChannelSynchronizerInterface
	mockSlackChannelRepository      *repositories_mocks.MockSlackChannelRepositoryInterface
	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester             *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenRemovingAgentFromChannelBeforeEach(t *testing.T) *WhenRemovingAgentFromChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockChannelMetrics := channel_metrics_mocks.NewMockChannelMetricsInterface(mockController)

	// Spans are not asserted here.
	mockChannelTracer := channel_traces_mocks.NewMockChannelTracerInterface(mockController)
	mockSpan := channel_traces_mocks.NewMockSpan(mockController)
	mockSpan.EXPECT().End().AnyTimes()
	mockSpan.EXPECT().Fail(gomock.Any()).AnyTimes()
	mockChannelTracer.EXPECT().
		StartCommandHandling(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ string) (context.Context, traces_interfaces.Span) {
			return ctx, mockSpan
		}).
		AnyTimes()

	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)
	mockSlackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)
	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := command_handlers.NewRemoveAgentFromChannel(
		mockChannelMetrics,
		mockChannelTracer,
		zap.NewNop(),
		mockSlackAgentRepository,
		mockSlackChannelRepository,
		mockSlackChannelSynchronizer,
		mockSlackConversationRepository,
		mockBrokerRequester,
	)

	return &WhenRemovingAgentFromChannelTestingSuite{
		sut:                             sut,
		mockChannelMetrics:              mockChannelMetrics,
		mockSlackAgentRepository:        mockSlackAgentRepository,
		mockSlackChannelSynchronizer:    mockSlackChannelSynchronizer,
		mockSlackChannelRepository:      mockSlackChannelRepository,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:             mockBrokerRequester,
	}
}

func TestWhenRemovingAgentFromChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given the request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and record a rejected outcome", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeRejected,
					gomock.Any(),
				).
				Times(1)

			err := suite.sut.Remove(context.Background(), &requests.RemoveAgentFromChannelRequest{
				ChannelID:             "chan",
				AgentCommandArguments: requests.AgentCommandArguments{},
			})

			var validationError *common_errors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "collectionAddress", validationError.Field)
		})
	})

	t.Run("Given error occurs while getting the channel from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the channel did not contain the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeHandled,
					gomock.Any(),
				).
				Times(1)

			slackChannel := model_builders.NewSlackChannelBuilder().
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)
		}

		t.Run("Should not attempt to remove the agent again and return no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an error occurs while syncing the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(nil, errors.New("sync error"))
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent from repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(nil, errors.New("repo error"))
		}

		t.Run("Should return repo error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent not found in repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeHandled,
					gomock.Any(),
				).
				Times(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the channel is not in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeHandled,
					gomock.Any(),
				).
				Times(1)

			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(nil, nil)

			slackChannel := model_builders.NewSlackChannelBuilder().Build()
			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should proceed with the sync and not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given error occurs while saving channel after removal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{
					*slackAgent,
				}).
				Build()
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackChannelSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.ChannelID).
				Return(slackChannel, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)

			suite.mockSlackChannelRepository.
				EXPECT().
				RemoveAgent(gomock.Any(), slackChannel, slackAgent).
				Return(nil, errors.New("save error"))
		}

		t.Run("Should return save error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given all steps succeed and the channel has existing broker conversations", func(t *testing.T) {
		t.Parallel()

		brokerAgentId := uuid.New()
		convId1 := uuid.New()
		convId2 := uuid.New()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockChannelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"agent.remove",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeHandled,
					gomock.Any(),
				).
				Times(1)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(request.CollectionAddress).
				WithCollectionTokenId(request.CollectionTokenId).
				Build()
			slackAgent.BrokerParticipantId = &brokerAgentId
			slackChannel := model_builders.NewSlackChannelBuilder().
				WithAgents([]models.SlackAgent{*slackAgent}).
				Build()
			suite.mockSlackChannelRepository.EXPECT().GetById(gomock.Any(), request.ChannelID).Return(slackChannel, nil)
			suite.mockSlackChannelSynchronizer.EXPECT().Sync(gomock.Any(), request.ChannelID).Return(slackChannel, nil)
			suite.mockSlackAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), request.CollectionAddress, request.CollectionTokenId).
				Return(slackAgent, nil)
			suite.mockSlackChannelRepository.EXPECT().
				RemoveAgent(gomock.Any(), slackChannel, slackAgent).
				Return(slackChannel, nil)

			suite.mockSlackConversationRepository.EXPECT().
				ListByChannelId(gomock.Any(), request.ChannelID).
				Return([]*models.SlackConversation{
					{SlackChannelId: request.ChannelID, SlackThreadTs: "1.0", BrokerConversationID: convId1},
					{SlackChannelId: request.ChannelID, SlackThreadTs: "2.0", BrokerConversationID: convId2},
				}, nil)

			suite.mockBrokerRequester.EXPECT().DetachParticipant(gomock.Any(), convId1, brokerAgentId).Return(nil)
			suite.mockBrokerRequester.EXPECT().DetachParticipant(gomock.Any(), convId2, brokerAgentId).Return(nil)
		}

		t.Run("Should detach the agent from every existing broker conversation in the channel", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID: "chan",
				AgentCommandArguments: requests.AgentCommandArguments{
					CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
					CollectionTokenId: 123,
				},
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)
			assert.NoError(t, err)
		})
	})
}
