package use_cases

import (
	"fmt"
	"net/url"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type AddMessageToAgentContextInterface interface {
	Add(request *requests.AddMessageToAgentContextRequest) error
}

type AddMessageToAgentContext struct {
	agentRegistryRequester       interfaces.AgentRegistryRequesterInterface
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface
}

func NewAddMessageToAgentContext(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	agentContextMessageRequester interfaces.AgentContextMessageRequesterInterface,
) *AddMessageToAgentContext {
	return &AddMessageToAgentContext{
		agentRegistryRequester:       agentRegistryRequester,
		agentContextMessageRequester: agentContextMessageRequester,
	}
}

func (c *AddMessageToAgentContext) Add(request *requests.AddMessageToAgentContextRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agent, err := c.agentRegistryRequester.GetAgent(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
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
