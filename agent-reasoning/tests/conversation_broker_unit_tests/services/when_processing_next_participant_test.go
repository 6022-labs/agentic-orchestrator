package services_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenProcessingNextParticipantTestingSuite struct {
	sut *services.ConversationProcessorDaemon

	mockConversationProcessor  *interfaces_mocks.MockConversationProcessorInterface
	mockParticipantRepository  *repositories_mocks.MockParticipantRepositoryInterface
	mockConversationRepository *repositories_mocks.MockConversationRepositoryInterface
}

func WhenProcessingNextParticipantBeforeEach(t *testing.T) *WhenProcessingNextParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationProcessor := interfaces_mocks.NewMockConversationProcessorInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)
	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)

	processingSettings := &settings.ProcessingSettings{
		ParticipantProcessingInterval:               20,
		FacilitatorTickInterval:                     30,
		ConversationForceConcludeByFacilitatorAfter: 300,
		ConversationForceConcludeAfter:              600,
	}

	sut := services.NewConversationProcessorDaemon(
		zap.NewNop(),
		processingSettings,
		mockConversationProcessor,
		mockParticipantRepository,
		mockConversationRepository,
	)

	return &WhenProcessingNextParticipantTestingSuite{
		sut: sut,

		mockConversationProcessor:  mockConversationProcessor,
		mockParticipantRepository:  mockParticipantRepository,
		mockConversationRepository: mockConversationRepository,
	}
}

func TestWhenProcessingNextParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given no participant is eligible", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRepository.EXPECT().
				GetNextAgentToProcess(gomock.Any()).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should do nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t)
			initSuite(suite)

			// No panics, no further interactions — mock controller enforces it.
			suite.sut.Tick()
		})
	})

	t.Run("Given the next participant has no processable conversations", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRepository.EXPECT().
				GetNextAgentToProcess(gomock.Any()).
				Return(participant, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(p *models.Participant) (*models.Participant, error) {
					assert.True(t, p.IsProcessing)
					return p, nil
				}).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(participant, uint64(30)).
				Return(nil, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				Save(gomock.Any()).
				DoAndReturn(func(p *models.Participant) (*models.Participant, error) {
					assert.False(t, p.IsProcessing)
					assert.NotNil(t, p.LastProcessedAt)
					return p, nil
				}).
				Times(1)
		}

		t.Run("Should mark it processing, find nothing, then release it", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t)
			initSuite(suite)

			suite.sut.Tick()
			suite.sut.WaitForInFlight()
		})
	})

	t.Run("Given the next participant has processable conversations", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conv1 := model_builders.NewConversationBuilder().Build()
		conv2 := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRepository.EXPECT().GetNextAgentToProcess(gomock.Any()).Return(participant, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().Save(gomock.Any()).Return(participant, nil).Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(participant, uint64(30)).
				Return([]*models.Conversation{conv1, conv2}, nil).
				Times(1)
			suite.mockConversationProcessor.EXPECT().ProcessConversation(participant, conv1).Return(nil).Times(1)
			suite.mockConversationProcessor.EXPECT().ProcessConversation(participant, conv2).Return(nil).Times(1)
			suite.mockParticipantRepository.EXPECT().Save(gomock.Any()).Return(participant, nil).Times(1)
		}

		t.Run("Should dispatch each conversation to the processor", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t)
			initSuite(suite)

			suite.sut.Tick()
			suite.sut.WaitForInFlight()
		})
	})

	t.Run("Given one conversation processor call fails", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conv := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRepository.EXPECT().GetNextAgentToProcess(gomock.Any()).Return(participant, nil).Times(1)
			suite.mockParticipantRepository.EXPECT().Save(gomock.Any()).Return(participant, nil).Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(participant, uint64(30)).
				Return([]*models.Conversation{conv}, nil).
				Times(1)
			suite.mockConversationProcessor.EXPECT().ProcessConversation(participant, conv).Return(errors.New("boom")).Times(1)
			suite.mockParticipantRepository.EXPECT().Save(gomock.Any()).Return(participant, nil).Times(1)
		}

		t.Run("Should still release the participant cleanly", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t)
			initSuite(suite)

			suite.sut.Tick()
			suite.sut.WaitForInFlight()
		})
	})
}
