package services

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SlackThreadMessagesAdder struct {
	logger                 *zap.Logger
	slackThreadRepository  repositories.SlackThreadRepositoryInterface
	slackChannelRepository repositories.SlackChannelRepositoryInterface
	slackMessageRepository repositories.SlackMessageRepositoryInterface
}

func NewSlackThreadMessagesAdder(
	logger *zap.Logger,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
	slackMessageRepository repositories.SlackMessageRepositoryInterface,
	slackChannelRepository repositories.SlackChannelRepositoryInterface,
) *SlackThreadMessagesAdder {
	return &SlackThreadMessagesAdder{
		logger:                 logger,
		slackThreadRepository:  slackThreadRepository,
		slackMessageRepository: slackMessageRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func (s *SlackThreadMessagesAdder) Add(
	channelId, parentMessageTs, messageId, text string,
	authorUserId, authorCollectionAddress *string, authorCollectionTokenId *uint64,
) error {
	if authorUserId == nil && (authorCollectionAddress == nil || authorCollectionTokenId == nil) {
		return fmt.Errorf("either authorUserId or both authorCollectionAddress and authorCollectionTokenId must be provided")
	}

	slackChannel, err := s.slackChannelRepository.GetById(channelId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get slack channel: %w", err)
	}
	if slackChannel == nil {
		return fmt.Errorf("slack channel with ID %s not found", channelId)
	}

	slackThread, err := s.slackThreadRepository.GetByIdAndChannelId(parentMessageTs, channelId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get slack thread: %w", err)
	}
	if slackThread == nil {
		slackThread = &models.SlackThread{
			Id:          parentMessageTs,
			ChannelId:   channelId,
			IsConcluded: false,
		}

		_, err = s.slackThreadRepository.Save(slackThread)
		if err != nil {
			return fmt.Errorf("failed to save slack thread: %w", err)
		}

		s.logger.Info("Created new slack thread", zap.String("thread_id", slackThread.Id))
	}

	slackMessage := &models.SlackMessage{
		Id:      messageId,
		Content: text,

		ChannelId: channelId,
		ThreadId:  parentMessageTs,

		IsFromUser:  lo.Ternary(authorUserId != nil, true, false),
		IsFromAgent: lo.Ternary(authorCollectionAddress != nil && authorCollectionTokenId != nil, true, false),

		AuthorUserId:            authorUserId,
		AuthorCollectionAddress: authorCollectionAddress,
		AuthorCollectionTokenId: authorCollectionTokenId,
	}

	_, err = s.slackMessageRepository.Save(slackMessage)
	if err != nil {
		return fmt.Errorf("failed to save slack message: %w", err)
	}

	s.logger.Info(
		"Added slack message to thread",
		zap.String("thread_id", parentMessageTs),
		zap.String("message_id", messageId),
	)

	return nil
}
