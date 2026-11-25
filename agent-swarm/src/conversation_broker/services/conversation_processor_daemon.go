package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"go.uber.org/zap"
)

const participantReleaseTimeout = 5 * time.Second

type ConversationProcessorDaemon struct {
	ctx          context.Context
	cancel       context.CancelFunc
	startOnce    sync.Once
	stopOnce     sync.Once
	loopWg       sync.WaitGroup
	processingWg sync.WaitGroup

	logger                 *zap.Logger
	brokerMetrics          metrics_interfaces.BrokerMetricsInterface
	processingSettings     *settings.ProcessingSettings
	conversationProcessor  interfaces.ConversationProcessorInterface
	participantRepository  repositories.ParticipantRepositoryInterface
	conversationRepository repositories.ConversationRepositoryInterface

	participantRoleResolver interfaces.ParticipantRoleResolverInterface
}

func NewConversationProcessorDaemon(
	logger *zap.Logger,
	brokerMetrics metrics_interfaces.BrokerMetricsInterface,
	processingSettings *settings.ProcessingSettings,
	conversationProcessor interfaces.ConversationProcessorInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
	conversationRepository repositories.ConversationRepositoryInterface,
	participantRoleResolver interfaces.ParticipantRoleResolverInterface,
) *ConversationProcessorDaemon {
	return &ConversationProcessorDaemon{
		logger:                 logger,
		brokerMetrics:          brokerMetrics,
		processingSettings:     processingSettings,
		conversationProcessor:  conversationProcessor,
		participantRepository:  participantRepository,
		conversationRepository: conversationRepository,

		participantRoleResolver: participantRoleResolver,
	}
}

func (d *ConversationProcessorDaemon) Start(ctx context.Context) {
	d.startOnce.Do(func() {
		// Cancellable by the caller's ctx or by Stop().
		d.ctx, d.cancel = context.WithCancel(ctx)

		d.logger.Info("Starting conversation processor daemon")
		d.loopWg.Add(1)
		go d.runLoop(d.ctx)
	})
}

func (d *ConversationProcessorDaemon) Stop() {
	d.stopOnce.Do(func() {
		d.logger.Info("Stopping conversation processor daemon")
		if d.cancel != nil {
			d.cancel()
		}
		d.loopWg.Wait()
		d.processingWg.Wait()
		d.logger.Info("Conversation processor daemon stopped cleanly")
	})
}

// Tick runs exactly one daemon tick synchronously. Exposed for unit tests
// that want to drive the loop without spawning the ticker.
func (d *ConversationProcessorDaemon) Tick(ctx context.Context) {
	d.processNextParticipant(ctx)
}

// WaitForInFlight blocks until every goroutine spawned by Tick/Start has
// finished. Exposed for unit tests.
func (d *ConversationProcessorDaemon) WaitForInFlight() {
	d.processingWg.Wait()
}

func (d *ConversationProcessorDaemon) runLoop(ctx context.Context) {
	defer d.loopWg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Info("Shutting down conversation processor daemon")
			return
		case <-ticker.C:
			d.processNextParticipant(ctx)
		}
	}
}

func (d *ConversationProcessorDaemon) processNextParticipant(ctx context.Context) {
	participant, err := d.participantRepository.GetNextAgentToProcess(ctx,
		d.processingSettings.ParticipantProcessingInterval,
	)
	if err != nil {
		d.logger.Error("Failed to fetch next participant to process", zap.Error(err))
		d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeFailed)
		return
	}
	if participant == nil {
		d.logger.Debug("No participants to process at the moment")
		d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeIdle)
		return
	}

	if err := d.setParticipantIsProcessing(ctx, participant); err != nil {
		d.logger.Error("Failed to mark participant as processing",
			d.participantLogFields(participant, zap.Error(err))...,
		)
		d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeFailed)
		return
	}

	d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeProcessed)

	d.processingWg.Add(1)
	go d.processConversationsForParticipant(ctx, participant)
}

func (d *ConversationProcessorDaemon) processConversationsForParticipant(
	ctx context.Context,
	participant *models.Participant,
) {
	defer d.processingWg.Done()
	defer func() {
		// Stop cancels ctx; the release must still reach the database or the participant is never picked again.
		releaseCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), participantReleaseTimeout)
		defer cancel()
		if err := d.setParticipantIsNotProcessing(releaseCtx, participant); err != nil {
			d.logger.Error("Failed to mark participant as not processing",
				d.participantLogFields(participant, zap.Error(err))...,
			)
		}
	}()

	role, err := d.participantRoleResolver.Resolve(ctx, participant)
	if err != nil {
		d.logger.Error("Failed to resolve participant role",
			d.participantLogFields(participant, zap.Error(err))...,
		)
		return
	}

	conversations, err := d.conversationRepository.ListProcessableForParticipant(ctx,
		participant,
		role,
		d.processingSettings.FacilitatorTickInterval,
	)
	if err != nil {
		d.logger.Error("Failed to fetch processable conversations for participant",
			d.participantLogFields(participant, zap.Error(err))...,
		)
		return
	}
	if len(conversations) == 0 {
		d.logger.Debug("No processable conversations found for participant",
			d.participantLogFields(participant)...,
		)
		return
	}

	d.logger.Info("Found processable conversations for participant",
		d.participantLogFields(participant, zap.Int("conversationCount", len(conversations)))...,
	)

	wg := &sync.WaitGroup{}
	wg.Add(len(conversations))

	for _, conversation := range conversations {
		go func(conv *models.Conversation) {
			defer wg.Done()

			d.logger.Info("Processing conversation for participant",
				d.participantLogFields(participant, zap.String("conversationId", conv.ID.String()))...,
			)

			if err := d.conversationProcessor.ProcessConversation(ctx, participant, conv); err != nil {
				d.logger.Error(
					"Failed to process conversation for participant",
					d.participantLogFields(participant,
						zap.String("conversationId", conv.ID.String()),
						zap.Error(err),
					)...,
				)
				return
			}

			d.logger.Info("Successfully processed conversation for participant",
				d.participantLogFields(participant, zap.String("conversationId", conv.ID.String()))...,
			)
		}(conversation)
	}

	wg.Wait()
}

func (d *ConversationProcessorDaemon) setParticipantIsProcessing(
	ctx context.Context,
	participant *models.Participant,
) error {
	participant.IsProcessing = true

	if _, err := d.participantRepository.Save(ctx, participant); err != nil {
		return fmt.Errorf("failed to set participant %s as processing: %w", participant.ID, err)
	}
	return nil
}

func (d *ConversationProcessorDaemon) setParticipantIsNotProcessing(
	ctx context.Context,
	participant *models.Participant,
) error {
	participant.IsProcessing = false
	now := time.Now()
	participant.LastProcessedAt = &now

	if _, err := d.participantRepository.Save(ctx, participant); err != nil {
		return fmt.Errorf("failed to set participant %s as not processing: %w", participant.ID, err)
	}
	return nil
}

func (d *ConversationProcessorDaemon) participantLogFields(
	participant *models.Participant,
	extra ...zap.Field,
) []zap.Field {
	fields := []zap.Field{zap.String("participantId", participant.ID.String())}
	if participant.EnsDomain != nil {
		fields = append(fields, zap.String("participantEnsDomain", *participant.EnsDomain))
	}
	if participant.CollectionAddress != nil {
		fields = append(fields, zap.String("participantCollectionAddress", *participant.CollectionAddress))
	}
	if participant.CollectionTokenId != nil {
		fields = append(fields, zap.Uint64("participantCollectionTokenId", *participant.CollectionTokenId))
	}
	return append(fields, extra...)
}
