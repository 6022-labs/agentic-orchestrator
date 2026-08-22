package prompts

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ConversationPromptBuilder struct {
	promptTemplateProvider *promptTemplateProvider
	registryAgentResolver  interfaces.RegistryAgentResolverInterface
}

func NewConversationPromptBuilder(
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface,
	registryAgentResolver interfaces.RegistryAgentResolverInterface,
) *ConversationPromptBuilder {
	return &ConversationPromptBuilder{
		promptTemplateProvider: &promptTemplateProvider{agentPromptTemplateRepository: agentPromptTemplateRepository},
		registryAgentResolver:  registryAgentResolver,
	}
}

func (b *ConversationPromptBuilder) BuildInput(
	ctx context.Context,
	platform string,
	agent *models.Participant,
	role *enums.ParticipantRole,
	conversation *models.Conversation,
	messages []*models.ConversationMessage,
	participants []*models.Participant,
) (*models.RemoteAgentInput, error) {
	if role == nil {
		return nil, fmt.Errorf("participant %s has no role", agent.ID)
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages to build a turn from")
	}

	systemKey, defaultKey, err := b.templateKeysForRole(*role)
	if err != nil {
		return nil, err
	}

	parameters, err := b.buildParameters(ctx, conversation, messages, participants, len(messages))
	if err != nil {
		return nil, err
	}

	instructions, err := b.parseTemplate(ctx, agent, systemKey, parameters)
	if err != nil {
		return nil, err
	}

	turnPrompt, err := b.parseTemplate(ctx, agent, defaultKey, parameters)
	if err != nil {
		return nil, err
	}

	input := b.mapMessagesToInput(agent, messages, participants)
	input = append(input, &models.RemoteAgentInputMessage{
		Role:    enums.RemoteAgentInputRoleUser,
		Content: turnPrompt,
	})

	metadata := map[string]string{
		"conversationId": conversation.ID.String(),
		"platform":       platform,
	}
	if agent.EnsDomain != nil {
		metadata["ensDomain"] = *agent.EnsDomain
	}

	return &models.RemoteAgentInput{
		Instructions: instructions,
		Input:        input,
		Metadata:     metadata,
	}, nil
}

func (b *ConversationPromptBuilder) BuildConcludeInput(
	ctx context.Context,
	turn *models.RemoteAgentInput,
	completions []string,
	agent *models.Participant,
	role *enums.ParticipantRole,
	conversation *models.Conversation,
	messages []*models.ConversationMessage,
	participants []*models.Participant,
) (*models.RemoteAgentInput, error) {
	parameters, err := b.buildParameters(ctx, conversation, messages, participants, len(turn.Input)+len(completions))
	if err != nil {
		return nil, err
	}

	concludePrompt, err := b.parseTemplate(ctx, agent, enums.AgentPromptTemplateKeyFacilitatorConclude, parameters)
	if err != nil {
		return nil, err
	}

	input := make([]*models.RemoteAgentInputMessage, 0, len(turn.Input)+len(completions)+1)
	input = append(input, turn.Input...)
	for _, completion := range completions {
		input = append(input, &models.RemoteAgentInputMessage{
			Role:    enums.RemoteAgentInputRoleAssistant,
			Content: completion,
		})
	}
	input = append(input, &models.RemoteAgentInputMessage{
		Role:    enums.RemoteAgentInputRoleUser,
		Content: concludePrompt,
	})

	return &models.RemoteAgentInput{
		Instructions: turn.Instructions,
		Input:        input,
		Metadata:     turn.Metadata,
	}, nil
}

func (b *ConversationPromptBuilder) parseTemplate(
	ctx context.Context,
	agent *models.Participant,
	key enums.AgentPromptTemplateKey,
	parameters map[string]any,
) (string, error) {
	tpl, err := b.promptTemplateProvider.templateFor(ctx, agent, key)
	if err != nil {
		return "", err
	}

	parsed, err := b.parsePrompt(tpl, parameters)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", key, err)
	}
	return parsed, nil
}

func (b *ConversationPromptBuilder) buildParameters(
	ctx context.Context,
	conversation *models.Conversation,
	messages []*models.ConversationMessage,
	participants []*models.Participant,
	messageCount int,
) (map[string]any, error) {
	channelMembers := make([]string, 0, len(participants))
	for _, p := range participants {
		if p.EnsDomain != nil {
			channelMembers = append(channelMembers, *p.EnsDomain)
		} else if p.DisplayName != nil {
			channelMembers = append(channelMembers, *p.DisplayName)
		}
	}

	var description string
	if conversation.Description != nil {
		description = *conversation.Description
	}

	participantsById := lo.SliceToMap(
		participants,
		func(participant *models.Participant) (uuid.UUID, *models.Participant) {
			return participant.ID, participant
		},
	)

	initialMessage := messages[0]
	latestMessage := messages[len(messages)-1]

	return NewPromptParametersBuilder(b.registryAgentResolver).
		WithParam([]string{"ChannelMembers"}, channelMembers, lo.ToPtr(PromptParameterTypePotentialAgentArray)).
		WithParam([]string{"ChannelDescription"}, description, nil).
		WithParam([]string{"ConversationId"}, conversation.ID.String(), nil).
		WithParam([]string{"ConversationMessageCount"}, messageCount, nil).
		WithParam([]string{"InitialRequest", "Content"}, initialMessage.Content, nil).
		WithParam([]string{"LatestRequest", "Content"}, latestMessage.Content, nil).
		WithParam(
			[]string{"InitialRequest", "Author"},
			b.authorName(initialMessage, participantsById),
			lo.ToPtr(PromptParameterTypePotentialAgent),
		).
		WithParam(
			[]string{"LatestRequest", "Author"},
			b.authorName(latestMessage, participantsById),
			lo.ToPtr(PromptParameterTypePotentialAgent),
		).
		Build(ctx)
}

func (b *ConversationPromptBuilder) mapMessagesToInput(
	agent *models.Participant,
	messages []*models.ConversationMessage,
	participants []*models.Participant,
) []*models.RemoteAgentInputMessage {
	participantsById := lo.SliceToMap(
		participants,
		func(participant *models.Participant) (uuid.UUID, *models.Participant) {
			return participant.ID, participant
		},
	)

	input := make([]*models.RemoteAgentInputMessage, len(messages))
	for i, message := range messages {
		if message.ParticipantID == agent.ID {
			input[i] = &models.RemoteAgentInputMessage{
				Role:    enums.RemoteAgentInputRoleAssistant,
				Content: message.Content,
			}
			continue
		}

		author := b.authorName(message, participantsById)
		content := message.Content
		if author != "" {
			content = fmt.Sprintf("%s: %s", author, content)
		}
		input[i] = &models.RemoteAgentInputMessage{
			Role:    enums.RemoteAgentInputRoleUser,
			Content: content,
		}
	}
	return input
}

func (b *ConversationPromptBuilder) authorName(
	message *models.ConversationMessage,
	participantsById map[uuid.UUID]*models.Participant,
) string {
	author := participantsById[message.ParticipantID]
	if author == nil {
		return ""
	}
	if author.EnsDomain != nil {
		return *author.EnsDomain
	}
	if author.DisplayName != nil {
		return *author.DisplayName
	}
	return ""
}

func (b *ConversationPromptBuilder) templateKeysForRole(
	role enums.ParticipantRole,
) (enums.AgentPromptTemplateKey, enums.AgentPromptTemplateKey, error) {
	switch role {
	case enums.ParticipantRoleClone:
		return enums.AgentPromptTemplateKeyCloneSystem, enums.AgentPromptTemplateKeyCloneDefault, nil
	case enums.ParticipantRoleExpert:
		return enums.AgentPromptTemplateKeyExpertSystem, enums.AgentPromptTemplateKeyExpertDefault, nil
	case enums.ParticipantRoleFacilitator:
		return enums.AgentPromptTemplateKeyFacilitatorSystem, enums.AgentPromptTemplateKeyFacilitatorDefault, nil
	default:
		return "", "", fmt.Errorf("no templates for role %s", role)
	}
}
