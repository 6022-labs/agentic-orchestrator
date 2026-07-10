package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenListingConversationParticipantsTestingSuite struct {
	sut *use_cases.ListConversationParticipants

	mockConversationRepository *repositories_mocks.MockConversationRepositoryInterface
	mockParticipantRepository  *repositories_mocks.MockParticipantRepositoryInterface
}

func WhenListingConversationParticipantsBeforeEach(t *testing.T) *WhenListingConversationParticipantsTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	sut := use_cases.NewListConversationParticipants(mockConversationRepository, mockParticipantRepository)

	return &WhenListingConversationParticipantsTestingSuite{
		sut: sut,

		mockConversationRepository: mockConversationRepository,
		mockParticipantRepository:  mockParticipantRepository,
	}
}

func TestWhenListingConversationParticipants(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)

			_, err := suite.sut.List(&requests.ConversationIdRequest{ConversationID: uuid.Nil})

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenListingConversationParticipantsTestingSuite, request *requests.ConversationIdRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(request.ConversationID).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)
			request := &requests.ConversationIdRequest{ConversationID: uuid.New()}
			initSuite(suite, request)

			_, err := suite.sut.List(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given conversation exists", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		participant := model_builders.NewAgentParticipantBuilder().Build()

		initSuite := func(suite *WhenListingConversationParticipantsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(conversation.ID).
				Return([]*models.Participant{participant}, nil).
				Times(1)
		}

		t.Run("Should return the conversation participants", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)
			initSuite(suite)

			participants, err := suite.sut.List(&requests.ConversationIdRequest{ConversationID: conversation.ID})

			assert.NoError(t, err)
			assert.Len(t, participants, 1)
			assert.Equal(t, participant.ID, participants[0].ID)
		})
	})

	t.Run("Given listing participants fails", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenListingConversationParticipantsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ListByConversationId(conversation.ID).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenListingConversationParticipantsBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.List(&requests.ConversationIdRequest{ConversationID: conversation.ID})

			assert.Error(t, err)
		})
	})
}
