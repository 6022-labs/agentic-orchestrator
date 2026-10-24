package services_test

import (
	"context"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_external_mocks/services_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/repositories_mocks"
	broker_requester_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSynchronizingSlackAgentTestingSuite struct {
	sut *services.SlackAgentSynchronizer

	mockAgentRequester       *interfaces_mocks.MockAgentRequesterInterface
	mockSlackAgentRepository *repositories_mocks.MockSlackAgentRepositoryInterface
	mockBrokerRequester      *broker_requester_mocks.MockConversationBrokerRequesterInterface
}

func WhenSynchronizingSlackAgentBeforeEach(t *testing.T) *WhenSynchronizingSlackAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentRequester := interfaces_mocks.NewMockAgentRequesterInterface(mockController)
	mockSlackAgentRepository := repositories_mocks.NewMockSlackAgentRepositoryInterface(mockController)
	mockBrokerRequester := broker_requester_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	sut := services.NewSlackAgentSynchronizer(
		mockAgentRequester,
		mockSlackAgentRepository,
		mockBrokerRequester,
	)

	return &WhenSynchronizingSlackAgentTestingSuite{
		sut: sut,

		mockAgentRequester:       mockAgentRequester,
		mockSlackAgentRepository: mockSlackAgentRepository,
		mockBrokerRequester:      mockBrokerRequester,
	}
}

func TestWhenSynchronizingSlackAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given error occurs while getting agent from orchestrator", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite) {
			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackAgentBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), common_tests.RandomEthAddress().Hex(), 123)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while getting agent from the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite) {
			agentFromOrchestrator := model_builders.NewSlackAgentBuilder().Build()
			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(agentFromOrchestrator, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, assert.AnError)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSynchronizingSlackAgentBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Sync(context.Background(), common_tests.RandomEthAddress().Hex(), 123)

			assert.Error(t, err)
		})
	})

	t.Run("Given agent is not found in the database", func(t *testing.T) {
		t.Parallel()

		brokerParticipantId := uuid.New()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite) {
			agentFromOrchestrator := model_builders.NewSlackAgentBuilder().Build()
			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(agentFromOrchestrator, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil)

			suite.mockBrokerRequester.
				EXPECT().
				RegisterAgent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(brokerParticipantId, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, saved *models.SlackAgent) (*models.SlackAgent, error) {
					assert.NotNil(t, saved.BrokerParticipantId)
					assert.Equal(t, brokerParticipantId, *saved.BrokerParticipantId)
					return saved, nil
				})
		}

		t.Run(
			"Should register the agent with the broker, store the broker participant id, and persist the row",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenSynchronizingSlackAgentBeforeEach(t)
				initSuite(suite)

				_, err := suite.sut.Sync(context.Background(), common_tests.RandomEthAddress().Hex(), 123)

				assert.Nil(t, err)
			},
		)
	})

	t.Run("Given agent is found in the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSynchronizingSlackAgentTestingSuite, collectionAddress string, collectionTokenId uint64) {
			updatedCid := "initial-cid"
			agentFromOrchestrator := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				WithIconImageCid(&updatedCid).
				WithOwnerWallet(common_tests.RandomEthAddress().Hex()).
				Build()

			suite.mockAgentRequester.
				EXPECT().
				GetAgent(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(agentFromOrchestrator, nil)

			agentInDatabase := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				WithIconImageCid(nil).
				WithOwnerWallet(common_tests.RandomEthAddress().Hex()).
				Build()
			suite.mockSlackAgentRepository.
				EXPECT().
				GetByCollectionAddressAndCollectionTokenId(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(agentInDatabase, nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				UpdateOwnerWallet(
					gomock.Any(),
					strings.ToLower(collectionAddress),
					collectionTokenId,
					strings.ToLower(agentFromOrchestrator.OwnerWallet),
				).
				Return(nil)

			suite.mockSlackAgentRepository.
				EXPECT().
				UpdateIconImageCid(
					gomock.Any(),
					strings.ToLower(collectionAddress),
					collectionTokenId,
					agentFromOrchestrator.IconImageCid,
				).
				Return(nil)
		}

		t.Run(
			"Should update the agent under its lowercased address with its owner wallet lowercased",
			func(t *testing.T) {
				t.Parallel()

				suite := WhenSynchronizingSlackAgentBeforeEach(t)

				collectionAddress := common_tests.RandomEthAddress().Hex()
				collectionTokenId := uint64(123)

				initSuite(suite, collectionAddress, collectionTokenId)

				agent, err := suite.sut.Sync(context.Background(), collectionAddress, collectionTokenId)

				assert.Nil(t, err)
				assert.NotNil(t, agent)
			},
		)
	})
}
