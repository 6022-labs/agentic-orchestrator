package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type CreateAgentContextInterface interface {
	Create(request *requests.CreateAgentContextRequest) error
}

type CreateAgentContext struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
	agentContextRequester  interfaces.AgentContextRequesterInterface
}

func NewCreateAgentContext(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	agentContextRequester interfaces.AgentContextRequesterInterface,
) *CreateAgentContext {
	return &CreateAgentContext{
		agentRegistryRequester: agentRegistryRequester,
		agentContextRequester:  agentContextRequester,
	}
}

func (c *CreateAgentContext) Create(request *requests.CreateAgentContextRequest) error {
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
		return errors.NewValidationError("agentCollectionTokenId", "agent not found")
	}
	if !agent.IsInternal {
		return errors.NewValidationError("agentCollectionTokenId", "Agent is not internal")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}
	if agent.Role != enums.AgentRoleExpert {
		return errors.NewValidationError("agentCollectionTokenId", "Agent is not a expert, cannot create context")
	}

	err = c.agentContextRequester.CreateAgentContext(agent.ID, request.Key)
	if err != nil {
		return err
	}

	return nil
}
