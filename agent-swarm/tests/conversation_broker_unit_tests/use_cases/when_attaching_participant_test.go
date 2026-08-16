package use_cases_test

import (
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

type WhenAttachingParticipantTestingSuite struct {
	sut *use_cases.AttachParticipant

	mockConversationRepository            *repositories_mocks.MockConversationRepositoryInterface
	mockParticipantRepository             *repositories_mocks.MockParticipantRepositoryInterface
	mockConversationParticipantRepository *repositories_mocks.MockConversationParticipantRepositoryInterface
}

func WhenAttachingParticipantBeforeEach(t *testing.T) *WhenAttachingParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)
	mockConversationParticipantRepository := repositories_mocks.NewMockConversationParticipantRepositoryInterface(
		mockController,
	)

	sut := use_cases.NewAttachParticipant(
		mockConversationRepository,
		mockParticipantRepository,
		mockConversationParticipantRepository,
	)

	return &WhenAttachingParticipantTestingSuite{
		sut: sut,

		mockConversationRepository:            mockConversationRepository,
		mockParticipantRepository:             mockParticipantRepository,
		mockConversationParticipantRepository: mockConversationParticipantRepository,
	}
}

func generateValidConversationParticipantRequest() *requests.ConversationParticipantRequest {
	return &requests.ConversationParticipantRequest{
		ConversationID: uuid.New(),
		ParticipantID:  uuid.New(),
	}
}

func TestWhenAttachingParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because conversationId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			request := generateValidConversationParticipantRequest()
			request.ConversationID = uuid.Nil

			err := suite.sut.Attach(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "conversationId")
		})

		t.Run("Because participantId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			request := generateValidConversationParticipantRequest()
			request.ParticipantID = uuid.Nil

			err := suite.sut.Attach(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "participantId")
		})
	})

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenAttachingParticipantTestingSuite, request *requests.ConversationParticipantRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(request.ConversationID).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			request := generateValidConversationParticipantRequest()
			initSuite(suite, request)

			err := suite.sut.Attach(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given participant does not exist", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenAttachingParticipantTestingSuite, request *requests.ConversationParticipantRequest) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().
				GetById(request.ParticipantID).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			request := &requests.ConversationParticipantRequest{
				ConversationID: conversation.ID,
				ParticipantID:  uuid.New(),
			}
			initSuite(suite, request)

			err := suite.sut.Attach(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given conversation and participant exist", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		participant := model_builders.NewAgentParticipantBuilder().Build()

		initSuite := func(suite *WhenAttachingParticipantTestingSuite) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().GetById(participant.ID).Return(participant, nil).Times(1)
			suite.mockConversationParticipantRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(link *models.ConversationParticipant) (*models.ConversationParticipant, error) {
					assert.Equal(t, conversation.ID, link.ConversationID)
					assert.Equal(t, participant.ID, link.ParticipantID)
					return link, nil
				}).
				Times(1)
		}

		t.Run("Should save the link", func(t *testing.T) {
			t.Parallel()

			suite := WhenAttachingParticipantBeforeEach(t)
			request := &requests.ConversationParticipantRequest{
				ConversationID: conversation.ID,
				ParticipantID:  participant.ID,
			}
			initSuite(suite)

			err := suite.sut.Attach(request)

			assert.NoError(t, err)
		})
	})
}
