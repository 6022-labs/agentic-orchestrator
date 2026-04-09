package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type CreateAgentMCPConnectionInterface interface {
	Create(request *requests.CreateAgentMCPConnectionRequest) error
}

type CreateAgentMCPConnection struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface
}

func NewCreateAgentMCPConnection(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface,
) *CreateAgentMCPConnection {
	return &CreateAgentMCPConnection{
		agentRegistryRequester: agentRegistryRequester,
		mcpConnectionRequester: mcpConnectionRequester,
	}
}

func (c *CreateAgentMCPConnection) Create(request *requests.CreateAgentMCPConnectionRequest) error {
	err := request.EnsureAuthenticatedWallet()
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	agent, err := c.agentRegistryRequester.GetAgent(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}
	if agent == nil {
		return errors.NewValidationError("agentCollectionTokenId", "agent not found")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}
	if agent.Role == enums.AgentRoleHuman {
		return errors.NewValidationError("agentCollectionTokenId", "Agent is human, cannot create MCP settings")
	}

	err = c.mcpConnectionRequester.Create(&interfaces.CreateMCPConnectionRequest{
		AgentId:    agent.ID,
		Url:        request.Url,
		Name:       request.Name,
		HeaderAuth: request.HeaderAuth,
	})
	if err != nil {
		return err
	}

	return nil
}
