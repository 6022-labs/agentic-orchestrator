package use_cases_test

import (
	"context"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenDetachingParticipantTestingSuite struct {
	sut *use_cases.DetachParticipant

	mockConversationParticipantRepository *repositories_mocks.MockConversationParticipantRepositoryInterface
}

func WhenDetachingParticipantBeforeEach(t *testing.T) *WhenDetachingParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationParticipantRepository := repositories_mocks.NewMockConversationParticipantRepositoryInterface(
		mockController,
	)

	sut := use_cases.NewDetachParticipant(mockConversationParticipantRepository)

	return &WhenDetachingParticipantTestingSuite{
		sut: sut,

		mockConversationParticipantRepository: mockConversationParticipantRepository,
	}
}

func generateValidDetachParticipantRequest() *requests.ConversationParticipantRequest {
	return &requests.ConversationParticipantRequest{
		ConversationID: uuid.New(),
		ParticipantID:  uuid.New(),
	}
}

func TestWhenDetachingParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenDetachingParticipantBeforeEach(t)
			request := generateValidDetachParticipantRequest()
			request.ConversationID = uuid.Nil

			err := suite.sut.Detach(context.Background(), request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "conversationId")
		})
	})

	t.Run("Given repository fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenDetachingParticipantTestingSuite, request *requests.ConversationParticipantRequest) {
			suite.mockConversationParticipantRepository.EXPECT().
				DeleteByConversationIdAndParticipantId(gomock.Any(), request.ConversationID, request.ParticipantID).
				Return(errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDetachingParticipantBeforeEach(t)
			request := generateValidDetachParticipantRequest()
			initSuite(suite, request)

			err := suite.sut.Detach(context.Background(), request)

			assert.Error(t, err)
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenDetachingParticipantTestingSuite, request *requests.ConversationParticipantRequest) {
			suite.mockConversationParticipantRepository.EXPECT().
				DeleteByConversationIdAndParticipantId(gomock.Any(), request.ConversationID, request.ParticipantID).
				Return(nil).
				Times(1)
		}

		t.Run("Should call repository delete", func(t *testing.T) {
			t.Parallel()

			suite := WhenDetachingParticipantBeforeEach(t)
			request := generateValidDetachParticipantRequest()
			initSuite(suite, request)

			err := suite.sut.Detach(context.Background(), request)

			assert.NoError(t, err)
		})
	})
}
