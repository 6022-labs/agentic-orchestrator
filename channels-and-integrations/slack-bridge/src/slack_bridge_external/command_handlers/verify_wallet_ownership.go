package command_handlers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers/requests"
	external_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"go.uber.org/zap"
)

const (
	VERIFY_WALLET_FRONT_PAGE_PATH = "/verify-wallet"
	RESULT_TEMPLATE               = "To verify your wallet, please click on the following link: %s"
)

type VerifyWalletOwnershipInterface interface {
	Verify(ctx context.Context, request *requests.VerifyWalletOwnershipRequest) (*string, error)
}

type VerifyWalletOwnership struct {
	channelMetrics metrics_interfaces.ChannelMetricsInterface
	channelTracer  traces_interfaces.ChannelTracerInterface

	logger                        *zap.Logger
	verifyWalletOwnershipSettings *settings.VerifyWalletOwnershipSettings
	walletOwnershipChallenge      interfaces.WalletOwnershipChallengeInterface
	slackUserSynchronizerCached   external_interfaces.SlackUserSynchronizerCachedInterface
}

func NewVerifyWalletOwnership(
	channelMetrics metrics_interfaces.ChannelMetricsInterface,
	channelTracer traces_interfaces.ChannelTracerInterface,
	logger *zap.Logger,
	verifyWalletOwnershipSettings *settings.VerifyWalletOwnershipSettings,
	walletOwnershipChallenge interfaces.WalletOwnershipChallengeInterface,
	slackUserSynchronizerCached external_interfaces.SlackUserSynchronizerCachedInterface,
) *VerifyWalletOwnership {
	return &VerifyWalletOwnership{
		channelMetrics: channelMetrics,
		channelTracer:  channelTracer,

		logger:                        logger,
		walletOwnershipChallenge:      walletOwnershipChallenge,
		slackUserSynchronizerCached:   slackUserSynchronizerCached,
		verifyWalletOwnershipSettings: verifyWalletOwnershipSettings,
	}
}

func (h *VerifyWalletOwnership) Verify(
	ctx context.Context,
	request *requests.VerifyWalletOwnershipRequest,
) (commandResult *string, err error) {
	ctx, span := h.channelTracer.StartCommandHandling(ctx, "wallet.verify")
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
		h.channelMetrics.RecordCommand(
			ctx,
			"wallet.verify",
			metrics_interfaces.CommandTriggerFrom(ctx),
			outcome,
			time.Since(start),
		)
		span.End()
	}()

	if err = request.Validate(); err != nil {
		return nil, err
	}

	userToVerify, err := h.slackUserSynchronizerCached.Sync(ctx, request.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to sync requester user : %w", err)
	}

	challenge, err := h.walletOwnershipChallenge.IssueChallenge(ctx, userToVerify.Id)
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	queryParams.Add("challenge", challenge.String())
	queryParams.Add("platform", "slack")

	verificationUrl := fmt.Sprintf(
		"%s%s?%s",
		h.verifyWalletOwnershipSettings.OrchestratorFrontBaseUrl,
		VERIFY_WALLET_FRONT_PAGE_PATH,
		queryParams.Encode(),
	)

	h.logger.Info("Generated wallet verification URL",
		zap.String("user_id", request.UserId),
		zap.String("url", verificationUrl),
	)

	result := fmt.Sprintf(RESULT_TEMPLATE, verificationUrl)
	return &result, nil
}
