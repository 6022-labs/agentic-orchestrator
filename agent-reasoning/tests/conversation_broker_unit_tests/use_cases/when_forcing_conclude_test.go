package use_cases_test

import (
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type WhenForcingConcludeTestingSuite struct {
	sut *use_cases.ForceConclude

	mockConversationRepository *repositories_mocks.MockConversationRepositoryInterface
}

func WhenForcingConcludeBeforeEach(t *testing.T) *WhenForcingConcludeTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)

	sut := use_cases.NewForceConclude(mockConversationRepository)

	return &WhenForcingConcludeTestingSuite{
		sut: sut,

		mockConversationRepository: mockConversationRepository,
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

			err := suite.sut.Conclude(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "conversationId")
		})
	})

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenForcingConcludeTestingSuite, request *requests.ConversationIdRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(request.ConversationID).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenForcingConcludeBeforeEach(t)
			request := generateValidForceConcludeRequest()
			initSuite(suite, request)

			err := suite.sut.Conclude(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given conversation exists", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenForcingConcludeTestingSuite) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(c *models.Conversation) (*models.Conversation, error) {
					assert.Equal(t, enums.ConversationStatusConcluded, c.Status)
					assert.NotNil(t, c.ConcludedAt)
					assert.Nil(t, c.UnconcludedAt)
					return c, nil
				}).
				Times(1)
		}

		t.Run("Should set status=concluded and concluded_at", func(t *testing.T) {
			t.Parallel()

			suite := WhenForcingConcludeBeforeEach(t)
			request := &requests.ConversationIdRequest{ConversationID: conversation.ID}
			initSuite(suite)

			err := suite.sut.Conclude(request)

			assert.NoError(t, err)
		})
	})
}
