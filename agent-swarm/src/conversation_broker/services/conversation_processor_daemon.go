package services

import (
	"context"
	"sync"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"go.uber.org/zap"
)

const (
	participantReleaseTimeout = 5 * time.Second
	// Two renewals may fail before the lease is lost.
	leaseRenewalsPerDuration = 3
)

type ConversationProcessorDaemon struct {
	ctx          context.Context
	cancel       context.CancelFunc
	startOnce    sync.Once
	stopOnce     sync.Once
	loopWg       sync.WaitGroup
	processingWg sync.WaitGroup
	// slots holds one token per participant in flight; one is only claimed for a free slot.
	slots chan struct{}

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
		slots: make(chan struct{}, max(processingSettings.MaxConcurrentParticipants, 1)),

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
	ticker := time.NewTicker(d.processingSettings.TickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Info("Shutting down conversation processor daemon")
			return
		case <-ticker.C:
			// select picks at random when Stop and a tick are both ready; a stopped daemon never ticks again.
			if ctx.Err() == nil {
				d.processNextParticipant(ctx)
			}
		}
	}
}

func (d *ConversationProcessorDaemon) processNextParticipant(ctx context.Context) {
	select {
	case d.slots <- struct{}{}:
	default:
		d.logger.Debug("Too many participants in flight")
		d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeSaturated)
		return
	}

	// The database starts the lease when it runs the claim, so the local clock starts it no later.
	claimedAt := time.Now()
	participant, err := d.participantRepository.ClaimNextAgentToProcess(ctx,
		d.processingSettings.ParticipantProcessingInterval,
		d.processingSettings.LeaseDuration,
	)
	if err != nil {
		<-d.slots
		d.logger.Error("Failed to claim next participant to process", zap.Error(err))
		d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeFailed)
		return
	}
	if participant == nil {
		<-d.slots
		d.logger.Debug("No participants to process at the moment")
		d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeIdle)
		return
	}

	d.brokerMetrics.RecordDaemonTick(ctx, metrics_interfaces.DaemonTickOutcomeProcessed)

	d.processingWg.Add(1)
	go d.processConversationsForParticipant(ctx, participant, claimedAt)
}

func (d *ConversationProcessorDaemon) processConversationsForParticipant(
	ctx context.Context,
	participant *models.Participant,
	claimedAt time.Time,
) {
	defer d.processingWg.Done()
	defer func() { <-d.slots }()
	defer func() {
		// Stop cancels ctx; the release must still reach the database or the participant waits for its lease to expire.
		releaseCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), participantReleaseTimeout)
		defer cancel()
		if err := d.participantRepository.ReleaseLease(releaseCtx, participant); err != nil {
			d.logger.Error("Failed to release participant",
				d.participantLogFields(participant, zap.Error(err))...,
			)
		}
	}()

	leaseWg := &sync.WaitGroup{}
	defer leaseWg.Wait()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	leaseWg.Go(func() { d.keepLease(ctx, cancel, participant, claimedAt) })

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
	defer wg.Wait()
	slots := make(chan struct{}, max(d.processingSettings.MaxConcurrentConversations, 1))

	for _, conv := range conversations {
		select {
		case <-ctx.Done():
			return
		case slots <- struct{}{}:
		}

		wg.Go(func() {
			defer func() { <-slots }()

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
		})
	}
}

// keepLease renews the lease until ctx is done; losing it cancels the work, which another replica may redo.
func (d *ConversationProcessorDaemon) keepLease(
	ctx context.Context,
	cancel context.CancelFunc,
	participant *models.Participant,
	claimedAt time.Time,
) {
	leaseDuration := d.processingSettings.LeaseDuration
	loseLease := func() {
		d.logger.Warn("Lost participant lease", d.participantLogFields(participant)...)
		cancel()
	}
	expiry := time.AfterFunc(leaseDuration-time.Since(claimedAt), loseLease)
	defer expiry.Stop()
	ticker := time.NewTicker(leaseDuration / leaseRenewalsPerDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		renewedAt := time.Now()
		held, err := d.participantRepository.RenewLease(ctx, participant, leaseDuration)
		switch {
		case ctx.Err() != nil:
			return
		case err != nil:
			d.logger.Warn("Failed to renew participant lease", d.participantLogFields(participant, zap.Error(err))...)
		case !held:
			loseLease()
			return
		default:
			expiry.Reset(leaseDuration - time.Since(renewedAt))
		}
	}
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
