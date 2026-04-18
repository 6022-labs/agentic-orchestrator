package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SlackUserSynchronizer struct {
	logger                      *zap.Logger
	slackUserRequester          interfaces.SlackUserRequesterInterface
	slackUserRepository         repositories.SlackUserRepositoryInterface
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
}

func NewSlackUserSynchronizer(
	logger *zap.Logger,
	slackUserRequester interfaces.SlackUserRequesterInterface,
	slackUserRepository repositories.SlackUserRepositoryInterface,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
) *SlackUserSynchronizer {
	return &SlackUserSynchronizer{
		logger:                      logger,
		slackUserRequester:          slackUserRequester,
		slackUserRepository:         slackUserRepository,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (s *SlackUserSynchronizer) Sync(id string) (*models.SlackUser, error) {
	s.logger.Debug("Starting user synchronization", zap.String("id", id))

	existingUserInDatabase, err := s.slackUserRepository.GetById(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existingUserInSlack, err := s.slackUserRequester.GetUserInformation(id)
	if err != nil {
		return nil, err
	}
	if existingUserInSlack == nil {
		return nil, fmt.Errorf("user with ID %s not found in Slack", id)
	}

	if existingUserInDatabase == nil {
		existingUserInDatabase = &models.SlackUser{
			Id: id,
		}
	}

	existingUserInDatabase.Name = existingUserInSlack.Name
	existingUserInDatabase.RealName = existingUserInSlack.RealName

	if existingUserInDatabase.BrokerParticipantId == nil {
		brokerParticipantId, err := s.conversationBrokerRequester.RegisterUser(s.pickDisplayName(existingUserInDatabase))
		if err != nil {
			return nil, fmt.Errorf("failed to register user participant with broker: %w", err)
		}
		existingUserInDatabase.BrokerParticipantId = &brokerParticipantId
	}

	saved, err := s.slackUserRepository.Save(existingUserInDatabase)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (s *SlackUserSynchronizer) pickDisplayName(slackUser *models.SlackUser) string {
	if name := strings.TrimSpace(slackUser.RealName); name != "" {
		return name
	}
	if name := strings.TrimSpace(slackUser.Name); name != "" {
		return name
	}
	return slackUser.Id
}
