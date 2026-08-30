package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenResolvingMemorySessionTestingSuite struct {
	sut *services.MemorySessionResolver

	mockConversationMemoryBindingRepository *repositories_mocks.MockConversationMemoryBindingRepositoryInterface
}

func WhenResolvingMemorySessionBeforeEach(t *testing.T) *WhenResolvingMemorySessionTestingSuite {
	mockController := gomock.NewController(t)

	mockConversationMemoryBindingRepository := repositories_mocks.NewMockConversationMemoryBindingRepositoryInterface(
		mockController,
	)

	return &WhenResolvingMemorySessionTestingSuite{
		sut: services.NewMemorySessionResolver(
			mockConversationMemoryBindingRepository,
		),
		mockConversationMemoryBindingRepository: mockConversationMemoryBindingRepository,
	}
}

func TestWhenResolvingMemorySession(t *testing.T) {
	t.Parallel()

	t.Run("Given an expert", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()
		sessionId := uuid.New()

		t.Run("Should return nothing until the agent has named a session", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), conversation.ID, participant.ID).
				Return(nil, nil)

			resolved, err := suite.sut.Resolve(
				context.Background(),
				lo.ToPtr(enums.ParticipantRoleExpert),
				conversation,
				participant,
				[]*models.Participant{participant},
			)

			assert.NoError(t, err)
			assert.Nil(t, resolved)
		})

		t.Run("Should return the bound session without deriving one", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), conversation.ID, participant.ID).
				Return(&models.ConversationMemoryBinding{MemorySessionID: sessionId}, nil)

			resolved, err := suite.sut.Resolve(
				context.Background(),
				lo.ToPtr(enums.ParticipantRoleExpert),
				conversation,
				participant,
				[]*models.Participant{participant},
			)

			assert.NoError(t, err)
			assert.Equal(t, sessionId, *resolved)
		})
	})

	t.Run("Given a facilitator", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()

		t.Run("Should bind the conversation itself on the first turn", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), conversation.ID, participant.ID).
				Return(nil, nil)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				Bind(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, binding *models.ConversationMemoryBinding) error {
					assert.Equal(t, conversation.ID, binding.MemorySessionID)
					return nil
				})

			resolved, err := suite.sut.Resolve(
				context.Background(),
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				conversation,
				participant,
				[]*models.Participant{participant},
			)

			assert.NoError(t, err)
			assert.Equal(t, conversation.ID, *resolved)
		})

		t.Run("Should reuse the stored binding without writing again", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), conversation.ID, participant.ID).
				Return(&models.ConversationMemoryBinding{MemorySessionID: conversation.ID}, nil)

			resolved, err := suite.sut.Resolve(
				context.Background(),
				lo.ToPtr(enums.ParticipantRoleFacilitator),
				conversation,
				participant,
				[]*models.Participant{participant},
			)

			assert.NoError(t, err)
			assert.Equal(t, conversation.ID, *resolved)
		})
	})

	t.Run("Given a clone", func(t *testing.T) {
		t.Parallel()

		clone := model_builders.NewAgentParticipantBuilder().Build()
		peer := model_builders.NewAgentParticipantBuilder().Build()
		human := model_builders.NewUserParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()

		resolveWith := func(
			suite *WhenResolvingMemorySessionTestingSuite,
			participants []*models.Participant,
		) (*uuid.UUID, error) {
			return suite.sut.Resolve(
				context.Background(),
				lo.ToPtr(enums.ParticipantRoleClone),
				conversation,
				clone,
				participants,
			)
		}

		t.Run("Should key the session on the agents present, ignoring humans and order", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(2)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				Bind(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(2)

			first, err := resolveWith(suite, []*models.Participant{clone, peer, human})
			assert.NoError(t, err)

			second, err := resolveWith(suite, []*models.Participant{peer, clone})
			assert.NoError(t, err)

			assert.Equal(t, *first, *second)
		})

		t.Run("Should move the binding when a participant joins", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			stale := uuid.New()
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), conversation.ID, clone.ID).
				Return(&models.ConversationMemoryBinding{MemorySessionID: stale}, nil)

			var rebound uuid.UUID
			suite.mockConversationMemoryBindingRepository.EXPECT().
				Bind(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, binding *models.ConversationMemoryBinding) error {
					rebound = binding.MemorySessionID
					return nil
				})

			resolved, err := resolveWith(suite, []*models.Participant{clone, peer})

			assert.NoError(t, err)
			assert.NotEqual(t, stale, *resolved)
			assert.Equal(t, rebound, *resolved)
		})

		t.Run("Should leave the binding alone while the set is unchanged", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil)
			var derived uuid.UUID
			suite.mockConversationMemoryBindingRepository.EXPECT().
				Bind(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, binding *models.ConversationMemoryBinding) error {
					derived = binding.MemorySessionID
					return nil
				})

			_, err := resolveWith(suite, []*models.Participant{clone, peer})
			assert.NoError(t, err)

			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&models.ConversationMemoryBinding{MemorySessionID: derived}, nil)

			resolved, err := resolveWith(suite, []*models.Participant{clone, peer})

			assert.NoError(t, err)
			assert.Equal(t, derived, *resolved)
		})
	})

	t.Run("Given a role that never runs a turn", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversation := model_builders.NewConversationBuilder().Build()

		t.Run("Should resolve nothing for a human", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingMemorySessionBeforeEach(t)
			suite.mockConversationMemoryBindingRepository.EXPECT().
				GetByConversationAndParticipant(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil)

			resolved, err := suite.sut.Resolve(
				context.Background(),
				lo.ToPtr(enums.ParticipantRoleHuman),
				conversation,
				participant,
				[]*models.Participant{participant},
			)

			assert.NoError(t, err)
			assert.Nil(t, resolved)
		})
	})
}
