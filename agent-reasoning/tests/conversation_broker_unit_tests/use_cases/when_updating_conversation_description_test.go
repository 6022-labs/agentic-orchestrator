package use_cases_test

import (
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

type WhenUpdatingConversationDescriptionTestingSuite struct {
	sut *use_cases.UpdateConversationDescription

	mockConversationRepository *repositories_mocks.MockConversationRepositoryInterface
}

func WhenUpdatingConversationDescriptionBeforeEach(t *testing.T) *WhenUpdatingConversationDescriptionTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)

	sut := use_cases.NewUpdateConversationDescription(mockConversationRepository)

	return &WhenUpdatingConversationDescriptionTestingSuite{
		sut: sut,

		mockConversationRepository: mockConversationRepository,
	}
}

func generateValidUpdateConversationDescriptionRequest() *requests.UpdateConversationDescriptionRequest {
	desc := "new description"
	return &requests.UpdateConversationDescriptionRequest{
		ConversationID: uuid.New(),
		Description:    &desc,
	}
}

func TestWhenUpdatingConversationDescription(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingConversationDescriptionBeforeEach(t)
			request := generateValidUpdateConversationDescriptionRequest()
			request.ConversationID = uuid.Nil

			err := suite.sut.Update(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "conversationId")
		})
	})

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenUpdatingConversationDescriptionTestingSuite, request *requests.UpdateConversationDescriptionRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(request.ConversationID).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingConversationDescriptionBeforeEach(t)
			request := generateValidUpdateConversationDescriptionRequest()
			initSuite(suite, request)

			err := suite.sut.Update(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given conversation exists and a new description is provided", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		newDescription := "renamed channel"

		initSuite := func(suite *WhenUpdatingConversationDescriptionTestingSuite) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(c *models.Conversation) (*models.Conversation, error) {
					assert.NotNil(t, c.Description)
					assert.Equal(t, newDescription, *c.Description)
					return c, nil
				}).
				Times(1)
		}

		t.Run("Should overwrite the description and persist", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingConversationDescriptionBeforeEach(t)
			request := &requests.UpdateConversationDescriptionRequest{
				ConversationID: conversation.ID,
				Description:    &newDescription,
			}
			initSuite(suite)

			err := suite.sut.Update(request)

			assert.NoError(t, err)
		})
	})

	t.Run("Given conversation exists and description is nil", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenUpdatingConversationDescriptionTestingSuite) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(c *models.Conversation) (*models.Conversation, error) {
					assert.Nil(t, c.Description)
					return c, nil
				}).
				Times(1)
		}

		t.Run("Should clear the description", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingConversationDescriptionBeforeEach(t)
			request := &requests.UpdateConversationDescriptionRequest{
				ConversationID: conversation.ID,
				Description:    nil,
			}
			initSuite(suite)

			err := suite.sut.Update(request)

			assert.NoError(t, err)
		})
	})
}
