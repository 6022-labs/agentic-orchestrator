package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/settings"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SlackThreadProcessor struct {
	logger                    *zap.Logger
	processingSettings        *settings.ProcessingSettings
	slackThreadConclude       interfaces.SlackThreadConcludeInterface
	slackThreadRequester      interfaces.SlackThreadRequesterInterface
	slackThreadRepository     repositories.SlackThreadRepositoryInterface
	slackMessageRepository    repositories.SlackMessageRepositoryInterface
	slackChannelRepository    repositories.SlackChannelRepositoryInterface
	agentCompletionsRequester interfaces.AgentCompletionsRequesterInterface
}

func NewSlackThreadProcessor(
	logger *zap.Logger,
	processingSettings *settings.ProcessingSettings,
	slackThreadConclude interfaces.SlackThreadConcludeInterface,
	slackThreadRequester interfaces.SlackThreadRequesterInterface,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
	slackMessageRepository repositories.SlackMessageRepositoryInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
	agentCompletionsRequester interfaces.AgentCompletionsRequesterInterface,
) *SlackThreadProcessor {
	return &SlackThreadProcessor{
		logger:                    logger,
		processingSettings:        processingSettings,
		slackThreadConclude:       slackThreadConclude,
		slackThreadRequester:      slackThreadRequester,
		slackThreadRepository:     slackThreadRepository,
		slackMessageRepository:    slackMessageRepository,
		slackChannelRepository:    slackChannelRepository,
		agentCompletionsRequester: agentCompletionsRequester,
	}
}

func (p *SlackThreadProcessor) ProcessThread(slackAgent *models.SlackAgent, slackThread *models.SlackThread) error {
	if p.shouldForceConcludeThread(slackThread) {
		p.logger.Warn("Force concluding thread",
			zap.String("threadId", slackThread.Id),
			zap.String("channelId", slackThread.ChannelId),
			zap.Uint64("agentCollectionTokenId", slackAgent.CollectionTokenId),
		)

		err := p.slackThreadConclude.Conclude(slackThread.Id, slackThread.ChannelId)
		if err != nil {
			return fmt.Errorf("failed to force conclude slack thread: %w", err)
		}

		return nil
	}

	slackChannel, err := p.slackChannelRepository.GetById(slackThread.ChannelId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to fetch slack channel: %w", err)
	}
	if slackChannel == nil {
		p.logger.Warn("Slack channel not found",
			zap.String("channelId", slackThread.ChannelId),
			zap.String("agentEnsDomain", slackAgent.EnsDomain),
			zap.Uint64("agentCollectionTokenId", slackAgent.CollectionTokenId),
		)
		slackThread.IsConcluded = true
		return nil
	}

	slackMessages, err := p.slackMessageRepository.GetProcessableMessagesForAgentAndThread(
		slackAgent,
		slackThread.Id,
		slackThread.ChannelId,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch slack messages for thread %s: %w", slackThread.Id, err)
	}
	if len(slackMessages) == 0 {
		p.logger.Info("No messages found in thread for agent",
			zap.String("threadId", slackThread.Id),
			zap.String("channelId", slackThread.ChannelId),
			zap.String("agentEnsDomain", slackAgent.EnsDomain),
			zap.Uint64("agentCollectionTokenId", slackAgent.CollectionTokenId),
		)

		return nil
	}

	forceCompletion := false
	if slackAgent.Role == enums.AgentRoleFacilitator && p.shouldForceConcludeByFacilitatorThread(slackThread) {
		forceCompletion = true
	}

	completions, isFinal, err := p.agentCompletionsRequester.Completions(
		slackChannel,
		slackThread,
		slackMessages,
		slackAgent.CollectionAddress,
		slackAgent.CollectionTokenId,
		forceCompletion,
	)
	if err != nil {
		return fmt.Errorf("failed to get completions: %v", err)
	}

	if slackAgent.Role == enums.AgentRoleHuman {
		// Stop here for human agents, no completions will be returned
		return nil
	}

	if len(completions) == 0 {
		return fmt.Errorf("no completions found for thread %v", slackThread.Id)
	}

	err = p.sendCompletionsToThread(completions, slackAgent, slackThread)
	if err != nil {
		return fmt.Errorf("failed to send completions to thread: %v", err)
	}

	if *isFinal {
		p.logger.Info("Thread is concluded", zap.String("threadId", slackThread.Id))

		err = p.slackThreadConclude.Conclude(slackThread.Id, slackThread.ChannelId)
		if err != nil {
			return fmt.Errorf("failed to conclude slack thread: %v", err)
		}
	}

	return nil
}

func (p *SlackThreadProcessor) sendCompletionsToThread(
	completions []*string,
	slackAgent *models.SlackAgent,
	slackThread *models.SlackThread,
) error {
	// Reload the thread to ensure it's not yet concluded
	reloadedSlackThread, err := p.slackThreadRepository.GetByIdAndChannelId(slackThread.Id, slackThread.ChannelId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to fetch slack thread: %w", err)
	}
	if reloadedSlackThread == nil {
		return fmt.Errorf("slack thread not found: %s %s", slackThread.Id, slackThread.ChannelId)
	}
	if reloadedSlackThread.IsConcluded {
		p.logger.Warn("Thread is already concluded, skipping completions",
			zap.String("threadId", slackThread.Id),
			zap.String("channelId", slackThread.ChannelId),
			zap.String("agentEnsDomain", slackAgent.EnsDomain),
			zap.Uint64("agentCollectionTokenId", slackAgent.CollectionTokenId),
		)

		return nil
	}

	p.logger.Info("Sending completions to thread",
		zap.String("threadId", slackThread.Id),
		zap.String("channelId", slackThread.ChannelId),
		zap.String("agentEnsDomain", slackAgent.EnsDomain),
		zap.Uint64("agentCollectionTokenId", slackAgent.CollectionTokenId),
	)

	for _, completion := range completions {
		err = p.slackThreadRequester.PushCompletionToThread(
			*completion,
			slackAgent,
			slackThread.ChannelId,
			slackThread.Id,
		)
		if err != nil {
			return fmt.Errorf("failed to push completion to thread: %v", err)
		}
	}

	return nil
}

func (p *SlackThreadProcessor) shouldForceConcludeByFacilitatorThread(slackThread *models.SlackThread) bool {
	return p.shouldForceConcludeThreadAfter(slackThread, p.processingSettings.ConversationForceConcludeByFacilitatorAfter)
}

func (p *SlackThreadProcessor) shouldForceConcludeThread(slackThread *models.SlackThread) bool {
	return p.shouldForceConcludeThreadAfter(slackThread, p.processingSettings.ConversationForceConcludeAfter)
}

func (p *SlackThreadProcessor) shouldForceConcludeThreadAfter(slackThread *models.SlackThread, delaySeconds uint64) bool {
	if delaySeconds == 0 {
		return false
	}

	concludeAfter := time.Duration(delaySeconds) * time.Second
	var referenceTime time.Time
	if slackThread.UnconcludedAt != nil {
		referenceTime = *slackThread.UnconcludedAt
	} else {
		referenceTime = slackThread.CreatedAt()
	}

	if referenceTime.IsZero() {
		return false
	}

	return referenceTime.Add(concludeAfter).Before(time.Now())
}
