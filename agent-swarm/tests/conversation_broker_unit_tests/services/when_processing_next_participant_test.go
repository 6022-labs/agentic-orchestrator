package services_test

import (
	"context"

	"errors"
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

type WhenProcessingNextParticipantTestingSuite struct {
	mockBrokerMetrics *broker_metrics_mocks.MockBrokerMetricsInterface
	sut               *services.ConversationProcessorDaemon

	mockConversationProcessor   *interfaces_mocks.MockConversationProcessorInterface
	mockParticipantRepository   *repositories_mocks.MockParticipantRepositoryInterface
	mockConversationRepository  *repositories_mocks.MockConversationRepositoryInterface
	mockParticipantRoleResolver *interfaces_mocks.MockParticipantRoleResolverInterface
}

func WhenProcessingNextParticipantBeforeEach(
	t *testing.T,
	leaseDuration time.Duration,
) *WhenProcessingNextParticipantTestingSuite {
	mockController := gomock.NewController(t)

	mockBrokerMetrics := broker_metrics_mocks.NewMockBrokerMetricsInterface(mockController)

	mockConversationProcessor := interfaces_mocks.NewMockConversationProcessorInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)
	mockConversationRepository := repositories_mocks.NewMockConversationRepositoryInterface(mockController)
	mockParticipantRoleResolver := interfaces_mocks.NewMockParticipantRoleResolverInterface(mockController)

	processingSettings := &settings.ProcessingSettings{
		ParticipantProcessingInterval:               20,
		FacilitatorTickInterval:                     30,
		ConversationForceConcludeByFacilitatorAfter: 300,
		ConversationForceConcludeAfter:              600,
		LeaseDuration:                               leaseDuration,
		MaxConcurrentParticipants:                   1,
		MaxConcurrentConversations:                  2,
	}

	sut := services.NewConversationProcessorDaemon(
		zap.NewNop(),
		mockBrokerMetrics,
		processingSettings,
		mockConversationProcessor,
		mockParticipantRepository,
		mockConversationRepository,
		mockParticipantRoleResolver,
	)

	return &WhenProcessingNextParticipantTestingSuite{
		mockBrokerMetrics: mockBrokerMetrics,
		sut:               sut,

		mockConversationProcessor:  mockConversationProcessor,
		mockParticipantRepository:  mockParticipantRepository,
		mockConversationRepository: mockConversationRepository,

		mockParticipantRoleResolver: mockParticipantRoleResolver,
	}

}

func TestWhenProcessingNextParticipant(t *testing.T) {
	t.Parallel()

	t.Run("Given no participant is eligible", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), uint64(20), time.Minute).
				Return(nil, nil).
				Times(1)
		}

		t.Run("Should do nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			initSuite(suite)

			// No panics, no further interactions — mock controller enforces it.
			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeIdle)

			suite.sut.Tick(context.Background())
		})
	})

	t.Run("Given the next participant has no processable conversations", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), uint64(20), time.Minute).
				Return(participant, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(gomock.Any(), participant, gomock.Any(), uint64(30)).
				Return(nil, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ReleaseLease(gomock.Any(), participant).Return(nil).Times(1)
		}

		t.Run("Should find nothing, then release it", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			initSuite(suite)

			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeProcessed)

			suite.sut.Tick(context.Background())
			suite.sut.WaitForInFlight()
		})
	})

	t.Run("Given the next participant has processable conversations", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conv1 := model_builders.NewConversationBuilder().Build()
		conv2 := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), uint64(20), time.Minute).
				Return(participant, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(gomock.Any(), participant, gomock.Any(), uint64(30)).
				Return([]*models.Conversation{conv1, conv2}, nil).
				Times(1)
			suite.mockConversationProcessor.EXPECT().
				ProcessConversation(gomock.Any(), participant, conv1).
				Return(nil).
				Times(1)
			suite.mockConversationProcessor.EXPECT().
				ProcessConversation(gomock.Any(), participant, conv2).
				Return(nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ReleaseLease(gomock.Any(), participant).Return(nil).Times(1)
		}

		t.Run("Should dispatch each conversation to the processor", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			initSuite(suite)

			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeProcessed)

			suite.sut.Tick(context.Background())
			suite.sut.WaitForInFlight()
		})
	})

	t.Run("Given one conversation processor call fails", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conv := model_builders.NewConversationBuilder().Build()

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), uint64(20), time.Minute).
				Return(participant, nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(gomock.Any(), participant, gomock.Any(), uint64(30)).
				Return([]*models.Conversation{conv}, nil).
				Times(1)
			suite.mockConversationProcessor.EXPECT().
				ProcessConversation(gomock.Any(), participant, conv).
				Return(errors.New("boom")).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ReleaseLease(gomock.Any(), participant).Return(nil).Times(1)
		}

		t.Run("Should still release the participant cleanly", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			initSuite(suite)

			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeProcessed)

			suite.sut.Tick(context.Background())
			suite.sut.WaitForInFlight()
		})
	})

	t.Run("Given a conversation outlasting its lease", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conv := model_builders.NewConversationBuilder().Build()
		leaseDuration := 300 * time.Millisecond

		initSuite := func(
			suite *WhenProcessingNextParticipantTestingSuite,
			process func(ctx context.Context, _ *models.Participant, _ *models.Conversation) error,
		) {
			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeProcessed)
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), leaseDuration).
				Return(participant, nil).
				Times(1)
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(gomock.Any(), participant, gomock.Any(), gomock.Any()).
				Return([]*models.Conversation{conv}, nil).
				Times(1)
			suite.mockConversationProcessor.EXPECT().
				ProcessConversation(gomock.Any(), participant, conv).
				DoAndReturn(process).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ReleaseLease(gomock.Any(), participant).Return(nil).Times(1)
		}
		untilCancelled := func(ctx context.Context, _ *models.Participant, _ *models.Conversation) error {
			<-ctx.Done()
			return ctx.Err()
		}

		t.Run("Should keep renewing it", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, leaseDuration)
			initSuite(suite, func(ctx context.Context, _ *models.Participant, _ *models.Conversation) error {
				time.Sleep(2 * leaseDuration)
				assert.NoError(t, ctx.Err())
				return nil
			})
			suite.mockParticipantRepository.EXPECT().
				RenewLease(gomock.Any(), participant, leaseDuration).
				Return(true, nil).
				MinTimes(3)

			suite.sut.Tick(context.Background())
			suite.sut.WaitForInFlight()
		})

		t.Run("Should cancel it once the lease is lost", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, leaseDuration)
			initSuite(suite, untilCancelled)
			suite.mockParticipantRepository.EXPECT().
				RenewLease(gomock.Any(), participant, leaseDuration).
				Return(false, nil).
				Times(1)

			suite.sut.Tick(context.Background())
			suite.sut.WaitForInFlight()
		})

		t.Run("Should cancel it once renewals failed for the whole lease", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, leaseDuration)
			initSuite(suite, untilCancelled)
			suite.mockParticipantRepository.EXPECT().
				RenewLease(gomock.Any(), participant, leaseDuration).
				Return(false, errors.New("db down")).
				MinTimes(2)

			suite.sut.Tick(context.Background())
			suite.sut.WaitForInFlight()
		})
	})

	t.Run("Given as many participants in flight as allowed", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()

		t.Run("Should claim another one only once a slot is free", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			unblock := make(chan struct{})
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(participant, nil).
				Times(1)
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				DoAndReturn(func(context.Context, *models.Participant) (*enums.ParticipantRole, error) {
					<-unblock
					return nil, errors.New("boom")
				}).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ReleaseLease(gomock.Any(), participant).Return(nil).Times(1)
			gomock.InOrder(
				suite.mockBrokerMetrics.EXPECT().
					RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeProcessed),
				suite.mockBrokerMetrics.EXPECT().
					RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeSaturated),
				suite.mockBrokerMetrics.EXPECT().
					RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeIdle),
			)

			suite.sut.Tick(context.Background())
			suite.sut.Tick(context.Background())
			close(unblock)
			suite.sut.WaitForInFlight()

			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil).
				Times(1)
			suite.sut.Tick(context.Background())
		})
	})

	t.Run("Given claims that start nothing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should free their slot", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			gomock.InOrder(
				suite.mockParticipantRepository.EXPECT().
					ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db down")),
				suite.mockParticipantRepository.EXPECT().
					ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(2),
			)
			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeFailed)
			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeIdle).
				Times(2)

			for range 3 {
				suite.sut.Tick(context.Background())
			}
		})
	})

	t.Run("Given more conversations than one participant may answer at once", func(t *testing.T) {
		t.Parallel()

		participant := model_builders.NewAgentParticipantBuilder().Build()
		conversations := lo.Times(6, func(int) *models.Conversation {
			return model_builders.NewConversationBuilder().Build()
		})

		initSuite := func(suite *WhenProcessingNextParticipantTestingSuite) {
			suite.mockBrokerMetrics.EXPECT().
				RecordDaemonTick(gomock.Any(), metrics_interfaces.DaemonTickOutcomeProcessed)
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(participant, nil).
				Times(1)
			suite.mockParticipantRoleResolver.EXPECT().
				Resolve(gomock.Any(), participant).
				Return(lo.ToPtr(enums.ParticipantRoleClone), nil).
				Times(1)
			suite.mockConversationRepository.EXPECT().
				ListProcessableForParticipant(gomock.Any(), participant, gomock.Any(), gomock.Any()).
				Return(conversations, nil).
				Times(1)
			suite.mockParticipantRepository.EXPECT().ReleaseLease(gomock.Any(), participant).Return(nil).Times(1)
		}

		t.Run("Should process them all without exceeding the limit", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			initSuite(suite)
			var running, peak atomic.Int32
			suite.mockConversationProcessor.EXPECT().
				ProcessConversation(gomock.Any(), participant, gomock.Any()).
				DoAndReturn(func(context.Context, *models.Participant, *models.Conversation) error {
					current := running.Add(1)
					defer running.Add(-1)
					for seen := peak.Load(); current > seen && !peak.CompareAndSwap(seen, current); seen = peak.Load() {
					}
					time.Sleep(10 * time.Millisecond)
					return nil
				}).
				Times(len(conversations))

			suite.sut.Tick(context.Background())
			suite.sut.WaitForInFlight()

			assert.LessOrEqual(t, peak.Load(), int32(2))
		})

		t.Run("Should not start the waiting ones once cancelled", func(t *testing.T) {
			t.Parallel()

			suite := WhenProcessingNextParticipantBeforeEach(t, time.Minute)
			initSuite(suite)
			ctx, cancel := context.WithCancel(context.Background())
			entered := make(chan struct{}, len(conversations))
			suite.mockConversationProcessor.EXPECT().
				ProcessConversation(gomock.Any(), participant, gomock.Any()).
				DoAndReturn(func(ctx context.Context, _ *models.Participant, _ *models.Conversation) error {
					entered <- struct{}{}
					<-ctx.Done()
					time.Sleep(50 * time.Millisecond)
					return ctx.Err()
				}).
				Times(2)

			suite.sut.Tick(ctx)
			<-entered
			<-entered
			cancel()
			suite.sut.WaitForInFlight()
		})
	})
}
