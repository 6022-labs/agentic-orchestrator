package services_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	broker_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenStoppingConversationProcessorDaemonTestingSuite struct {
	mockBrokerMetrics *broker_metrics_mocks.MockBrokerMetricsInterface
	sut               *services.ConversationProcessorDaemon

	mockConversationProcessor   *interfaces_mocks.MockConversationProcessorInterface
	mockParticipantRepository   *repositories_mocks.MockParticipantRepositoryInterface
	mockConversationRepository  *repositories_mocks.MockConversationRepositoryInterface
	mockParticipantRoleResolver *interfaces_mocks.MockParticipantRoleResolverInterface
}

func WhenStoppingConversationProcessorDaemonBeforeEach(
	t *testing.T,
) *WhenStoppingConversationProcessorDaemonTestingSuite {
	mockController := gomock.NewController(t)

	mockBrokerMetrics := broker_metrics_mocks.NewMockBrokerMetricsInterface(mockController)
	mockConversationProcessor := interfaces_mocks.NewMockConversationProcessorInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)
	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockParticipantRoleResolver := interfaces_mocks.NewMockParticipantRoleResolverInterface(mockController)

	return &WhenStoppingConversationProcessorDaemonTestingSuite{
		mockBrokerMetrics: mockBrokerMetrics,
		sut: services.NewConversationProcessorDaemon(
			zap.NewNop(),
			mockBrokerMetrics,
			&settings.ProcessingSettings{
				ParticipantProcessingInterval: 20,
				FacilitatorTickInterval:       30,
				TickInterval:                  10 * time.Millisecond,
				LeaseDuration:                 time.Minute,
			},
			mockConversationProcessor,
			mockParticipantRepository,
			mockConversationRepository,
			mockParticipantRoleResolver,
		),

		mockConversationProcessor:   mockConversationProcessor,
		mockParticipantRepository:   mockParticipantRepository,
		mockConversationRepository:  mockConversationRepository,
		mockParticipantRoleResolver: mockParticipantRoleResolver,
	}
}

func TestWhenStoppingConversationProcessorDaemon(t *testing.T) {
	t.Parallel()

	t.Run("Given a daemon never started", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return at once", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoppingConversationProcessorDaemonBeforeEach(t)

			assert.NotPanics(t, suite.sut.Stop)
		})
	})

	t.Run("Given a conversation in flight", func(t *testing.T) {
		t.Parallel()

		t.Run("Should release its participant before returning", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoppingConversationProcessorDaemonBeforeEach(t)
			participant := model_builders.NewAgentParticipantBuilder().Build()
			conversation := model_builders.NewConversationBuilder().Build()
			entered := make(chan struct{})
			released := &atomic.Bool{}

			suite.mockBrokerMetrics.EXPECT().RecordDaemonTick(gomock.Any(), gomock.Any()).AnyTimes()
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(participant, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(gomock.Any(), participant, gomock.Any(), gomock.Any()).
				Return([]*models.Conversation{conversation}, nil).
				Times(1)
			suite.mockConversationProcessor.EXPECT().
				ProcessConversation(gomock.Any(), participant, conversation).
				DoAndReturn(func(ctx context.Context, _ *models.Participant, _ *models.Conversation) error {
					close(entered)
					<-ctx.Done()
					return ctx.Err()
				}).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ReleaseLease(gomock.Any(), participant).
				DoAndReturn(func(ctx context.Context, _ *models.Participant) error {
					assert.NoError(t, ctx.Err())
					_, bounded := ctx.Deadline()
					assert.True(t, bounded)
					time.Sleep(20 * time.Millisecond)
					released.Store(true)
					return nil
				}).
				Times(1)

			suite.sut.Start(context.Background())
			<-entered
			suite.sut.Stop()

			assert.True(t, released.Load())
		})
	})

	t.Run("Given a tick in flight", func(t *testing.T) {
		t.Parallel()

		t.Run("Should wait for it", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoppingConversationProcessorDaemonBeforeEach(t)
			entered := make(chan struct{})
			finished := &atomic.Bool{}

			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(context.Context, uint64, time.Duration) (*models.Participant, error) {
					close(entered)
					time.Sleep(20 * time.Millisecond)
					return nil, nil
				}).
				Times(1)
			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeIdle).
				Do(func(context.Context, string) { finished.Store(true) }).
				Times(1)

			suite.sut.Start(context.Background())
			<-entered
			suite.sut.Stop()

			assert.True(t, finished.Load())
		})
	})
}
