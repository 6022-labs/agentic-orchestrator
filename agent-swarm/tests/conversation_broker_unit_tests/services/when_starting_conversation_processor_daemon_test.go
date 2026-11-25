package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	broker_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenStartingConversationProcessorDaemonTestingSuite struct {
	mockBrokerMetrics *broker_metrics_mocks.MockBrokerMetricsInterface
	sut               *services.ConversationProcessorDaemon

	mockParticipantRepository *repositories_mocks.MockParticipantRepositoryInterface
}

func WhenStartingConversationProcessorDaemonBeforeEach(
	t *testing.T,
) *WhenStartingConversationProcessorDaemonTestingSuite {
	mockController := gomock.NewController(t)

	mockBrokerMetrics := broker_metrics_mocks.NewMockBrokerMetricsInterface(mockController)
	mockParticipantRepository := repositories_mocks.NewMockParticipantRepositoryInterface(mockController)

	return &WhenStartingConversationProcessorDaemonTestingSuite{
		mockBrokerMetrics: mockBrokerMetrics,
		sut: services.NewConversationProcessorDaemon(
			zap.NewNop(),
			mockBrokerMetrics,
			&settings.ProcessingSettings{TickInterval: 10 * time.Millisecond},
			interfaces_mocks.NewMockConversationProcessorInterface(mockController),
			mockParticipantRepository,
			repositories_mocks.NewMockConversationRepositoryInterface(mockController),
			interfaces_mocks.NewMockParticipantRoleResolverInterface(mockController),
		),

		mockParticipantRepository: mockParticipantRepository,
	}
}

func TestWhenStartingConversationProcessorDaemon(t *testing.T) {
	t.Parallel()

	t.Run("Given a tick interval", func(t *testing.T) {
		t.Parallel()

		t.Run("Should tick at that interval", func(t *testing.T) {
			t.Parallel()

			suite := WhenStartingConversationProcessorDaemonBeforeEach(t)
			ticks := make(chan struct{}, 100)
			suite.mockBrokerMetrics.EXPECT().RecordDaemonTick(gomock.Any(), gomock.Any()).AnyTimes()
			suite.mockParticipantRepository.EXPECT().
				ClaimNextAgentToProcess(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(context.Context, uint64, time.Duration) (*models.Participant, error) {
					ticks <- struct{}{}
					return nil, nil
				}).
				MinTimes(3)

			started := time.Now()
			suite.sut.Start(context.Background())
			for range 3 {
				<-ticks
			}
			suite.sut.Stop()

			assert.Less(t, time.Since(started), time.Second)
		})
	})
}
