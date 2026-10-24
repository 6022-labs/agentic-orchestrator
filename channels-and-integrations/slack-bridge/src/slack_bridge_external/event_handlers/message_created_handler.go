package event_handlers

import (
	"context"
	"fmt"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	command_handlers_requests "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	external_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	NO_FACILITATOR_FOUND_ERROR = "Error: No facilitator found in this channel. " +
		"Please add a facilitator to manage conversations."
)

const (
	AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE = "agent_message_created"
)

type MessageCreatedHandlerInterface interface {
	Handle(ctx context.Context, request *requests.MessageCreatedRequest) error
}

type MessageCreatedHandler struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	logger                         *zap.Logger
	slackSettings                  *settings.SlackSettings
	slackThreadRequester           interfaces.SlackThreadRequesterInterface
	slackAgentRepository           repositories.SlackAgentRepositoryInterface
	addAgentToChannel              command_handlers.AddAgentToChannelInterface
	brokerMessageForwarder         external_interfaces.BrokerMessageForwarderInterface
	slackUserSynchronizerCached    external_interfaces.SlackUserSynchronizerCachedInterface
	slackChannelSynchronizerCached external_interfaces.SlackChannelSynchronizerCachedInterface
}

func NewMessageCreatedHandler(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	logger *zap.Logger,
	slackSettings *settings.SlackSettings,
	addAgentToChannel command_handlers.AddAgentToChannelInterface,
	slackThreadRequester interfaces.SlackThreadRequesterInterface,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
	brokerMessageForwarder external_interfaces.BrokerMessageForwarderInterface,
	slackUserSynchronizerCached external_interfaces.SlackUserSynchronizerCachedInterface,
	slackChannelSynchronizerCached external_interfaces.SlackChannelSynchronizerCachedInterface,
) *MessageCreatedHandler {
	return &MessageCreatedHandler{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		logger:                         logger,
		slackSettings:                  slackSettings,
		addAgentToChannel:              addAgentToChannel,
		slackThreadRequester:           slackThreadRequester,
		slackAgentRepository:           slackAgentRepository,
		brokerMessageForwarder:         brokerMessageForwarder,
		slackUserSynchronizerCached:    slackUserSynchronizerCached,
		slackChannelSynchronizerCached: slackChannelSynchronizerCached,
	}
}

func (p *MessageCreatedHandler) Handle(ctx context.Context, request *requests.MessageCreatedRequest) (err error) {
	ctx, span := p.channelTracer.StartEventHandling(ctx, "message.created")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		if err != nil {
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		p.channelMetrics.RecordEvent(ctx, "message.created", outcome, time.Since(start))
		span.End()
	}()

	p.logger.Debug(
		"Handling message created event",
		zap.String("messageTs", request.Event.Ts),
		zap.String("threadTs", lo.FromPtrOr(request.Event.ThreadTs, "nil")),
		zap.String("channelId", request.Event.Channel),
		zap.String("channelType", request.Event.ChannelType),
	)

	slackChannel, err := p.slackChannelSynchronizerCached.Sync(ctx, request.Event.Channel)
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
			p.logger.Info(
				"Skipping no facilitator error message as it is already posted",
				zap.String("channelId", slackChannel.Id),
			)
			return nil
		}

		err = p.slackThreadRequester.PushCompletionToThread(ctx,
			NO_FACILITATOR_FOUND_ERROR,
			nil,
			request.Event.Channel,
			request.Event.Ts,
		)
		if err != nil {
			return fmt.Errorf("failed to push completion to thread: %w", err)
		}

		return nil
	}

	if request.Event.AppId != nil && *request.Event.AppId == p.slackSettings.AppId {
		return p.handleAsAgent(ctx, request, slackChannel)
	}

	if request.Event.User != nil {
		return p.handleAsSlackUser(ctx, request, slackChannel)
	}

	return fmt.Errorf("message created event does not have a user or app_id: %+v", request.Event)
}

func (p *MessageCreatedHandler) handleAsSlackUser(
	ctx context.Context,
	request *requests.MessageCreatedRequest,
	slackChannel *models.SlackChannel,
) error {
	slackUser, err := p.slackUserSynchronizerCached.Sync(ctx, *request.Event.User)
	if err != nil {
		return fmt.Errorf("failed to synchronize slack user and slack user not found in persistence: %w", err)
	}
	if slackUser == nil {
		return fmt.Errorf("slack user not found in persistence: %s", *request.Event.User)
	}

	if slackUser.PairedHumanAgentCollectionAddress != nil && slackUser.PairedHumanAgentCollectionTokenId != nil {
		return p.handleAsSlackUserAgent(ctx, request, slackUser, slackChannel)
	}

	parentMessageTs := p.getParentMessageTs(request)

	err = p.brokerMessageForwarder.Forward(ctx,
		slackChannel,
		parentMessageTs,
		request.Event.Ts,
		request.Event.Text,
		slackUser,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to forward slack user message to broker: %w", err)
	}

	return nil
}

func (p *MessageCreatedHandler) handleAsAgent(
	ctx context.Context,
	request *requests.MessageCreatedRequest,
	slackChannel *models.SlackChannel,
) error {
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

	agent, err := p.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		*collectionTokenAddress,
		*collectionTokenId,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get agent with metadata %s and %d : %w",
			*collectionTokenAddress,
			*collectionTokenId,
			err,
		)
	}
	if agent == nil {
		return fmt.Errorf(
			"agent not found in persistence with metadata %s and %d",
			*collectionTokenAddress,
			*collectionTokenId,
		)
	}

	parentMessageTs := p.getParentMessageTs(request)

	err = p.brokerMessageForwarder.Forward(ctx,
		slackChannel,
		parentMessageTs,
		request.Event.Ts,
		request.Event.Text,
		nil,
		agent,
	)
	if err != nil {
		return fmt.Errorf("failed to forward agent message to broker: %w", err)
	}

	return nil
}

func (p *MessageCreatedHandler) handleAsSlackUserAgent(
	ctx context.Context,
	request *requests.MessageCreatedRequest,
	slackUser *models.SlackUser,
	slackChannel *models.SlackChannel,
) error {
	pairedHumanAgent, err := p.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		*slackUser.PairedHumanAgentCollectionAddress,
		*slackUser.PairedHumanAgentCollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get paired human agent: %w", err)
	}
	if pairedHumanAgent == nil {
		return fmt.Errorf(
			"paired human agent not found in persistence with collection address and token ID: %s %d",
			*slackUser.PairedHumanAgentCollectionAddress,
			*slackUser.PairedHumanAgentCollectionTokenId,
		)
	}

	err = p.addAgentToChannel.Add(metrics_interfaces.WithSystemCommandTrigger(ctx),
		&command_handlers_requests.AddAgentToChannelRequest{
			ChannelID: slackChannel.Id,
			AgentCommandArguments: command_handlers_requests.AgentCommandArguments{
				CollectionAddress: *slackUser.PairedHumanAgentCollectionAddress,
				CollectionTokenId: *slackUser.PairedHumanAgentCollectionTokenId,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add paired human agent to channel: %w", err)
	}

	// Re-sync channel to pick up the freshly attached paired human agent so the
	// broker registration includes them in the conversation participant set.
	freshChannel, err := p.slackChannelSynchronizerCached.Sync(ctx, slackChannel.Id)
	if err != nil {
		return fmt.Errorf("failed to refresh slack channel after attaching paired human agent: %w", err)
	}
	if freshChannel != nil {
		slackChannel = freshChannel
	}

	parentMessageTs := p.getParentMessageTs(request)

	err = p.brokerMessageForwarder.Forward(ctx,
		slackChannel,
		parentMessageTs,
		request.Event.Ts,
		request.Event.Text,
		nil,
		pairedHumanAgent,
	)
	if err != nil {
		return fmt.Errorf("failed to forward paired-human-agent message to broker: %w", err)
	}

	return nil
}

func (p *MessageCreatedHandler) getParentMessageTs(request *requests.MessageCreatedRequest) string {
	if request.Event.ThreadTs != nil {
		return *request.Event.ThreadTs
	}

	return request.Event.Ts
}
