package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SlackUserSynchronizer struct {
	logger              *zap.Logger
	slackUserRequester  interfaces.SlackUserRequesterInterface
	slackUserRepository repositories.SlackUserRepositoryInterface
}

func NewSlackUserSynchronizer(
	logger *zap.Logger,
	slackUserRequester interfaces.SlackUserRequesterInterface,
	slackUserRepository repositories.SlackUserRepositoryInterface,
) *SlackUserSynchronizer {
	return &SlackUserSynchronizer{
		logger:              logger,
		slackUserRequester:  slackUserRequester,
		slackUserRepository: slackUserRepository,
	}
}

func (s *SlackUserSynchronizer) Sync(id string) (*models.SlackUser, error) {
	s.logger.Debug("Starting user synchronization", zap.String("id", id))

	existingUserInDatabase, err := s.slackUserRepository.GetById(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUserInDatabase == nil {
		existingUserInDatabase = &models.SlackUser{
			Id: id,
		}
	}

	existingUserInSlack, err := s.slackUserRequester.GetUserInformation(id)
	if err != nil {
		return nil, err
	}
	if existingUserInSlack == nil {
		return nil, fmt.Errorf("user with ID %s not found in Slack", id)
	}

	existingUserInDatabase.Name = existingUserInSlack.Name
	existingUserInDatabase.RealName = existingUserInSlack.RealName

	existingUserInDatabase, err = s.slackUserRepository.Save(existingUserInDatabase)
	if err != nil {
		return nil, err
	}

	return existingUserInDatabase, nil
}
