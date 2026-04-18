package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/clients"
	http_compute_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type AgentCompletionsRequester struct {
	client clients.AgentReasoningAgentComputeClientInterface
}

func NewAgentCompletionsRequester(client clients.AgentReasoningAgentComputeClientInterface) *AgentCompletionsRequester {
	return &AgentCompletionsRequester{
		client: client,
	}
}

func (r *AgentCompletionsRequester) Completions(
	platform string,
	agent *models.Participant,
	conversation *models.Conversation,
	messages []*models.ConversationMessage,
	participants []*models.Participant,
	forceCompletion bool,
) ([]*string, bool, error) {
	if agent == nil || agent.CollectionAddress == nil || agent.CollectionTokenId == nil {
		return nil, false, fmt.Errorf("participant is not an agent (missing on-chain identity)")
	}

	participantsById := lo.SliceToMap(participants, func(participant *models.Participant) (uuid.UUID, *models.Participant) {
		return participant.ID, participant
	})

	inputs := make([]*http_compute_models.ExternalAgentComputeInputRequest, len(messages))
	for i, message := range messages {
		var author *string
		authoringParticipant := participantsById[message.ParticipantID]
		if authoringParticipant != nil && authoringParticipant.EnsDomain != nil {
			ens := *authoringParticipant.EnsDomain
			author = &ens
		}
		inputs[i] = &http_compute_models.ExternalAgentComputeInputRequest{
			Author:    author,
			Content:   message.Content,
			Timestamp: message.CreatedAt,
		}
	}

	agentRequests := make([]*http_compute_models.ExternalAgentComputeContextConversationAgentRequest, 0, len(participants))
	channelMembers := make([]string, 0, len(participants))
	for _, p := range participants {
		if p.CollectionAddress != nil && p.CollectionTokenId != nil {
			agentRequests = append(agentRequests, &http_compute_models.ExternalAgentComputeContextConversationAgentRequest{
				CollectionAddress: *p.CollectionAddress,
				CollectionTokenId: *p.CollectionTokenId,
			})
			if p.EnsDomain != nil {
				channelMembers = append(channelMembers, *p.EnsDomain)
			}
		} else if p.DisplayName != nil {
			channelMembers = append(channelMembers, *p.DisplayName)
		}
	}

	var description string
	if conversation.Description != nil {
		description = *conversation.Description
	}

	request := &http_compute_models.ExternalAgentComputeRequest{
		Platform: platform,
		ConversationContext: &http_compute_models.ExternalAgentComputeContextConversationContextRequest{
			Id:     conversation.ID.String(),
			Agents: agentRequests,
		},
		ForceConclusion: forceCompletion,
		Inputs:          inputs,
		PromptParameters: map[string]any{
			"ChannelMembers":      channelMembers,
			"ChannelMembers.type": "potential_agent[]",
			"ChannelDescription":  description,
		},
	}

	result, err := r.client.Completions(*agent.CollectionAddress, *agent.CollectionTokenId, request)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get completions: %w", err)
	}

	return result.Responses, result.IsFinal, nil
}
