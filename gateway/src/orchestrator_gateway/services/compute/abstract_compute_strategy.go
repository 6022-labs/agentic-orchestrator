package compute

import (
	"fmt"
	"slices"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type ConversationContextAgent struct {
	CollectionAddress string
	CollectionTokenId uint64
}

type ConversationContext struct {
	Id     string
	Agents []*ConversationContextAgent
}

type ComputeContext struct {
	Platform            enums.AgentContextMessagePlatform
	ForceConclusion     bool
	ConversationContext *ConversationContext
	PromptParameters    map[string]any
}

type AbstractComputeStrategy struct {
	logger                       *zap.Logger
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface
}

func NewAbstractComputeStrategy(
	logger *zap.Logger,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
) *AbstractComputeStrategy {
	return &AbstractComputeStrategy{
		logger:                       logger,
		agentContextMessageRequester: agentContextMessageRequester,
	}
}

func (cs *AbstractComputeStrategy) saveMessages(
	agent *models.Agent,
	agentContextKey string,
	messages []*models.AgentContextMessage,
) error {
	requestMessages := lo.Map(messages, func(message *models.AgentContextMessage, _ int) *interfaces.AddMessageToAgentContextRequestMessage {
		return &interfaces.AddMessageToAgentContextRequestMessage{
			Role:           message.Role,
			Author:         message.Author,
			Content:        message.Content,
			Platform:       message.Platform,
			ConversationId: message.ConversationId,
		}
	})

	err := cs.agentContextMessageRequester.AddMessagesToAgentContext(&interfaces.AddMessageToAgentContextRequest{
		AgentId:    agent.ID,
		Messages:   requestMessages,
		ContextKey: &agentContextKey,
	})
	if err != nil {
		return fmt.Errorf("failed to save non-stored messages: %w", err)
	}

	return nil
}

func (cs *AbstractComputeStrategy) unionMessagesByHash(
	agentContextId uuid.UUID,
	firstArray []*models.AgentContextMessage,
	secondArray []*models.AgentContextMessage,
) []*models.AgentContextMessage {
	seen := make(map[string]struct{})
	result := make([]*models.AgentContextMessage, 0, len(firstArray)+len(secondArray))

	for _, msg := range append(firstArray, secondArray...) {
		msg.EnsureHash(agentContextId)
		if _, exists := seen[msg.Hash]; !exists {
			seen[msg.Hash] = struct{}{}
			result = append(result, msg)
		}
	}

	return result
}

func (cs *AbstractComputeStrategy) diffMessagesByHash(
	agentContextId uuid.UUID,
	firstArray []*models.AgentContextMessage,
	secondArray []*models.AgentContextMessage,
) []*models.AgentContextMessage {
	seen := make(map[string]struct{}, len(secondArray))
	for _, msg := range secondArray {
		msg.EnsureHash(agentContextId)
		seen[msg.Hash] = struct{}{}
	}

	result := make([]*models.AgentContextMessage, 0)
	for _, msg := range firstArray {
		msg.EnsureHash(agentContextId)
		if _, exists := seen[msg.Hash]; !exists {
			result = append(result, msg)
		}
	}

	return result
}

func (cs *AbstractComputeStrategy) groupMessages(messages []*models.AgentContextMessage, currentGroupingKey string) []*models.AgentContextMessage {
	slices.SortStableFunc(messages, func(a, b *models.AgentContextMessage) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		} else if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}
		return 0
	})

	groupedMessages := lo.GroupBy(messages, func(message *models.AgentContextMessage) string {
		if message.ConversationId != nil {
			return *message.ConversationId
		}

		return ""
	})

	type grouped struct {
		key       string
		group     []*models.AgentContextMessage
		createdAt time.Time
	}

	var groupsToInsertBeforeCurrentKey []grouped
	for key, group := range groupedMessages {
		if key == "" || key == currentGroupingKey {
			continue
		}

		groupsToInsertBeforeCurrentKey = append(groupsToInsertBeforeCurrentKey, grouped{
			key:       key,
			group:     group,
			createdAt: group[0].CreatedAt,
		})
	}

	slices.SortStableFunc(groupsToInsertBeforeCurrentKey, func(a, b grouped) int {
		if a.createdAt.Before(b.createdAt) {
			return -1
		} else if a.createdAt.After(b.createdAt) {
			return 1
		}
		return 0
	})

	result := groupedMessages[""]
	result = append(result, lo.Map(groupsToInsertBeforeCurrentKey, func(group grouped, _ int) *models.AgentContextMessage {
		return group.group[0]
	})...)

	slices.SortStableFunc(result, func(a, b *models.AgentContextMessage) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		} else if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}
		return 0
	})

	for _, group := range groupsToInsertBeforeCurrentKey {
		firstMessage := group.group[0]
		remainingToInsert := group.group[1:]

		_, insertIndex, found := lo.FindIndexOf(result, func(msg *models.AgentContextMessage) bool {
			return msg.Hash == firstMessage.Hash
		})

		if !found || insertIndex+1 == len(result) {
			result = append(result, remainingToInsert...)
			continue
		}

		result = append(result[:insertIndex+1], append(remainingToInsert, result[insertIndex+1:]...)...)
	}

	result = append(result, groupedMessages[currentGroupingKey]...)

	return result
}
