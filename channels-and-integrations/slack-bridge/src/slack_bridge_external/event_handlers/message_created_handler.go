package event_handlers

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	command_handlers_requests "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	external_interfaces "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	NO_FACILITATOR_FOUND_ERROR = "Error: No facilitator found in this channel. Please add a facilitator to manage conversations."
)

const (
	AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE = "agent_message_created"
)

type MessageCreatedHandlerInterface interface {
	Handle(request *requests.MessageCreatedRequest) error
}

type MessageCreatedHandler struct {
	logger                         *zap.Logger
	slackSettings                  *settings.SlackSettings
	slackThreadRequester           interfaces.SlackThreadRequesterInterface
	slackAgentRepository           repositories.SlackAgentRepositoryInterface
	slackThreadRepository          repositories.SlackThreadRepositoryInterface
	addAgentToChannel              command_handlers.AddAgentToChannelInterface
	slackThreadMessagesAdder       external_interfaces.SlackThreadMessagesAdderInterface
	slackUserSynchronizerCached    external_interfaces.SlackUserSynchronizerCachedInterface
	slackChannelSynchronizerCached external_interfaces.SlackChannelSynchronizerCachedInterface
}

func NewMessageCreatedHandler(
	logger *zap.Logger,
	slackSettings *settings.SlackSettings,
	addAgentToChannel command_handlers.AddAgentToChannelInterface,
	slackThreadRequester interfaces.SlackThreadRequesterInterface,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
	slackThreadRepository repositories.SlackThreadRepositoryInterface,
	slackThreadMessagesAdder external_interfaces.SlackThreadMessagesAdderInterface,
	slackUserSynchronizerCached external_interfaces.SlackUserSynchronizerCachedInterface,
	slackChannelSynchronizerCached external_interfaces.SlackChannelSynchronizerCachedInterface,
) *MessageCreatedHandler {
	return &MessageCreatedHandler{
		logger:                         logger,
		slackSettings:                  slackSettings,
		addAgentToChannel:              addAgentToChannel,
		slackThreadRequester:           slackThreadRequester,
		slackAgentRepository:           slackAgentRepository,
		slackThreadRepository:          slackThreadRepository,
		slackThreadMessagesAdder:       slackThreadMessagesAdder,
		slackUserSynchronizerCached:    slackUserSynchronizerCached,
		slackChannelSynchronizerCached: slackChannelSynchronizerCached,
	}
}

func (p *MessageCreatedHandler) Handle(request *requests.MessageCreatedRequest) error {
	p.logger.Debug(
		"Handling message created event",
		zap.String("messageTs", request.Event.Ts),
		zap.String("threadTs", lo.FromPtrOr(request.Event.ThreadTs, "nil")),
		zap.String("channelId", request.Event.Channel),
		zap.String("channelType", request.Event.ChannelType),
	)

	slackChannel, err := p.slackChannelSynchronizerCached.Sync(request.Event.Channel)
	if err != nil {
		return fmt.Errorf("failed to synchronize slack channel and slack channel not found in persistence: %w", err)
	}
	if slackChannel == nil {
		return fmt.Errorf("slack channel not found in persistence: %s", request.Event.Channel)
	}

	containFacilitator := lo.ContainsBy(slackChannel.Agents, func(agent models.SlackAgent) bool {
		return agent.Role == enums.AgentRoleFacilitator
	})
	if !containFacilitator {
		p.logger.Info("No facilitator found in channel", zap.String("channelId", slackChannel.Id))

		p.logger.Debug("Checking if the error message is already posted",
			zap.String("channelId", slackChannel.Id),
			zap.String("ts", request.Event.Ts),
			zap.String("eventText", request.Event.Text),
		)

		if request.ApiAppId == p.slackSettings.AppId && request.Event.Text == NO_FACILITATOR_FOUND_ERROR {
			p.logger.Info("Skipping no facilitator error message as it is already posted", zap.String("channelId", slackChannel.Id))
			return nil
		}

		err = p.slackThreadRequester.PushCompletionToThread(NO_FACILITATOR_FOUND_ERROR, nil, request.Event.Channel, request.Event.Ts)
		if err != nil {
			return fmt.Errorf("failed to push completion to thread: %w", err)
		}

		return nil
	}

	if request.Event.AppId != nil && *request.Event.AppId == p.slackSettings.AppId {
		return p.handleAsAgent(request, slackChannel)
	}

	if request.Event.User != nil {
		return p.handleAsSlackUser(request, slackChannel)
	}

	return fmt.Errorf("message created event does not have a user or app_id: %+v", request.Event)
}

func (p *MessageCreatedHandler) handleAsSlackUser(request *requests.MessageCreatedRequest, slackChannel *models.SlackChannel) error {
	slackUser, err := p.slackUserSynchronizerCached.Sync(*request.Event.User)
	if err != nil {
		return fmt.Errorf("failed to synchronize slack user and slack user not found in persistence: %w", err)
	}
	if slackUser == nil {
		return fmt.Errorf("slack user not found in persistence: %s", *request.Event.User)
	}

	if slackUser.PairedHumanAgentCollectionAddress != nil && slackUser.PairedHumanAgentCollectionTokenId != nil {
		return p.handleAsSlackUserAgent(request, slackUser, slackChannel)
	}

	parentMessageTs := p.getParentMessageTs(request)

	err = p.slackThreadMessagesAdder.Add(
		slackChannel.Id,
		parentMessageTs,
		request.Event.Ts,
		request.Event.Text,
		request.Event.User,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to add message to thread: %w", err)
	}

	err = p.unconcludeThreadIfNeeded(parentMessageTs, slackChannel.Id)
	if err != nil {
		return fmt.Errorf("failed to unconclude thread: %w", err)
	}

	return nil
}

func (p *MessageCreatedHandler) handleAsAgent(request *requests.MessageCreatedRequest, slackChannel *models.SlackChannel) error {
	if request.Event.Metadata == nil {
		return fmt.Errorf("message created event from agent does not have metadata: %+v", request.Event)
	}

	if request.Event.Metadata.EventType != AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE {
		return fmt.Errorf("unsupported agent message created event type: %s", request.Event.Metadata.EventType)
	}

	collectionTokenAddress, collectionTokenId, err := request.Event.Metadata.ToAgentIdentifiers()
	if err != nil {
		return err
	}

	agent, err := p.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(
		*collectionTokenAddress,
		*collectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get agent with metadata %s and %d : %w", *collectionTokenAddress, *collectionTokenId, err)
	}
	if agent == nil {
		return fmt.Errorf("agent not found in persistence with metadata %s and %d", *collectionTokenAddress, *collectionTokenId)
	}

	parentMessageTs := p.getParentMessageTs(request)

	err = p.slackThreadMessagesAdder.Add(
		slackChannel.Id,
		parentMessageTs,
		request.Event.Ts,
		request.Event.Text,
		nil,
		&agent.CollectionAddress,
		&agent.CollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to add message to thread: %w", err)
	}

	return nil
}

func (p *MessageCreatedHandler) handleAsSlackUserAgent(request *requests.MessageCreatedRequest, slackUser *models.SlackUser, slackChannel *models.SlackChannel) error {
	pairedHumanAgent, err := p.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(
		*slackUser.PairedHumanAgentCollectionAddress,
		*slackUser.PairedHumanAgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get paired human agent: %w", err)
	}
	if pairedHumanAgent == nil {
		return fmt.Errorf("paired human agent not found in persistence with collection address and token ID: %s %d", *slackUser.PairedHumanAgentCollectionAddress, *slackUser.PairedHumanAgentCollectionTokenId)
	}

	parentMessageTs := p.getParentMessageTs(request)

	err = p.slackThreadMessagesAdder.Add(
		slackChannel.Id,
		parentMessageTs,
		request.Event.Ts,
		request.Event.Text,
		request.Event.User,
		&pairedHumanAgent.CollectionAddress,
		&pairedHumanAgent.CollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to add message to thread: %w", err)
	}

	err = p.addAgentToChannel.Add(&command_handlers_requests.AddAgentToChannelRequest{
		CollectionAddress: *slackUser.PairedHumanAgentCollectionAddress,
		CollectionTokenId: *slackUser.PairedHumanAgentCollectionTokenId,
		ChannelID:         slackChannel.Id,
	})
	if err != nil {
		return fmt.Errorf("failed to add paired human agent to channel: %w", err)
	}

	err = p.unconcludeThreadIfNeeded(parentMessageTs, slackChannel.Id)
	if err != nil {
		return fmt.Errorf("failed to unconclude thread: %w", err)
	}

	return nil
}

func (p *MessageCreatedHandler) getParentMessageTs(request *requests.MessageCreatedRequest) string {
	if request.Event.ThreadTs != nil {
		return *request.Event.ThreadTs
	}

	return request.Event.Ts
}

func (p *MessageCreatedHandler) unconcludeThreadIfNeeded(slackThreadId, slackChannelId string) error {
	slackThread, err := p.slackThreadRepository.GetByIdAndChannelId(slackThreadId, slackChannelId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get slack thread by id and channel id: %w", err)
	}
	if slackThread == nil {
		return fmt.Errorf("slack thread not found after adding message to thread: %s in channel %s", slackThreadId, slackChannelId)
	}

	if slackThread.IsConcluded {
		p.logger.Info("Reopening concluded thread due to new message from slack user", zap.String("threadId", slackThread.Id), zap.String("channelId", slackChannelId))

		now := time.Now()
		slackThread.IsConcluded = false
		slackThread.UnconcludedAt = &now

		_, err = p.slackThreadRepository.Save(slackThread)
		if err != nil {
			return fmt.Errorf("failed to reopen concluded slack thread: %w", err)
		}
	}

	return nil
}
