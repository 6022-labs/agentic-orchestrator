package command_handlers_test

import (
	"context"
	"strings"
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	channel_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	channel_traces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/traces_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenPairingSlackUserToHumanAgentTestingSuite struct {
	sut *command_handlers.PairSlackUserToHumanAgent

	channelMetrics *channel_metrics_mocks.MockChannelMetricsInterface

	slackUserRepository         repositories_mocks.MockSlackUserRepositoryInterface
	slackAgentSynchronizer      interfaces_mocks.MockSlackAgentSynchronizerInterface
	slackUserSynchronizerCached interfaces_mocks.MockSlackUserSynchronizerCachedInterface
}

func WhenPairingSlackUserToHumanAgentBeforeEach(t *testing.T) *WhenPairingSlackUserToHumanAgentTestingSuite {
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

	mockSlackUserRepository := repositories_mocks.NewMockSlackUserRepositoryInterface(mockController)
	mockSlackAgentSynchronizer := interfaces_mocks.NewMockSlackAgentSynchronizerInterface(mockController)
	mockSlackUserSynchronizerCached := interfaces_mocks.NewMockSlackUserSynchronizerCachedInterface(mockController)

	sut := command_handlers.NewPairSlackUserToHumanAgent(
		mockChannelMetrics,
		mockChannelTracer,
		zap.NewNop(),
		mockSlackUserRepository,
		mockSlackAgentSynchronizer,
		mockSlackUserSynchronizerCached,
	)

	return &WhenPairingSlackUserToHumanAgentTestingSuite{
		sut:            sut,
		channelMetrics: mockChannelMetrics,

		slackUserRepository:         *mockSlackUserRepository,
		slackAgentSynchronizer:      *mockSlackAgentSynchronizer,
		slackUserSynchronizerCached: *mockSlackUserSynchronizerCached,
	}
}

func generateValidPairSlackUserToHumanAgentRequest() *requests.PairSlackUserToHumanAgentRequest {
	return &requests.PairSlackUserToHumanAgentRequest{
		RequesterSlackUserId:        "U123456",
		SlackUserId:                 "U654321",
		HumanAgentCollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
		HumanAgentCollectionTokenId: uint64(1),
	}
}

func TestWhenPairingSlackUserToHumanAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given the request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a validation error and record a rejected outcome", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeRejected,
					gomock.Any(),
				).
				Times(1)

			err := suite.sut.Pair(context.Background(), &requests.PairSlackUserToHumanAgentRequest{
				RequesterSlackUserId: "U123456",
			})

			assert.ErrorIs(t, err, common_errors.ErrValidation)
		})
	})

	t.Run("Given an error occurs while syncing the requester slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failing", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to sync requester user")
		})
	})

	t.Run("Given the requester user doesn't have a verified wallet", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(nil).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)
		}

		t.Run("Should return an error indicating the requester needs to verify their wallet", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "requester needs to verify their wallet first before pairing a human agent")
		})
	})

	t.Run("Given an error occurs while syncing the human agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)

			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failure to sync human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to sync human agent")
		})
	})

	t.Run("Given the specified agent is not a human agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleExpert).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)
		}

		t.Run("Should return an error indicating the agent is not a human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "is not a human agent")
		})
	})

	t.Run("Given the requester does not own the human agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(common_tests.RandomEthAddress().Hex()).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)
		}

		t.Run("Should return an error indicating the requester does not own the human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "requester does not own the human agent")
		})
	})

	t.Run("Given an error occurs while getting existing already paired user for this agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failing to get existing pairing for agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to check if human agent is already paired")
		})
	})

	t.Run("Given the specified human agent is already paired to another slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			alreadyExistingSlackUserForAgent := model_builders.NewSlackUserBuilder().
				WithId("U999999").
				WithPairedHumanAgentCollectionAddress(&request.HumanAgentCollectionAddress).
				WithPairedHumanAgentCollectionTokenId(&request.HumanAgentCollectionTokenId).
				Build()
			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(alreadyExistingSlackUserForAgent, nil)
		}

		t.Run(
			"Should return an error indicating the human agent is already paired to another slack user",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

				request := generateValidPairSlackUserToHumanAgentRequest()
				initSuite(suite, request)

				err := suite.sut.Pair(context.Background(), request)

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "is already paired to another slack user")
			},
		)
	})

	t.Run("Given an error occurs while getting existing pairing for slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeFailed,
					gomock.Any(),
				).
				Times(1)

			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, nil)

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.SlackUserId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failing to get existing pairing for slack user", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to sync slack user to pair")
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.channelMetrics.EXPECT().
				RecordCommand(
					gomock.Any(),
					"user.pair",
					metrics_interfaces.CommandTriggerUser,
					metrics_interfaces.ChannelOutcomeHandled,
					gomock.Any(),
				).
				Times(1)

			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(gomock.Any(), request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, nil)

			slackUserToPair := model_builders.NewSlackUserBuilder().
				WithId(request.SlackUserId).
				WithWalletAddress(nil).
				WithPairedHumanAgentCollectionAddress(nil).
				WithPairedHumanAgentCollectionTokenId(nil).
				Build()
			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(gomock.Any(), request.SlackUserId).
				Return(slackUserToPair, nil)

			suite.slackUserRepository.
				EXPECT().
				Save(gomock.Any(), slackUserToPair).
				Return(nil, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(context.Background(), request)

			assert.NoError(t, err)
		})
	})
}
