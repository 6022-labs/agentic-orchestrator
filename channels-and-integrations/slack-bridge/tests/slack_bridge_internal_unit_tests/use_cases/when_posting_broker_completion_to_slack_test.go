package use_cases_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenPostingBrokerCompletionToSlackTestingSuite struct {
	sut *use_cases.PostBrokerCompletionToSlack

	mockSlackConversationRepository *repositories_mocks.MockSlackConversationRepositoryInterface
	mockSlackAgentRepository        *repositories_mocks.MockSlackAgentRepositoryInterface
	mockSlackThreadRequester        *interfaces_mocks.MockSlackThreadRequesterInterface
}

func WhenPostingBrokerCompletionToSlackBeforeEach(t *testing.T) *WhenPostingBrokerCompletionToSlackTestingSuite {
	mockController := gomock.NewController(t)

	mockSlackConversationRepository := repositories_mocks.NewMockSlackConversationRepositoryInterface(mockController)
	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)
	mockSlackThreadRequester := interfaces_mocks.NewMockSlackThreadRequesterInterface(mockController)

	sut := use_cases.NewPostBrokerCompletionToSlack(
		mockSlackConversationRepository,
		mockSlackAgentRepository,
		mockSlackThreadRequester,
	)

	return &WhenPostingBrokerCompletionToSlackTestingSuite{
		sut:                             sut,
		mockSlackConversationRepository: mockSlackConversationRepository,
		mockSlackAgentRepository:        mockSlackAgentRepository,
		mockSlackThreadRequester:        mockSlackThreadRequester,
	}
}

func generateValidPostBrokerCompletionToSlackRequest() *requests.PostBrokerCompletionToSlackRequest {
	return &requests.PostBrokerCompletionToSlackRequest{
		ConversationId:    uuid.New(),
		CollectionAddress: strings.ToLower(common_tests.RandomEthAddress().Hex()),
		CollectionTokenId: 7,
		Content:           "hello back",
	}
}

func TestWhenPostingBrokerCompletionToSlack(t *testing.T) {
	t.Parallel()

	t.Run("Given the request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because collectionAddress has no 0x prefix", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingBrokerCompletionToSlackBeforeEach(t)
			req := generateValidPostBrokerCompletionToSlackRequest()
			req.CollectionAddress = strings.TrimPrefix(req.CollectionAddress, "0x")

			err := suite.sut.Post(context.Background(), req)
			assert.ErrorIs(t, err, common_errors.ErrValidation)
		})
	})

	t.Run("Given a checksummed collection address", func(t *testing.T) {
		t.Parallel()

		checksumAddress := common_tests.RandomEthAddress().Hex()

		initSuite := func(suite *WhenPostingBrokerCompletionToSlackTestingSuite, req *requests.PostBrokerCompletionToSlackRequest) {
			suite.mockSlackConversationRepository.EXPECT().
				GetByBrokerConversationId(gomock.Any(), req.ConversationId).
				Return(&models.SlackConversation{
					SlackChannelId:       "C-target",
					SlackThreadTs:        "1.0",
					BrokerConversationID: req.ConversationId,
				}, nil)
			suite.mockSlackAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), strings.ToLower(checksumAddress), req.CollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should look the agent up by its lowercased address", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingBrokerCompletionToSlackBeforeEach(t)
			req := generateValidPostBrokerCompletionToSlackRequest()
			req.CollectionAddress = checksumAddress
			initSuite(suite, req)

			err := suite.sut.Post(context.Background(), req)
			assert.ErrorIs(t, err, use_cases.ErrSlackAgentNotFound)
		})
	})

	t.Run("Given the conversation link does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingBrokerCompletionToSlackTestingSuite, req *requests.PostBrokerCompletionToSlackRequest) {
			suite.mockSlackConversationRepository.EXPECT().
				GetByBrokerConversationId(gomock.Any(), req.ConversationId).
				Return(nil, nil)
		}

		t.Run("Should return ErrConversationLinkNotFound", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingBrokerCompletionToSlackBeforeEach(t)
			req := generateValidPostBrokerCompletionToSlackRequest()
			initSuite(suite, req)

			err := suite.sut.Post(context.Background(), req)
			assert.ErrorIs(t, err, use_cases.ErrConversationLinkNotFound)
		})
	})

	t.Run("Given the slack agent does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingBrokerCompletionToSlackTestingSuite, req *requests.PostBrokerCompletionToSlackRequest) {
			suite.mockSlackConversationRepository.EXPECT().
				GetByBrokerConversationId(gomock.Any(), req.ConversationId).
				Return(&models.SlackConversation{
					SlackChannelId:       "C-target",
					SlackThreadTs:        "1.0",
					BrokerConversationID: req.ConversationId,
				}, nil)
			suite.mockSlackAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), req.CollectionAddress, req.CollectionTokenId).
				Return(nil, nil)
		}

		t.Run("Should return ErrSlackAgentNotFound", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingBrokerCompletionToSlackBeforeEach(t)
			req := generateValidPostBrokerCompletionToSlackRequest()
			initSuite(suite, req)

			err := suite.sut.Post(context.Background(), req)
			assert.ErrorIs(t, err, use_cases.ErrSlackAgentNotFound)
		})
	})

	t.Run("Given posting to slack fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingBrokerCompletionToSlackTestingSuite, req *requests.PostBrokerCompletionToSlackRequest) {
			suite.mockSlackConversationRepository.EXPECT().
				GetByBrokerConversationId(gomock.Any(), req.ConversationId).
				Return(&models.SlackConversation{
					SlackChannelId:       "C-target",
					SlackThreadTs:        "1.0",
					BrokerConversationID: req.ConversationId,
				}, nil)
			slackAgent := &models.SlackAgent{
				CollectionAddress: req.CollectionAddress,
				CollectionTokenId: req.CollectionTokenId,
			}
			suite.mockSlackAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), req.CollectionAddress, req.CollectionTokenId).
				Return(slackAgent, nil)
			suite.mockSlackThreadRequester.EXPECT().
				PushCompletionToThread(gomock.Any(), req.Content, slackAgent, "C-target", "1.0").
				Return(errors.New("slack down"))
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingBrokerCompletionToSlackBeforeEach(t)
			req := generateValidPostBrokerCompletionToSlackRequest()
			initSuite(suite, req)

			err := suite.sut.Post(context.Background(), req)
			assert.Error(t, err)
		})
	})

	t.Run("Given everything succeeds", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenPostingBrokerCompletionToSlackTestingSuite, req *requests.PostBrokerCompletionToSlackRequest) {
			suite.mockSlackConversationRepository.EXPECT().
				GetByBrokerConversationId(gomock.Any(), req.ConversationId).
				Return(&models.SlackConversation{
					SlackChannelId:       "C-target",
					SlackThreadTs:        "1.0",
					BrokerConversationID: req.ConversationId,
				}, nil)
			slackAgent := &models.SlackAgent{
				CollectionAddress: req.CollectionAddress,
				CollectionTokenId: req.CollectionTokenId,
			}
			suite.mockSlackAgentRepository.EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), req.CollectionAddress, req.CollectionTokenId).
				Return(slackAgent, nil)
			suite.mockSlackThreadRequester.EXPECT().
				PushCompletionToThread(gomock.Any(), req.Content, slackAgent, "C-target", "1.0").
				Return(nil)
		}

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenPostingBrokerCompletionToSlackBeforeEach(t)
			req := generateValidPostBrokerCompletionToSlackRequest()
			initSuite(suite, req)

			err := suite.sut.Post(context.Background(), req)
			assert.NoError(t, err)
		})
	})
}
