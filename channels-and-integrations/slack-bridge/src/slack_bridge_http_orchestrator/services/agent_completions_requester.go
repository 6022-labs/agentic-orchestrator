package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/clients"
	http_orchestrator_models "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/models"
	"gorm.io/gorm"
)

type AgentCompletionsRequester struct {
	client clients.OrchestratorGatewayClientInterface

	slackAgentRepository repositories.SlackAgentRepositoryInterface
}

func NewAgentCompletionsRequester(
	client clients.OrchestratorGatewayClientInterface,
	slackAgentRepository repositories.SlackAgentRepositoryInterface,
) *AgentCompletionsRequester {
	return &AgentCompletionsRequester{
		client:               client,
		slackAgentRepository: slackAgentRepository,
	}
}

func (a *AgentCompletionsRequester) Completions(
	slackChannel *models.SlackChannel,
	slackThread *models.SlackThread,
	slackMessages []*models.SlackMessage,
	collectionAddress string,
	collectionTokenId uint64,
	forceCompletion bool,
) ([]*string, *bool, error) {
	inputs := make([]*http_orchestrator_models.ExternalAgentComputeInputRequest, len(slackMessages))
	for i, message := range slackMessages {
		var author *string

		if message.AuthorCollectionAddress != nil && message.AuthorCollectionTokenId != nil {
			authorAgent, err := a.slackAgentRepository.GetByCollectionAddressAndCollectionTokenId(*message.AuthorCollectionAddress, *message.AuthorCollectionTokenId)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, nil, fmt.Errorf("failed to fetch author agent: %w", err)
			}
			if authorAgent == nil {
				return nil, nil, fmt.Errorf("author agent not found for address %s and token ID %d", *message.AuthorCollectionAddress, *message.AuthorCollectionTokenId)
			}
			author = &authorAgent.EnsDomain
		}

		inputs[i] = &http_orchestrator_models.ExternalAgentComputeInputRequest{
			Author:    author,
			Content:   message.Content,
			Timestamp: message.CreatedAt(),
		}
	}

	channelMembers := make([]string, len(slackChannel.Members)+len(slackChannel.Agents))
	copy(channelMembers, slackChannel.Members)
	for i, agent := range slackChannel.Agents {
		channelMembers[len(slackChannel.Members)+i] = agent.EnsDomain
	}

	externalAgentComputeContextConversationContextRequest := &http_orchestrator_models.ExternalAgentComputeContextConversationContextRequest{
		Id: slackThread.Id,
	}
	externalAgentComputeContextConversationContextRequest.SetAgents(slackChannel.Agents)

	completionsRequest := &http_orchestrator_models.ExternalAgentComputeRequest{
		Platform:            "slack",
		Inputs:              inputs,
		ForceConclusion:     forceCompletion,
		ConversationContext: externalAgentComputeContextConversationContextRequest,
		PromptParameters: map[string]any{
			"ChannelMembers":      channelMembers,
			"ChannelMembers.type": "potential_agent[]",
			"ChannelDescription":  slackChannel.Description,
		},
	}

	completionsResult, err := a.client.Completions(collectionAddress, collectionTokenId, completionsRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get completions: %w", err)
	}

	return completionsResult.Responses, &completionsResult.IsFinal, nil
}
