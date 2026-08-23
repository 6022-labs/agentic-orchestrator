package use_cases_test

import (
	"context"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	broker_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenForcingConcludeTestingSuite struct {
	sut *use_cases.ForceConclude

	mockConversationConclude *interfaces_mocks.MockConversationConcludeInterface
	mockBrokerMetrics        *broker_metrics_mocks.MockBrokerMetricsInterface
}

func WhenForcingConcludeBeforeEach(t *testing.T) *WhenForcingConcludeTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationConclude := interfaces_mocks.NewMockConversationConcludeInterface(mockController)

	mockBrokerMetrics := broker_metrics_mocks.NewMockBrokerMetricsInterface(mockController)

	sut := use_cases.NewForceConclude(mockBrokerMetrics, mockConversationConclude)

	return &WhenForcingConcludeTestingSuite{
		sut: sut,

		mockConversationConclude: mockConversationConclude,
		mockBrokerMetrics:        mockBrokerMetrics,
	}
}

func generateValidForceConcludeRequest() *requests.ConversationIdRequest {
	return &requests.ConversationIdRequest{
		ConversationID: uuid.New(),
	}
}

func TestWhenForcingConclude(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenForcingConcludeBeforeEach(t)
			request := generateValidForceConcludeRequest()
			request.ConversationID = uuid.Nil

			err := suite.sut.Conclude(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "conversationId")
		})
	})

	t.Run("Given the conclude service fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForcingConcludeTestingSuite, request *requests.ConversationIdRequest) {
			suite.mockConversationConclude.EXPECT().
				ConcludeIfExists(gomock.Any(), request.ConversationID).
				Return(false, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForcingConcludeBeforeEach(t)
			request := generateValidForceConcludeRequest()
			initSuite(suite, request)

			err := suite.sut.Conclude(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForcingConcludeTestingSuite, request *requests.ConversationIdRequest) {
			suite.mockConversationConclude.EXPECT().
				ConcludeIfExists(gomock.Any(), request.ConversationID).
				Return(true, nil).
				Times(1)
		}

		t.Run("Should delegate to the conclude service", func(t *testing.T) {
			t.Parallel()

			suite := WhenForcingConcludeBeforeEach(t)
			suite.mockBrokerMetrics.EXPECT().
				RecordConversationConcluded(gomock.Any(), metrics_interfaces.ConclusionReasonForced)
			request := generateValidForceConcludeRequest()
			initSuite(suite, request)

			err := suite.sut.Conclude(context.Background(), request)

			assert.NoError(t, err)
		})
	})
}
