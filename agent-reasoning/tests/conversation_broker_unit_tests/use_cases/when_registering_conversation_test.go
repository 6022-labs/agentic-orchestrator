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

type WhenRegisteringConversationTestingSuite struct {
	sut *use_cases.RegisterConversation

	mockBridgeRepository                  *repositories_mocks.MockBridgeRepositoryInterface
	mockConversationRepository            *repositories_mocks.MockConversationRepositoryInterface
	mockConversationParticipantRepository *repositories_mocks.MockConversationParticipantRepositoryInterface
	mockParticipantRepository             *repositories_mocks.MockParticipantRepositoryInterface
}

func WhenRegisteringConversationBeforeEach(t *testing.T) *WhenRegisteringConversationTestingSuite {
	mockController := gomock.NewController(t)

	mockBridgeRepository := repositories_mocks.NewMockBridgeRepositoryInterface(mockController)
	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockConversationParticipantRepository := repositories_mocks.NewMockConversationParticipantRepositoryInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	sut := use_cases.NewRegisterConversation(
		mockBridgeRepository,
		mockConversationRepository,
		mockConversationParticipantRepository,
		mockParticipantRepository,
	)

	return &WhenRegisteringConversationTestingSuite{
		sut: sut,

		mockBridgeRepository:                  mockBridgeRepository,
		mockConversationRepository:            mockConversationRepository,
		mockConversationParticipantRepository: mockConversationParticipantRepository,
		mockParticipantRepository:             mockParticipantRepository,
	}
}

func generateValidRegisterConversationRequest() *requests.RegisterConversationRequest {
	return &requests.RegisterConversationRequest{
		BridgeID: uuid.New(),
	}
}

func TestWhenRegisteringConversation(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because bridgeId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringConversationBeforeEach(t)
			request := generateValidRegisterConversationRequest()
			request.BridgeID = uuid.Nil

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "bridgeId")
		})

		t.Run("Because a participantId entry is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringConversationBeforeEach(t)
			request := generateValidRegisterConversationRequest()
			request.ParticipantIDs = []uuid.UUID{uuid.Nil}

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.ValidationError{}, err)
			assert.Contains(t, err.Error(), "participantIds")
		})
	})

	t.Run("Given the bridge does not exist", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRegisteringConversationTestingSuite, request *requests.RegisterConversationRequest) {
			suite.mockBridgeRepository.EXPECT().
				GetById(request.BridgeID).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringConversationBeforeEach(t)
			request := generateValidRegisterConversationRequest()
			initSuite(suite, request)

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given a participant does not exist", func(t *testing.T) {
		t.Parallel()

		bridge := model_builders.NewBridgeBuilder().Build()
		missingParticipantId := uuid.New()

		initSuite := func(suite *WhenRegisteringConversationTestingSuite) {
			suite.mockBridgeRepository.EXPECT().GetById(bridge.ID).Return(bridge, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().
				GetById(missingParticipantId).
				Return(nil, gorm.ErrRecordNotFound).
				Times(1)
		}

		t.Run("Should return a not found error before saving the conversation", func(t *testing.T) {
			t.Parallel()

			suite := WhenRegisteringConversationBeforeEach(t)
			request := &requests.RegisterConversationRequest{
				BridgeID:       bridge.ID,
				ParticipantIDs: []uuid.UUID{missingParticipantId},
			}
			initSuite(suite)

			_, err := suite.sut.Register(request)

			assert.IsType(t, &cb_errors.NotFoundError{}, err)
		})
	})

	t.Run("Given the bridge and participants exist", func(t *testing.T) {
		t.Parallel()

		bridge := model_builders.NewBridgeBuilder().Build()
		participant1 := model_builders.NewParticipantBuilder().Build()
		participant2 := model_builders.NewNonMintedParticipantBuilder().Build()

		initSuite := func(suite *WhenRegisteringConversationTestingSuite, request *requests.RegisterConversationRequest) {
			suite.mockBridgeRepository.EXPECT().GetById(bridge.ID).Return(bridge, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().GetById(participant1.ID).Return(participant1, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().GetById(participant2.ID).Return(participant2, nil).Times(1)

			suite.mockConversationRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(c *models.Conversation) (*models.Conversation, error) {
					assert.Equal(t, bridge.ID, c.BridgeID)
					assert.Equal(t, request.Description, c.Description)
					c.ID = uuid.New()
					return c, nil
				}).
				Times(1)

			suite.mockConversationParticipantRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(link *models.ConversationParticipant) (*models.ConversationParticipant, error) {
					assert.Contains(t, []uuid.UUID{participant1.ID, participant2.ID}, link.ParticipantID)
					return link, nil
				}).
				Times(2)
		}

		t.Run("Should create the conversation and attach all participants", func(t *testing.T) {
			t.Parallel()

			description := "channel description"
			suite := WhenRegisteringConversationBeforeEach(t)
			request := &requests.RegisterConversationRequest{
				BridgeID:       bridge.ID,
				Description:    &description,
				ParticipantIDs: []uuid.UUID{participant1.ID, participant2.ID},
			}
			initSuite(suite, request)

			conversation, err := suite.sut.Register(request)

			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, conversation.ID)
			assert.Equal(t, bridge.ID, conversation.BridgeID)
		})
	})
}
