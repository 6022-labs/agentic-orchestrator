package services

import (
	"context"

	"fmt"
	"slices"
	"strings"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	slack_bridge_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	broker_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// BrokerMessageForwarder translates a Slack message into the broker: dedup,
// register-on-first-sight of the conversation, then ingest.
type BrokerMessageForwarder struct {
	logger                      *zap.Logger
	channelMetrics              metrics_interfaces.ChannelMetricsInterface
	brokerSettings              *broker_settings.ConversationBrokerSettings
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	slackMessageRepository      repositories.SlackMessageRepositoryInterface
}

func NewBrokerMessageForwarder(
	logger *zap.Logger,
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	brokerSettings *broker_settings.ConversationBrokerSettings,
	conversationBrokerRequester slack_bridge_interfaces.ConversationBrokerRequesterInterface,
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	slackMessageRepository repositories.SlackMessageRepositoryInterface,
) *BrokerMessageForwarder {
	return &BrokerMessageForwarder{
		logger:                      logger,
		channelMetrics:              channelMetrics,
		brokerSettings:              brokerSettings,
		conversationBrokerRequester: conversationBrokerRequester,
		slackConversationRepository: slackConversationRepository,
		slackMessageRepository:      slackMessageRepository,
	}
}

// Forward translates a Slack message into the broker. Exactly one of
// authorSlackUser or authorSlackAgent must be non-nil, and its BrokerParticipantId
// must already be populated by its synchronizer.
func (f *BrokerMessageForwarder) Forward(
	ctx context.Context,
	slackChannel *models.SlackChannel,
	parentMessageTs, messageTs, content string,
	authorSlackUser *models.SlackUser,
	authorSlackAgent *models.SlackAgent,
) (err error) {
	// Only a real forward is counted: a deduplicated redelivery never reaches the broker.
	forwarded := false
	defer func() {
		if forwarded || err != nil {
			f.channelMetrics.RecordMessageForwarded(ctx, err)
		}
	}()

	alreadyForwarded, err := f.alreadyForwarded(ctx, slackChannel.Id, messageTs)
	if err != nil {
		return err
	}
	if alreadyForwarded {
		f.logger.Debug("Skipping slack message already forwarded to broker",
			zap.String("channelId", slackChannel.Id),
			zap.String("messageTs", messageTs),
		)
		return nil
	}

	authorParticipantId, err := f.resolveAuthorParticipant(authorSlackUser, authorSlackAgent)
	if err != nil {
		return err
	}

	conversationId, err := f.resolveOrRegisterConversation(ctx, slackChannel, parentMessageTs, *authorParticipantId)
	if err != nil {
		return err
	}

	brokerMessageId, err := f.conversationBrokerRequester.IngestMessage(
		ctx,
		*conversationId,
		*authorParticipantId,
		content,
	)
	if err != nil {
		return fmt.Errorf("failed to ingest message into broker: %w", err)
	}

	_, err = f.slackMessageRepository.Save(ctx,
		&models.SlackMessage{
			SlackChannelId:  slackChannel.Id,
			SlackMessageTs:  messageTs,
			BrokerMessageID: brokerMessageId,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to save slack message mapping: %w", err)
	}

	forwarded = true

	return nil
}

func (f *BrokerMessageForwarder) alreadyForwarded(ctx context.Context, channelId, messageTs string) (bool, error) {
	message, err := f.slackMessageRepository.GetBySlackIds(ctx, channelId, messageTs)
	if err != nil {
		return false, fmt.Errorf("failed to lookup slack message mapping: %w", err)
	}
	return message != nil, nil
}

func (f *BrokerMessageForwarder) resolveAuthorParticipant(
	authorSlackUser *models.SlackUser,
	authorSlackAgent *models.SlackAgent,
) (*uuid.UUID, error) {
	if authorSlackAgent != nil {
		if authorSlackAgent.BrokerParticipantId == nil {
			return nil, fmt.Errorf(
				"slack agent %s/%d has no broker participant id; SlackAgentSynchronizer must run before forwarding",
				authorSlackAgent.CollectionAddress,
				authorSlackAgent.CollectionTokenId,
			)
		}
		return authorSlackAgent.BrokerParticipantId, nil
	}
	if authorSlackUser == nil {
		return nil, fmt.Errorf("forward called with neither agent nor slack user")
	}
	if authorSlackUser.BrokerParticipantId == nil {
		return nil, fmt.Errorf(
			"slack user %s has no broker participant id; SlackUserSynchronizer must run before forwarding",
			authorSlackUser.Id,
		)
	}
	return authorSlackUser.BrokerParticipantId, nil
}

func (f *BrokerMessageForwarder) resolveOrRegisterConversation(
	ctx context.Context,
	slackChannel *models.SlackChannel,
	parentMessageTs string,
	authorParticipantId uuid.UUID,
) (*uuid.UUID, error) {
	conversation, err := f.slackConversationRepository.GetBySlackIds(ctx, slackChannel.Id, parentMessageTs)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup slack conversation mapping: %w", err)
	}
	if conversation != nil {
		return &conversation.BrokerConversationID, nil
	}

	participantIds := make([]uuid.UUID, 0, len(slackChannel.Agents)+1)
	for _, agent := range slackChannel.Agents {
		if agent.BrokerParticipantId == nil {
			return nil, fmt.Errorf(
				"channel agent %s/%d has no broker participant id; SlackAgentSynchronizer must run before forwarding",
				agent.CollectionAddress,
				agent.CollectionTokenId,
			)
		}
		participantIds = append(participantIds, *agent.BrokerParticipantId)
	}
	if !slices.Contains(participantIds, authorParticipantId) {
		participantIds = append(participantIds, authorParticipantId)
	}

	var description *string
	if d := strings.TrimSpace(slackChannel.Description); d != "" {
		description = &slackChannel.Description
	}

	conversationId, err := f.conversationBrokerRequester.RegisterConversation(ctx,
		f.brokerSettings.BridgeID,
		description,
		participantIds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register conversation with broker: %w", err)
	}

	_, err = f.slackConversationRepository.Save(ctx,
		&models.SlackConversation{
			SlackChannelId:       slackChannel.Id,
			SlackThreadTs:        parentMessageTs,
			BrokerConversationID: conversationId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save slack conversation mapping: %w", err)
	}

	f.logger.Info("Registered new broker conversation",
		zap.String("channelId", slackChannel.Id),
		zap.String("threadTs", parentMessageTs),
		zap.String("brokerConversationId", conversationId.String()),
		zap.Int("attachedParticipants", len(participantIds)),
	)

	return &conversationId, nil
}
