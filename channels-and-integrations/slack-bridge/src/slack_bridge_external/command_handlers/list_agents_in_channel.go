package command_handlers

import (
	"context"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/responses"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
)

type ListAgentInChannelInterface interface {
	List(
		ctx context.Context,
		request *requests.ListAgentInChannelRequest,
	) ([]*responses.ListAgentsInChannelResponse, error)
}

type ListAgentInChannel struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface
}

func NewListAgentInChannel(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	slackChannelSynchronizer interfaces.SlackChannelSynchronizerInterface,
) *ListAgentInChannel {
	return &ListAgentInChannel{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		slackChannelSynchronizer: slackChannelSynchronizer,
	}
}

func (handler *ListAgentInChannel) List(
	ctx context.Context,
	request *requests.ListAgentInChannelRequest,
) (commandResult []*responses.ListAgentsInChannelResponse, err error) {
	ctx, span := handler.channelTracer.StartCommandHandling(ctx, "agent.list")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		if err != nil {
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		handler.channelMetrics.RecordCommand(
			ctx,
			"agent.list",
			metrics_interfaces.CommandTriggerFrom(ctx),
			outcome,
			time.Since(start),
		)
		span.End()
	}()

	slackChannel, err := handler.slackChannelSynchronizer.Sync(ctx, request.ChannelID)
	if err != nil {
		return nil, err
	}

	result := make([]*responses.ListAgentsInChannelResponse, len(slackChannel.Agents))
	for i, agent := range slackChannel.Agents {
		result[i] = &responses.ListAgentsInChannelResponse{
			Role:              agent.Role,
			EnsDomain:         agent.EnsDomain,
			CollectionAddress: agent.CollectionAddress,
			CollectionTokenId: agent.CollectionTokenId,
		}
	}

	return result, nil
}
