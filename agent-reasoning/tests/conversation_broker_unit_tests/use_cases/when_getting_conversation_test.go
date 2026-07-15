package use_cases_test

import (
	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingConversationTestingSuite struct {
	sut *use_cases.GetConversation

	mockConversationRepository *repositories_mocks.MockConversationRepositoryInterface
}

func WhenGettingConversationBeforeEach(t *testing.T) *WhenGettingConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)

	sut := use_cases.NewGetConversation(mockConversationRepository)

	return &WhenGettingConversationTestingSuite{
		sut: sut,

		mockConversationRepository: mockConversationRepository,
	}
}

func generateValidGetConversationRequest() *requests.ConversationIdRequest {
	return &requests.ConversationIdRequest{
		ConversationID: uuid.New(),
	}
}

func TestWhenGettingConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			request := generateValidGetConversationRequest()
			request.ConversationID = uuid.Nil

			_, err := suite.sut.Get(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "conversationId")
		})
	})

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationTestingSuite, request *requests.ConversationIdRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(request.ConversationID).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			request := generateValidGetConversationRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given repository fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationTestingSuite, request *requests.ConversationIdRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(request.ConversationID).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			request := generateValidGetConversationRequest()
			initSuite(suite, request)

			_, err := suite.sut.Get(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given conversation exists", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenGettingConversationTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(conversation.ID).
				Return(conversation, nil).
				Times(1)
		}

		t.Run("Should return the conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationBeforeEach(t)
			request := &requests.ConversationIdRequest{ConversationID: conversation.ID}
			initSuite(suite)

			got, err := suite.sut.Get(request)

			assert.NoError(t, err)
			assert.Equal(t, conversation.ID, got.ID)
		})
	})
}
