package use_cases

import (
	"fmt"
	"net/url"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"gorm.io/gorm"
)

type AddMessageToAgentContextInterface interface {
	Add(request *requests.AddMessageToAgentContextRequest) error
}

type AddMessageToAgentContext struct {
	agentRepository              repositories.AgentRepositoryInterface
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface
}

func NewAddMessageToAgentContext(
	agentRepository repositories.AgentRepositoryInterface,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
) *AddMessageToAgentContext {
	return &AddMessageToAgentContext{
		agentRepository:              agentRepository,
		agentContextMessageRequester: agentContextMessageRequester,
	}
}

func (c *AddMessageToAgentContext) Add(request *requests.AddMessageToAgentContextRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agent, err := c.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if agent == nil {
		return errors.NewValidationError("agentCollectionTokenId", "Agent not found")
	}
	if !agent.IsInternal {
		return errors.NewValidationError("agentCollectionTokenId", "Agent is not internal")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}

	contextKeyDecoded, err := url.QueryUnescape(request.ContextKey)
	if err != nil {
		return fmt.Errorf("failed to decode context key: %w", err)
	}

	err = c.agentContextMessageRequester.AddMessagesToAgentContext(&interfaces.AddMessageToAgentContextRequest{
		AgentId:    agent.ID,
		ContextKey: &contextKeyDecoded,
		Messages: []*interfaces.AddMessageToAgentContextRequestMessage{
			{
				Role:     request.Role,
				Content:  request.Message,
				Author:   request.AuthenticatedWallet.Hex(),
				Platform: enums.AgentContextMessagePlatformNone,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
