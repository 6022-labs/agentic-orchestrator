package services_test

import (
	"context"

	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenConcludingConversationIfExistsTestingSuite struct {
	sut *services.ConversationConclude

	mockConversationRepository *repositories_mocks.MockConversationRepositoryInterface
}

func WhenConcludingConversationIfExistsBeforeEach(t *testing.T) *WhenConcludingConversationIfExistsTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)

	sut := services.NewConversationConclude(mockConversationRepository)

	return &WhenConcludingConversationIfExistsTestingSuite{
		sut: sut,

		mockConversationRepository: mockConversationRepository,
	}
}

func TestWhenConcludingConversationIfExists(t *testing.T) {
	t.Parallel()

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		conversationId := uuid.New()

		initSuite := func(suite *WhenConcludingConversationIfExistsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversationId).
				Return(nil, nil).
				Times(1)
			// No Save expectation — best-effort: missing conversation is a no-op.
		}

		t.Run("Should leave the system untouched", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludingConversationIfExistsBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ConcludeIfExists(context.Background(), conversationId)

			assert.NoError(t, err)
		})
	})

	t.Run("Given the repository fails on save", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenConcludingConversationIfExistsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludingConversationIfExistsBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ConcludeIfExists(context.Background(), conversation.ID)

			assert.Error(t, err)
		})
	})

	t.Run("Given conversation exists", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenConcludingConversationIfExistsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, c *models.Conversation) (*models.Conversation, error) {
					assert.Equal(t, enums.ConversationStatusConcluded, c.Status)
					assert.NotNil(t, c.ConcludedAt)
					assert.Nil(t, c.UnconcludedAt)
					return c, nil
				}).
				Times(1)
		}

		t.Run("Should mark the conversation as concluded", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludingConversationIfExistsBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ConcludeIfExists(context.Background(), conversation.ID)

			assert.NoError(t, err)
		})
	})

	t.Run("Given conversation is already concluded", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().
			WithStatus(enums.ConversationStatusConcluded).
			Build()

		initSuite := func(suite *WhenConcludingConversationIfExistsTestingSuite) {
			suite.mockConversationRepository.EXPECT().
				GetById(gomock.Any(), conversation.ID).
				Return(conversation, nil).
				Times(1)

			// No Save expectation — idempotent.
		}

		t.Run("Should leave the conversation untouched", func(t *testing.T) {
			t.Parallel()

			suite := WhenConcludingConversationIfExistsBeforeEach(t)
			initSuite(suite)

			err := suite.sut.ConcludeIfExists(context.Background(), conversation.ID)

			assert.NoError(t, err)
		})
	})
}
