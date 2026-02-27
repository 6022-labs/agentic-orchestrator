package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
)

type RemoveAgentMCPConnectionInterface interface {
	Remove(request *requests.RemoveAgentMCPConnectionRequest) error
}

type RemoveAgentMCPConnection struct {
	agentRepository        repositories.AgentRepositoryInterface
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface
}

func NewRemoveAgentMCPConnection(
	agentRepository repositories.AgentRepositoryInterface,
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface,
) *RemoveAgentMCPConnection {
	return &RemoveAgentMCPConnection{
		agentRepository:        agentRepository,
		mcpConnectionRequester: mcpConnectionRequester,
	}
}

func (r *RemoveAgentMCPConnection) Remove(request *requests.RemoveAgentMCPConnectionRequest) error {
	err := request.EnsureAuthenticatedWallet()
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	agent, err := r.agentRepository.GetByCollectionAddressAndCollectionTokenId(
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
		return errors.NewValidationError("agentCollectionTokenId", "Agent is human, cannot have MCP settings")
	}

	err = r.mcpConnectionRequester.Remove(request.MCPConnectionId)
	if err != nil {
		return err
	}

	return nil
}
