package requests

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/compute"
	"github.com/ethereum/go-ethereum/common"
)

// For the moment agents must follow the ENS domain pattern formatted as:
// <agent-name>.protocol6022.<domain-extension>
var agentRegexPattern = regexp.MustCompile(`^[a-z0-9-\.]+\.protocol6022\.[a-z]{2,}$`)

type AgentComputeInputRequest struct {
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type AgentComputeConversationContextAgentRequest struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionTokenId uint64 `json:"collectionTokenId"`
}

type AgentComputeConversationContextRequest struct {
	Id     string                                         `json:"id"`
	Agents []*AgentComputeConversationContextAgentRequest `json:"agents"`
}

type AgentComputeRequest struct {
	AuthenticatedWallet    *common.Address                         `swaggerignore:"true"`
	Platform               enums.AgentContextMessagePlatform       `json:"platform"`
	ForceConclusion        bool                                    `json:"forceConclusion"`
	Inputs                 []*AgentComputeInputRequest             `json:"inputs"`
	ConversationContext    *AgentComputeConversationContextRequest `json:"conversationContext"`
	PromptParameters       map[string]any                          `json:"promptParameters"`
	AgentCollectionAddress string                                  `param:"agentCollectionAddress" swaggerignore:"true"`
	AgentCollectionTokenId uint64                                  `param:"agentCollectionTokenId" swaggerignore:"true"`
}

func (r *AgentComputeRequest) ValidateAndSanitize() error {
	if len(r.Inputs) == 0 {
		return errors.NewValidationError("inputs", "cannot be empty")
	}

	for index, input := range r.Inputs {
		if len(strings.TrimSpace(input.Author)) != 0 && !agentRegexPattern.MatchString(input.Author) {
			return errors.NewValidationError("inputs.author", "must be a valid ENS domain if setup, check input n°"+strconv.Itoa(index+1))
		}

		if len(strings.TrimSpace(input.Content)) == 0 {
			return errors.NewValidationError("inputs.content", "content cannot be empty, check input n°"+strconv.Itoa(index+1))
		}

		if input.Timestamp.IsZero() {
			return errors.NewValidationError("inputs.timestamp", "timestamp cannot be empty, check¨ input n°"+strconv.Itoa(index+1))
		}
	}

	if !enums.GetAllAgentContextMessagePlatforms().Contains(r.Platform) {
		r.Platform = enums.AgentContextMessagePlatformNone
	}

	if len(strings.TrimSpace(r.ConversationContext.Id)) == 0 {
		return errors.NewValidationError("conversationContext.id", "cannot be empty")
	}

	return nil
}

func (r *AgentComputeRequest) ToComputeContext() *compute.ComputeContext {
	result := &compute.ComputeContext{
		Platform:        r.Platform,
		ForceConclusion: r.ForceConclusion,
		ConversationContext: &compute.ConversationContext{
			Id: r.ConversationContext.Id,
		},
		PromptParameters: r.PromptParameters,
	}

	for _, agent := range r.ConversationContext.Agents {
		result.ConversationContext.Agents = append(result.ConversationContext.Agents, &compute.ConversationContextAgent{
			CollectionAddress: agent.CollectionAddress,
			CollectionTokenId: agent.CollectionTokenId,
		})
	}

	return result
}
