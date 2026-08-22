package use_cases_test

import (
	"context"

	"errors"
	"testing"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenGettingConversationMessagesTestingSuite struct {
	sut *use_cases.GetConversationMessages

	mockConversationRepository        *repositories_mocks.MockConversationRepositoryInterface
	mockConversationMessageRepository *repositories_mocks.MockConversationMessageRepositoryInterface
}

func WhenGettingConversationMessagesBeforeEach(t *testing.T) *WhenGettingConversationMessagesTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockConversationMessageRepository := repositories_mocks.NewMockConversationMessageRepositoryInterface(
		mockController,
	)

	sut := use_cases.NewGetConversationMessages(mockConversationRepository, mockConversationMessageRepository)

	return &WhenGettingConversationMessagesTestingSuite{
		sut: sut,

		mockConversationRepository:        mockConversationRepository,
		mockConversationMessageRepository: mockConversationMessageRepository,
	}
}

func TestWhenGettingConversationMessages(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)

			_, err := suite.sut.Get(context.Background(), &requests.ConversationIdRequest{ConversationID: uuid.Nil})

			assert.IsType(t, &cb_errors.ValidationError{}, err)
		})
	})

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenGettingConversationMessagesTestingSuite, request *requests.ConversationIdRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), request.ConversationID).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)
			request := &requests.ConversationIdRequest{ConversationID: uuid.New()}
			initSuite(suite, request)

			_, err := suite.sut.Get(context.Background(), request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given conversation exists", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		message := model_builders.NewConversationMessageBuilder().WithConversationID(conversation.ID).Build()

		initSuite := func(suite *WhenGettingConversationMessagesTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return([]*models.ConversationMessage{message}, nil).
				Times(1)
		}

		t.Run("Should return the conversation messages", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)
			initSuite(suite)

			messages, err := suite.sut.Get(
				context.Background(),
				&requests.ConversationIdRequest{ConversationID: conversation.ID},
			)

			assert.NoError(t, err)
			assert.Len(t, messages, 1)
			assert.Equal(t, message.ID, messages[0].ID)
		})
	})

	t.Run("Given listing messages fails", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenGettingConversationMessagesTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				ListByConversationId(gomock.Any(), conversation.ID).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingConversationMessagesBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Get(
				context.Background(),
				&requests.ConversationIdRequest{ConversationID: conversation.ID},
			)

			assert.Error(t, err)
		})
	})
}
