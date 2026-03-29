package services

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SlackThreadConclude struct {
	logger                *zap.Logger
	slackThreadRepository repositories.SlackThreadRepositoryInterface
}

func NewSlackThreadConclude(
	logger *zap.Logger,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
) *SlackThreadConclude {
	return &SlackThreadConclude{
		logger:                logger,
		slackThreadRepository: slackThreadRepository,
	}
}

func (s *SlackThreadConclude) concludeThread(slackThread *models.SlackThread) error {
	if slackThread.IsConcluded {
		s.logger.Info(
			"Thread is already concluded, skipping save",
			zap.String("threadId", slackThread.Id),
		)
		return nil
	}

	slackThread.IsConcluded = true
	_, err := s.slackThreadRepository.Save(slackThread)
	if err != nil {
		return fmt.Errorf("failed to save slack thread state: %w", err)
	}

	s.logger.Info("Successfully saved thread state",
		zap.String("threadId", slackThread.Id),
		zap.Bool("isConcluded", slackThread.IsConcluded),
	)

	return nil
}

func (s *SlackThreadConclude) Conclude(threadId, channelId string) error {
	slackThread, err := s.slackThreadRepository.GetByIdAndChannelId(threadId, channelId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to fetch slack thread from database: %w", err)
	}
	if slackThread == nil {
		return fmt.Errorf("slack thread not found in database: %s", threadId)
	}

	return s.concludeThread(slackThread)
}

func (s *SlackThreadConclude) ConcludeIfExists(threadId, channelId string) error {
	slackThread, err := s.slackThreadRepository.GetByIdAndChannelId(threadId, channelId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to fetch slack thread from database: %w", err)
	}
	if slackThread == nil {
		return nil
	}

	return s.concludeThread(slackThread)
}
