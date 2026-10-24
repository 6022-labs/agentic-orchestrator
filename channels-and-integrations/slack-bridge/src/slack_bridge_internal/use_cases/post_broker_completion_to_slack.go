package use_cases

import (
	"context"
	"fmt"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases/requests"
)

type PostBrokerCompletionToSlackInterface interface {
	Post(ctx context.Context, request *requests.PostBrokerCompletionToSlackRequest) error
}

type PostBrokerCompletionToSlack struct {
	slackConversationRepository repositories.SlackConversationRepositoryInterface
	slackAgentRepository        repositories.SlackAgentRepositoryInterface
	slackThreadRequester        interfaces.SlackThreadRequesterInterface
}

func NewPostBrokerCompletionToSlack(
	slackConversationRepository repositories.SlackConversationRepositoryInterface,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
	slackThreadRequester interfaces.SlackThreadRequesterInterface,
) *PostBrokerCompletionToSlack {
	return &PostBrokerCompletionToSlack{
		slackConversationRepository: slackConversationRepository,
		slackAgentRepository:        slackAgentRepository,
		slackThreadRequester:        slackThreadRequester,
	}
}

func (uc *PostBrokerCompletionToSlack) Post(
	ctx context.Context,
	request *requests.PostBrokerCompletionToSlackRequest,
) error {
	if err := request.ValidateAndSanitize(); err != nil {
		return err
	}

	link, err := uc.slackConversationRepository.GetByBrokerConversationId(ctx, request.ConversationId)
	if err != nil {
		return fmt.Errorf("failed to lookup slack conversation link: %w", err)
	}
	if link == nil {
		return common_errors.NewNotFoundError(
			"slack_conversation_link_not_found",
			"no slack conversation is linked to this broker conversation",
		)
	}

	slackAgent, err := uc.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		request.CollectionAddress,
		request.CollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to lookup slack agent: %w", err)
	}
	if slackAgent == nil {
		return common_errors.NewNotFoundError(
			"slack_agent_not_found",
			"no slack agent exists for this collection address and token id",
		)
	}

	if err := uc.slackThreadRequester.PushCompletionToThread(ctx,
		request.Content,
		slackAgent,
		link.SlackChannelId,
		link.SlackThreadTs,
	); err != nil {
		return fmt.Errorf("failed to post completion to slack thread: %w", err)
	}

	return nil
}
