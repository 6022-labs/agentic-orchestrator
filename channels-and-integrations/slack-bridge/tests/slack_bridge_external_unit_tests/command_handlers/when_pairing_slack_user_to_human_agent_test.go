package command_handlers_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenPairingSlackUserToHumanAgentTestingSuite struct {
	sut *command_handlers.PairSlackUserToHumanAgent

	slackUserRepository         repositories_mocks.MockSlackUserRepositoryInterface
	slackAgentSynchronizer      interfaces_mocks.MockSlackAgentSynchronizerInterface
	slackUserSynchronizerCached interfaces_mocks.MockSlackUserSynchronizerCachedInterface
}

func WhenPairingSlackUserToHumanAgentBeforeEach(t *testing.T) *WhenPairingSlackUserToHumanAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackUserRepository := repositories_mocks.NewMockSlackUserRepositoryInterface(mockController)
	mockSlackAgentSynchronizer := interfaces_mocks.NewMockSlackAgentSynchronizerInterface(mockController)
	mockSlackUserSynchronizerCached := interfaces_mocks.NewMockSlackUserSynchronizerCachedInterface(mockController)

	sut := command_handlers.NewPairSlackUserToHumanAgent(
		zap.NewNop(),
		mockSlackUserRepository,
		mockSlackAgentSynchronizer,
		mockSlackUserSynchronizerCached,
	)

	return &WhenPairingSlackUserToHumanAgentTestingSuite{
		sut: sut,

		slackUserRepository:         *mockSlackUserRepository,
		slackAgentSynchronizer:      *mockSlackAgentSynchronizer,
		slackUserSynchronizerCached: *mockSlackUserSynchronizerCached,
	}
}

func generateValidPairSlackUserToHumanAgentRequest() *requests.PairSlackUserToHumanAgentRequest {
	return &requests.PairSlackUserToHumanAgentRequest{
		RequesterSlackUserId:        "U123456",
		SlackUserId:                 "U654321",
		HumanAgentCollectionAddress: common_tests.RandomEthAddress().Hex(),
		HumanAgentCollectionTokenId: uint64(1),
	}
}

func TestWhenPairingSlackUserToHumanAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given an error occurs while syncing the requester slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failing", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to sync requester user")
		})
	})

	t.Run("Given the requester user doesn't have a verified wallet", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(nil).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)
		}

		t.Run("Should return an error indicating the requester needs to verify their wallet", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "requester needs to verify their wallet first before pairing a human agent")
		})
	})

	t.Run("Given an error occurs while syncing the human agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)

			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failure to sync human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to sync human agent")
		})
	})

	t.Run("Given the specified agent is not a human agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleExpert).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)
		}

		t.Run("Should return an error indicating the agent is not a human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "is not a human agent")
		})
	})

	t.Run("Given the requester does not own the human agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(common_tests.RandomEthAddress().Hex()).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)
		}

		t.Run("Should return an error indicating the requester does not own the human agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "requester does not own the human agent")
		})
	})

	t.Run("Given an error occurs while getting existing already paired user for this agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failing to get existing pairing for agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to check if human agent is already paired")
		})
	})

	t.Run("Given the specified human agent is already paired to another slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			alreadyExistingSlackUserForAgent := model_builders.NewSlackUserBuilder().
				WithId("U999999").
				WithPairedHumanAgentCollectionAddress(&request.HumanAgentCollectionAddress).
				WithPairedHumanAgentCollectionTokenId(&request.HumanAgentCollectionTokenId).
				Build()
			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(alreadyExistingSlackUserForAgent, nil)
		}

		t.Run("Should return an error indicating the human agent is already paired to another slack user", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "is already paired to another slack user")
		})
	})

	t.Run("Given an error occurs while getting existing pairing for slack user", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, nil)

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.SlackUserId).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error indicating failing to get existing pairing for slack user", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to sync slack user to pair")
		})
	})

	t.Run("Given no errors occur", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPairingSlackUserToHumanAgentTestingSuite, request *requests.PairSlackUserToHumanAgentRequest) {
			slackUserWalletAddress := common_tests.RandomEthAddress().Hex()
			slackUser := model_builders.NewSlackUserBuilder().
				WithId(request.RequesterSlackUserId).
				WithWalletAddress(&slackUserWalletAddress).
				Build()

			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.RequesterSlackUserId).
				Return(slackUser, nil)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithRole(enums.AgentRoleHuman).
				WithOwnerWallet(slackUserWalletAddress).
				WithCollectionAddress(request.HumanAgentCollectionAddress).
				WithCollectionTokenId(request.HumanAgentCollectionTokenId).
				Build()
			suite.slackAgentSynchronizer.
				EXPECT().
				Sync(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(slackAgent, nil)

			suite.slackUserRepository.
				EXPECT().
				GetByPairedHumanAgent(request.HumanAgentCollectionAddress, request.HumanAgentCollectionTokenId).
				Return(nil, nil)

			slackUserToPair := model_builders.NewSlackUserBuilder().
				WithId(request.SlackUserId).
				WithWalletAddress(nil).
				WithPairedHumanAgentCollectionAddress(nil).
				WithPairedHumanAgentCollectionTokenId(nil).
				Build()
			suite.slackUserSynchronizerCached.
				EXPECT().
				Sync(request.SlackUserId).
				Return(slackUserToPair, nil)

			suite.slackUserRepository.
				EXPECT().
				Save(slackUserToPair).
				Return(nil, nil)
		}

		t.Run("Should not return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPairingSlackUserToHumanAgentBeforeEach(t)

			request := generateValidPairSlackUserToHumanAgentRequest()
			initSuite(suite, request)

			err := suite.sut.Pair(request)

			assert.NoError(t, err)
		})
	})
}
