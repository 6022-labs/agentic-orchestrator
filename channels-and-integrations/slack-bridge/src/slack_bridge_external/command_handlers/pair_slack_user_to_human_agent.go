package command_handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"go.uber.org/zap"
)

type PairSlackUserToHumanAgentInterface interface {
	Pair(ctx context.Context, request *requests.PairSlackUserToHumanAgentRequest) error
}

type PairSlackUserToHumanAgent struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	logger                      *zap.Logger
	slackUserRepository         repositories.SlackUserRepositoryInterface
	slackAgentSynchronizer      interfaces.SlackAgentSynchronizerInterface
	slackUserSynchronizerCached interfaces.SlackUserSynchronizerCachedInterface
}

func NewPairSlackUserToHumanAgent(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	logger *zap.Logger,
	slackUserRepository repositories.SlackUserRepositoryInterface,
	slackAgentSynchronizer interfaces.SlackAgentSynchronizerInterface,
	slackUserSynchronizerCached interfaces.SlackUserSynchronizerCachedInterface,
) *PairSlackUserToHumanAgent {
	return &PairSlackUserToHumanAgent{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		logger:                      logger,
		slackUserRepository:         slackUserRepository,
		slackAgentSynchronizer:      slackAgentSynchronizer,
		slackUserSynchronizerCached: slackUserSynchronizerCached,
	}
}

func (p *PairSlackUserToHumanAgent) Pair(
	ctx context.Context,
	request *requests.PairSlackUserToHumanAgentRequest,
) (err error) {
	ctx, span := p.channelTracer.StartCommandHandling(ctx, "user.pair")
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.ChannelOutcomeHandled
		switch {
		case errors.As(err, new(*common_errors.ValidationError)):
			outcome = metrics_interfaces.ChannelOutcomeRejected
		case err != nil:
			outcome = metrics_interfaces.ChannelOutcomeFailed
			span.Fail(err)
		}
		p.channelMetrics.RecordCommand(
			ctx,
			"user.pair",
			metrics_interfaces.CommandTriggerFrom(ctx),
			outcome,
			time.Since(start),
		)
		span.End()
	}()

	if err = request.ValidateAndSanitize(); err != nil {
		return err
	}

	requesterUser, err := p.slackUserSynchronizerCached.Sync(ctx, request.RequesterSlackUserId)
	if err != nil {
		return fmt.Errorf("failed to sync requester user : %w", err)
	}
	if requesterUser.WalletAddress == nil {
		return fmt.Errorf("requester needs to verify their wallet first before pairing a human agent")
	}

	agent, err := p.slackAgentSynchronizer.Sync(ctx,
		request.HumanAgentCollectionAddress,
		request.HumanAgentCollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to sync human agent: %w", err)
	}
	if agent.Role != enums.AgentRoleHuman {
		return fmt.Errorf(
			"the specified agent %s %d is not a human agent, role is: %s",
			request.HumanAgentCollectionAddress,
			request.HumanAgentCollectionTokenId,
			agent.Role,
		)
	}
	if agent.OwnerWallet != *requesterUser.WalletAddress {
		return fmt.Errorf(
			"requester does not own the human agent %s %d, owner is: %s",
			request.HumanAgentCollectionAddress,
			request.HumanAgentCollectionTokenId,
			agent.OwnerWallet,
		)
	}

	alreadyPairedUserForAgent, err := p.slackUserRepository.GetByPairedHumanAgent(ctx,
		request.HumanAgentCollectionAddress,
		request.HumanAgentCollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to check if human agent is already paired: %w", err)
	}
	if alreadyPairedUserForAgent != nil {
		return fmt.Errorf(
			"the specified human agent %s %d is already paired to another slack user",
			request.HumanAgentCollectionAddress,
			request.HumanAgentCollectionTokenId,
		)
	}

	slackUserToPair, err := p.slackUserSynchronizerCached.Sync(ctx, request.SlackUserId)
	if err != nil {
		return fmt.Errorf("failed to sync slack user to pair: %w", err)
	}

	slackUserToPair.PairedHumanAgentCollectionAddress = &request.HumanAgentCollectionAddress
	slackUserToPair.PairedHumanAgentCollectionTokenId = &request.HumanAgentCollectionTokenId

	_, err = p.slackUserRepository.Save(ctx, slackUserToPair)
	if err != nil {
		return fmt.Errorf("failed to save pairing: %w", err)
	}

	return nil
}
