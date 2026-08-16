package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"go.uber.org/zap"
)

type ConversationProcessorDaemon struct {
	ctx          context.Context
	cancel       context.CancelFunc
	startOnce    sync.Once
	stopOnce     sync.Once
	processingWg sync.WaitGroup

	logger                 *zap.Logger
	processingSettings     *settings.ProcessingSettings
	conversationProcessor  interfaces.ConversationProcessorInterface
	participantRepository  repositories.ParticipantRepositoryInterface
	conversationRepository repositories.ConversationRepositoryInterface

	participantRoleResolver interfaces.ParticipantRoleResolverInterface
}

func NewConversationProcessorDaemon(
	logger *zap.Logger,
	processingSettings *settings.ProcessingSettings,
	conversationProcessor interfaces.ConversationProcessorInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
	conversationRepository repositories.ConversationRepositoryInterface,
	participantRoleResolver interfaces.ParticipantRoleResolverInterface,
) *ConversationProcessorDaemon {
	ctx, cancel := context.WithCancel(context.Background())

	return &ConversationProcessorDaemon{
		logger:                 logger,
		processingSettings:     processingSettings,
		conversationProcessor:  conversationProcessor,
		participantRepository:  participantRepository,
		conversationRepository: conversationRepository,

		participantRoleResolver: participantRoleResolver,

		ctx:    ctx,
		cancel: cancel,
	}
}

func (d *ConversationProcessorDaemon) Start() {
	d.startOnce.Do(func() {
		d.logger.Info("Starting conversation processor daemon")
		go d.runLoop()
	})
}

func (d *ConversationProcessorDaemon) Stop() {
	d.stopOnce.Do(func() {
		d.logger.Info("Stopping conversation processor daemon")
		d.cancel()
		d.processingWg.Wait()
		d.logger.Info("Conversation processor daemon stopped cleanly")
	})
}

// Tick runs exactly one daemon tick synchronously. Exposed for unit tests
// that want to drive the loop without spawning the ticker.
func (d *ConversationProcessorDaemon) Tick() {
	d.processNextParticipant()
}

// WaitForInFlight blocks until every goroutine spawned by Tick/Start has
// finished. Exposed for unit tests.
func (d *ConversationProcessorDaemon) WaitForInFlight() {
	d.processingWg.Wait()
}

func (d *ConversationProcessorDaemon) runLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Info("Shutting down conversation processor daemon")
			return
		case <-ticker.C:
			d.processNextParticipant()
		}
	}
}

func (d *ConversationProcessorDaemon) processNextParticipant() {
	participant, err := d.participantRepository.GetNextAgentToProcess(
		d.processingSettings.ParticipantProcessingInterval,
	)
	if err != nil {
		d.logger.Error("Failed to fetch next participant to process", zap.Error(err))
		return
	}
	if participant == nil {
		d.logger.Debug("No participants to process at the moment")
		return
	}

	if err := d.setParticipantIsProcessing(participant); err != nil {
		d.logger.Error("Failed to mark participant as processing",
			d.participantLogFields(participant, zap.Error(err))...,
		)
		return
	}

	d.processingWg.Add(1)
	go d.processConversationsForParticipant(participant)
}

func (d *ConversationProcessorDaemon) processConversationsForParticipant(participant *models.Participant) {
	defer func() {
		d.processingWg.Done()
		if err := d.setParticipantIsNotProcessing(participant); err != nil {
			d.logger.Error("Failed to mark participant as not processing",
				d.participantLogFields(participant, zap.Error(err))...,
			)
		}
	}()

	role, err := d.participantRoleResolver.Resolve(participant)
	if err != nil {
		d.logger.Error("Failed to resolve participant role",
			d.participantLogFields(participant, zap.Error(err))...,
		)
		return
	}

	conversations, err := d.conversationRepository.ListProcessableForParticipant(
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

			if err := d.conversationProcessor.ProcessConversation(participant, conv); err != nil {
				d.logger.Error(
					"Failed to process conversation for participant",
					d.participantLogFields(
						participant,
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

func (d *ConversationProcessorDaemon) setParticipantIsProcessing(participant *models.Participant) error {
	participant.IsProcessing = true

	if _, err := d.participantRepository.Save(participant); err != nil {
		return fmt.Errorf("failed to set participant %s as processing: %w", participant.ID, err)
	}
	return nil
}

func (d *ConversationProcessorDaemon) setParticipantIsNotProcessing(participant *models.Participant) error {
	participant.IsProcessing = false
	now := time.Now()
	participant.LastProcessedAt = &now

	if _, err := d.participantRepository.Save(participant); err != nil {
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
