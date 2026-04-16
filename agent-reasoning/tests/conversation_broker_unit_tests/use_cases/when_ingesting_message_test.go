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

type WhenIngestingMessageTestingSuite struct {
	sut *use_cases.IngestMessage

	mockConversationRepository        *repositories_mocks.MockConversationRepositoryInterface
	mockParticipantRepository         *repositories_mocks.MockParticipantRepositoryInterface
	mockConversationMessageRepository *repositories_mocks.MockConversationMessageRepositoryInterface
}

func WhenIngestingMessageBeforeEach(t *testing.T) *WhenIngestingMessageTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)
	mockConversationMessageRepository := repositories_mocks.NewMockConversationMessageRepositoryInterface(mockController)

	sut := use_cases.NewIngestMessage(
		mockConversationRepository,
		mockParticipantRepository,
		mockConversationMessageRepository,
	)

	return &WhenIngestingMessageTestingSuite{
		sut: sut,

		mockConversationRepository:        mockConversationRepository,
		mockParticipantRepository:         mockParticipantRepository,
		mockConversationMessageRepository: mockConversationMessageRepository,
	}
}

func generateValidIngestMessageRequest() *requests.IngestMessageRequest {
	return &requests.IngestMessageRequest{
		ConversationID: uuid.New(),
		ParticipantID:  uuid.New(),
		Content:        "hello",
	}
}

func TestWhenIngestingMessage(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because content is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := generateValidIngestMessageRequest()
			request.Content = ""

			_, err := suite.sut.Ingest(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "content")
		})

		t.Run("Because participantId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := generateValidIngestMessageRequest()
			request.ParticipantID = uuid.Nil

			_, err := suite.sut.Ingest(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "participantId")
		})
	})

	t.Run("Given conversation does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenIngestingMessageTestingSuite, request *requests.IngestMessageRequest) {
			suite.mockConversationRepository.EXPECT().
				GetById(request.ConversationID).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := generateValidIngestMessageRequest()
			initSuite(suite, request)

			_, err := suite.sut.Ingest(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given participant does not exist", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenIngestingMessageTestingSuite, request *requests.IngestMessageRequest) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().
				GetById(request.ParticipantID).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := &requests.IngestMessageRequest{
				ConversationID: conversation.ID,
				ParticipantID:  uuid.New(),
				Content:        "hi",
			}
			initSuite(suite, request)

			_, err := suite.sut.Ingest(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given conversation is concluded and message is from an unpaired user", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		conversation.Status = enums.ConversationStatusConcluded
		userParticipant := model_builders.NewUserParticipantBuilder().Build()

		initSuite := func(suite *WhenIngestingMessageTestingSuite, request *requests.IngestMessageRequest) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().GetById(userParticipant.ID).Return(userParticipant, nil).Times(1)
			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(c *models.Conversation) (*models.Conversation, error) {
					assert.Equal(t, enums.ConversationStatusActive, c.Status)
					assert.NotNil(t, c.UnconcludedAt)
					return c, nil
				}).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().Save(gomock.Any()).
				DoAndReturn(func(m *models.ConversationMessage) (*models.ConversationMessage, error) { return m, nil }).
				Times(1)
		}

		t.Run("Should reopen the conversation and ingest the message", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := &requests.IngestMessageRequest{
				ConversationID: conversation.ID,
				ParticipantID:  userParticipant.ID,
				Content:        "back online",
			}
			initSuite(suite, request)

			message, err := suite.sut.Ingest(request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, message.Content)
		})
	})

	t.Run("Given conversation is concluded and message is from a paired human-agent", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		conversation.Status = enums.ConversationStatusConcluded
		humanAgent := model_builders.NewAgentParticipantBuilder().WithRole(enums.ParticipantRoleHuman).Build()

		initSuite := func(suite *WhenIngestingMessageTestingSuite, request *requests.IngestMessageRequest) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().GetById(humanAgent.ID).Return(humanAgent, nil).Times(1)
			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(c *models.Conversation) (*models.Conversation, error) {
					assert.Equal(t, enums.ConversationStatusActive, c.Status)
					assert.NotNil(t, c.UnconcludedAt)
					return c, nil
				}).
				Times(1)
			suite.mockConversationMessageRepository.EXPECT().Save(gomock.Any()).
				DoAndReturn(func(m *models.ConversationMessage) (*models.ConversationMessage, error) { return m, nil }).
				Times(1)
		}

		t.Run("Should reopen the conversation and ingest the message", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := &requests.IngestMessageRequest{
				ConversationID: conversation.ID,
				ParticipantID:  humanAgent.ID,
				Content:        "back online",
			}
			initSuite(suite, request)

			message, err := suite.sut.Ingest(request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, message.Content)
		})
	})

	t.Run("Given conversation is concluded and message is from an AI agent", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		conversation.Status = enums.ConversationStatusConcluded
		aiAgent := model_builders.NewAgentParticipantBuilder().WithRole(enums.ParticipantRoleClone).Build()

		initSuite := func(suite *WhenIngestingMessageTestingSuite, request *requests.IngestMessageRequest) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().GetById(aiAgent.ID).Return(aiAgent, nil).Times(1)
			// No Save on the conversation repository — AI agents must not reopen.
			suite.mockConversationMessageRepository.EXPECT().Save(gomock.Any()).
				DoAndReturn(func(m *models.ConversationMessage) (*models.ConversationMessage, error) { return m, nil }).
				Times(1)
		}

		t.Run("Should ingest the message without reopening the conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := &requests.IngestMessageRequest{
				ConversationID: conversation.ID,
				ParticipantID:  aiAgent.ID,
				Content:        "final word",
			}
			initSuite(suite, request)

			message, err := suite.sut.Ingest(request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, message.Content)
		})
	})

	t.Run("Given conversation and participant exist", func(t *testing.T) {
		t.Parallel()

		conversation := model_builders.NewConversationBuilder().Build()
		participant := model_builders.NewAgentParticipantBuilder().Build()

		initSuite := func(suite *WhenIngestingMessageTestingSuite, request *requests.IngestMessageRequest) {
			suite.mockConversationRepository.EXPECT().GetById(conversation.ID).Return(conversation, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().GetById(participant.ID).Return(participant, nil).Times(1)
			suite.mockConversationMessageRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(m *models.ConversationMessage) (*models.ConversationMessage, error) {
					assert.Equal(t, conversation.ID, m.ConversationID)
					assert.Equal(t, participant.ID, m.ParticipantID)
					assert.Equal(t, request.Content, m.Content)
					return m, nil
				}).
				Times(1)
		}

		t.Run("Should persist and return the message", func(t *testing.T) {
			t.Parallel()

			suite := WhenIngestingMessageBeforeEach(t)
			request := &requests.IngestMessageRequest{
				ConversationID: conversation.ID,
				ParticipantID:  participant.ID,
				Content:        "hi",
			}
			initSuite(suite, request)

			message, err := suite.sut.Ingest(request)

			assert.NoError(t, err)
			assert.Equal(t, request.Content, message.Content)
		})
	})
}
