package services_test

import (
	"context"
	"strings"
	"testing"

	gateway_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenEnsuringConversationAccessTestingSuite struct {
	sut *services.ConversationAccessVerifier

	mockAgentAccessVerifier         *interfaces_mocks.MockAgentAccessVerifierInterface
	mockConversationBrokerRequester *interfaces_mocks.MockConversationBrokerRequesterInterface
}

func WhenEnsuringConversationAccessBeforeEach(
	t *testing.T,
) *WhenEnsuringConversationAccessTestingSuite {
	mockController := gomock.NewController(t)

	mockAgentAccessVerifier := interfaces_mocks.NewMockAgentAccessVerifierInterface(mockController)
	mockConversationBrokerRequester := interfaces_mocks.NewMockConversationBrokerRequesterInterface(mockController)

	return &WhenEnsuringConversationAccessTestingSuite{
		sut: services.NewConversationAccessVerifier(
			zap.NewNop(),
			mockAgentAccessVerifier,
			mockConversationBrokerRequester,
		),

		mockAgentAccessVerifier:         mockAgentAccessVerifier,
		mockConversationBrokerRequester: mockConversationBrokerRequester,
	}
}

func TestWhenEnsuringConversationAccess(t *testing.T) {
	t.Parallel()

	agentParticipant := func(collectionAddress string, collectionTokenId uint64) *responses.BrokerParticipantResponse {
		return &responses.BrokerParticipantResponse{
			ID:                uuid.New(),
			CollectionAddress: lo.ToPtr(collectionAddress),
			CollectionTokenId: lo.ToPtr(collectionTokenId),
		}
	}

	t.Run("Given the conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringConversationAccessTestingSuite, conversationId uuid.UUID) {
			suite.mockConversationBrokerRequester.EXPECT().
				ListConversationParticipants(gomock.Any(), conversationId).
				Return(nil, gateway_errors.NewNotFoundError("conversation_not_found", "conversation not found")).
				Times(1)
		}

		t.Run("Should return the not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringConversationAccessBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				conversationId,
			)

			var notFoundErr *gateway_errors.NotFoundError
			assert.ErrorAs(t, err, &notFoundErr)
			assert.Equal(t, "conversation_not_found", notFoundErr.Code)
		})
	})

	t.Run("Given the broker fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenEnsuringConversationAccessTestingSuite, conversationId uuid.UUID) {
			suite.mockConversationBrokerRequester.EXPECT().
				ListConversationParticipants(gomock.Any(), conversationId).
				Return(nil, assert.AnError).
				Times(1)
		}

		t.Run("Should return an unavailable error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringConversationAccessBeforeEach(t)
			conversationId := uuid.New()
			initSuite(suite, conversationId)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				conversationId,
			)

			var unavailable *gateway_errors.UnavailableError
			assert.ErrorAs(t, err, &unavailable)
			assert.Equal(t, "conversation_broker_unavailable", unavailable.Code)
			assert.NotContains(t, err.Error(), assert.AnError.Error())
		})
	})

	t.Run("Given the conversation has no agent participants", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenEnsuringConversationAccessTestingSuite,
			conversationId uuid.UUID,
			collectionAddress string,
		) {
			suite.mockConversationBrokerRequester.EXPECT().
				ListConversationParticipants(gomock.Any(), conversationId).
				Return([]*responses.BrokerParticipantResponse{
					nil,
					{ID: uuid.New()},
					{ID: uuid.New(), CollectionAddress: lo.ToPtr(collectionAddress)},
					{ID: uuid.New(), CollectionTokenId: lo.ToPtr(uint64(1))},
				}, nil).
				Times(1)
		}

		t.Run("Should return a forbidden error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringConversationAccessBeforeEach(t)
			conversationId := uuid.New()
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, conversationId, collectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				conversationId,
			)

			var forbiddenErr *gateway_errors.ForbiddenError
			assert.ErrorAs(t, err, &forbiddenErr)
			assert.Equal(t, services.ConversationAccessForbiddenCode, forbiddenErr.Code)
		})
	})

	t.Run("Given the registry fails while checking a participant", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenEnsuringConversationAccessTestingSuite,
			conversationId uuid.UUID,
			collectionAddress string,
		) {
			suite.mockConversationBrokerRequester.EXPECT().
				ListConversationParticipants(gomock.Any(), conversationId).
				Return([]*responses.BrokerParticipantResponse{agentParticipant(collectionAddress, 1)}, nil).
				Times(1)
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), gomock.Any(), collectionAddress, uint64(1)).
				Return(gateway_errors.NewUnavailableError("agent_registry_unavailable", "registry down")).
				Times(1)
		}

		t.Run("Should return the error instead of forbidding", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringConversationAccessBeforeEach(t)
			conversationId := uuid.New()
			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, conversationId, collectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				conversationId,
			)

			var unavailable *gateway_errors.UnavailableError
			assert.ErrorAs(t, err, &unavailable)
			assert.Equal(t, "agent_registry_unavailable", unavailable.Code)
		})
	})

	t.Run("Given the caller has no access to any participating agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenEnsuringConversationAccessTestingSuite,
			conversationId uuid.UUID,
			firstCollectionAddress string,
			secondCollectionAddress string,
		) {
			suite.mockConversationBrokerRequester.EXPECT().
				ListConversationParticipants(gomock.Any(), conversationId).
				Return([]*responses.BrokerParticipantResponse{
					agentParticipant(firstCollectionAddress, 1),
					agentParticipant(secondCollectionAddress, 2),
				}, nil).
				Times(1)
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), gomock.Any(), firstCollectionAddress, uint64(1)).
				Return(gateway_errors.NewForbiddenError(services.AgentAccessForbiddenCode, "forbidden")).
				Times(1)
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), gomock.Any(), secondCollectionAddress, uint64(2)).
				Return(gateway_errors.NewNotFoundError("agent_not_found", "agent not found")).
				Times(1)
		}

		t.Run("Should return a forbidden error", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringConversationAccessBeforeEach(t)
			conversationId := uuid.New()
			firstCollectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			secondCollectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, conversationId, firstCollectionAddress, secondCollectionAddress)

			err := suite.sut.EnsureAccess(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				conversationId,
			)

			var forbiddenErr *gateway_errors.ForbiddenError
			assert.ErrorAs(t, err, &forbiddenErr)
			assert.Equal(t, services.ConversationAccessForbiddenCode, forbiddenErr.Code)
		})
	})

	t.Run("Given the caller has access to one of the participating agents", func(t *testing.T) {
		t.Parallel()

		initSuite := func(
			suite *WhenEnsuringConversationAccessTestingSuite,
			callerId string,
			conversationId uuid.UUID,
			firstCollectionAddress string,
			secondCollectionAddress string,
		) {
			suite.mockConversationBrokerRequester.EXPECT().
				ListConversationParticipants(gomock.Any(), conversationId).
				Return([]*responses.BrokerParticipantResponse{
					{ID: uuid.New()},
					agentParticipant(firstCollectionAddress, 1),
					agentParticipant(secondCollectionAddress, 2),
				}, nil).
				Times(1)
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, firstCollectionAddress, uint64(1)).
				Return(gateway_errors.NewForbiddenError(services.AgentAccessForbiddenCode, "forbidden")).
				Times(1)
			suite.mockAgentAccessVerifier.EXPECT().
				EnsureAccess(gomock.Any(), callerId, secondCollectionAddress, uint64(2)).
				Return(nil).
				Times(1)
		}

		t.Run("Should succeed", func(t *testing.T) {
			t.Parallel()

			suite := WhenEnsuringConversationAccessBeforeEach(t)
			callerId := strings.ToLower(common_tests.RandomEthAddress().Hex())
			conversationId := uuid.New()
			firstCollectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			secondCollectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			initSuite(suite, callerId, conversationId, firstCollectionAddress, secondCollectionAddress)

			err := suite.sut.EnsureAccess(context.Background(), callerId, conversationId)

			assert.NoError(t, err)
		})
	})
}
