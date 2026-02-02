package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_internal/settings"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SlackThreadProcessorDaemon struct {
	ctx                context.Context
	cancel             context.CancelFunc
	startOnce          sync.Once
	stopOnce           sync.Once
	processingAgents   map[string]bool
	processingAgentsWg sync.WaitGroup

	logger                *zap.Logger
	processingSettings    *settings.ProcessingSettings
	slackThreadProcessor  interfaces.SlackThreadProcessorInterface
	slackAgentRepository  repositories.SlackAgentRepositoryInterface
	slackThreadRepository repositories.SlackThreadRepositoryInterface
}

func NewSlackThreadProcessorDaemon(
	logger *zap.Logger,
	processingSettings *settings.ProcessingSettings,
	slackThreadProcessor interfaces.SlackThreadProcessorInterface,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
) *SlackThreadProcessorDaemon {
	ctx, cancel := context.WithCancel(context.Background())

	return &SlackThreadProcessorDaemon{
		logger:                logger,
		processingSettings:    processingSettings,
		slackThreadProcessor:  slackThreadProcessor,
		slackAgentRepository:  slackAgentRepository,
		slackThreadRepository: slackThreadRepository,

		ctx:              ctx,
		cancel:           cancel,
		processingAgents: make(map[string]bool),
	}
}

func (p *SlackThreadProcessorDaemon) Start() {
	p.logger.Info("Starting Slack thread processor")

	p.startOnce.Do(func() {
		p.logger.Info("Starting Slack thread processor daemon")
		go p.runLoop()
	})
}

func (p *SlackThreadProcessorDaemon) Stop() {
	p.logger.Info("Stopping Slack thread processor")

	p.stopOnce.Do(func() {
		p.logger.Info("Stopping Slack thread processor daemon")
		p.cancel()
		p.processingAgentsWg.Wait()
		p.logger.Info("Slack thread processor daemon stopped cleanly")
	})
}

func (p *SlackThreadProcessorDaemon) runLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("Shutting down Slack thread processor daemon")
			return
		case <-ticker.C:
			p.processNextAgent()
		}
	}
}

func (p *SlackThreadProcessorDaemon) processNextAgent() {
	agent, err := p.slackAgentRepository.GetNextAgentToProcess(p.processingSettings.AgentProcessingInterval)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		p.logger.Error("Failed to fetch next agent to process threads", zap.Error(err))
		return
	}
	if agent == nil {
		p.logger.Debug("No agents to process at the moment")
		return
	}

	err = p.setAgentIsProcessing(agent)
	if err != nil {
		p.logger.Error("Failed to mark agent as processing",
			zap.String("agentEnsDomain", agent.EnsDomain),
			zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
			zap.Error(err),
		)
		return
	}

	go p.processThreadsForAgent(agent)
}

func (p *SlackThreadProcessorDaemon) processThreadsForAgent(agent *models.SlackAgent) {
	p.processingAgentsWg.Add(1)
	defer func() {
		p.processingAgentsWg.Done()
		err := p.setAgentIsNotProcessing(agent)
		if err != nil {
			p.logger.Error("Failed to mark agent as not processing",
				zap.String("agentEnsDomain", agent.EnsDomain),
				zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
				zap.Error(err),
			)
		}
	}()

	threads, err := p.slackThreadRepository.GetProcessableThreadsForAgent(agent, p.processingSettings.FacilitatorTickInterval)
	if err != nil {
		p.logger.Error("Failed to fetch processable threads for agent",
			zap.String("agentEnsDomain", agent.EnsDomain),
			zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
			zap.Error(err),
		)
		return
	}
	if len(threads) == 0 {
		p.logger.Debug("No processable threads found for agent",
			zap.String("agentEnsDomain", agent.EnsDomain),
			zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
		)
		return
	}

	p.logger.Info("Found processable threads for agent",
		zap.Int("threadCount", len(threads)),
		zap.String("agentEnsDomain", agent.EnsDomain),
		zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
	)

	wg := &sync.WaitGroup{}
	wg.Add(len(threads))

	for _, slackThread := range threads {
		go func(thread *models.SlackThread) {
			defer wg.Done()

			p.logger.Info("Processing thread for agent",
				zap.String("threadId", thread.Id),
				zap.String("channelId", thread.ChannelId),
				zap.String("agentEnsDomain", agent.EnsDomain),
				zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
			)

			err := p.slackThreadProcessor.ProcessThread(agent, thread)
			if err != nil {
				p.logger.Error("Failed to process thread for agent",
					zap.String("threadId", thread.Id),
					zap.String("channelId", thread.ChannelId),
					zap.String("agentEnsDomain", agent.EnsDomain),
					zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
					zap.Error(err),
				)
				return
			}

			p.logger.Info("Successfully processed thread for agent",
				zap.String("threadId", thread.Id),
				zap.String("channelId", thread.ChannelId),
				zap.String("agentEnsDomain", agent.EnsDomain),
				zap.Uint64("agentCollectionTokenId", agent.CollectionTokenId),
			)
		}(slackThread)
	}

	wg.Wait()
}

func (p *SlackThreadProcessorDaemon) setAgentIsProcessing(agent *models.SlackAgent) error {
	agent.IsProcessing = true
	p.processingAgents[agent.EnsDomain] = true

	_, err := p.slackAgentRepository.Save(agent)
	if err != nil {
		return fmt.Errorf("failed to set agent %s as processing: %w", agent.EnsDomain, err)
	}

	return nil
}

func (p *SlackThreadProcessorDaemon) setAgentIsNotProcessing(agent *models.SlackAgent) error {
	agent.IsProcessing = false
	now := time.Now()
	agent.LastProcessedAt = &now

	delete(p.processingAgents, agent.EnsDomain)

	_, err := p.slackAgentRepository.Save(agent)
	if err != nil {
		return fmt.Errorf("failed to set agent %s as not processing: %w", agent.EnsDomain, err)
	}

	return nil
}
