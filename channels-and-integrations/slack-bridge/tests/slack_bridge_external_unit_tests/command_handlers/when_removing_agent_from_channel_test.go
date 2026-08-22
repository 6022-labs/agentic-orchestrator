package command_handlers_test

import (
	"context"

	"errors"
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
	"go.uber.org/zap"
)

type WhenRemovingAgentFromChannelTestingSuite struct {
	sut *command_handlers.RemoveAgentFromChannel

	mockSlackAgentRepository        *repositories_mocks.MockSlackAgentRepositoryInterface
	mockSlackChannelSynchronizer    *interfaces_mocks.MockSlackChannelSynchronizerInterface
	mockSlackChannelRepository      *repositories_mocks.MockSlackChannelRepositoryInterface
	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockBrokerRequester             *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenRemovingAgentFromChannelBeforeEach(t *testing.T) *WhenRemovingAgentFromChannelTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)
	mockSlackChannelSynchronizer := interfaces_mocks.NewMockSlackChannelSynchronizerInterface(mockController)
	mockSlackChannelRepository := repositories_mocks.NewMockSlackChannelRepositoryInterface(mockController)
	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := command_handlers.NewRemoveAgentFromChannel(
		zap.NewNop(),
		mockSlackAgentRepository,
		mockSlackChannelRepository,
		mockSlackChannelSynchronizer,
		mockSlackConversationRepository,
		mockBrokerRequester,
	)

	return &WhenRemovingAgentFromChannelTestingSuite{
		sut:                             sut,
		mockSlackAgentRepository:        mockSlackAgentRepository,
		mockSlackChannelSynchronizer:    mockSlackChannelSynchronizer,
		mockSlackChannelRepository:      mockSlackChannelRepository,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockBrokerRequester:             mockBrokerRequester,
	}
}

func TestWhenRemovingAgentFromChannel(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting the channel from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
			suite.mockSlackChannelRepository.
				EXPECT().
				GetById(gomock.Any(), request.ChannelID).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRemovingAgentFromChannelBeforeEach(t)

			request := &requests.RemoveAgentFromChannelRequest{
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given the channel did not contain the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
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
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given an error occurs while syncing the channel", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
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
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent from repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
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
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent not found in repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
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
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the channel is not in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
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
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given error occurs while saving channel after removal", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRemovingAgentFromChannelTestingSuite, request *requests.RemoveAgentFromChannelRequest) {
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
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
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
				ChannelID:         "chan",
				CollectionTokenId: 123,
				CollectionAddress: common_tests.RandomEthAddress().Hex(),
			}

			initSuite(suite, request)

			err := suite.sut.Remove(context.Background(), request)
			assert.NoError(t, err)
		})
	})
}
