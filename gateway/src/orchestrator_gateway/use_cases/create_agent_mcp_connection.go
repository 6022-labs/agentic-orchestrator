package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"gorm.io/gorm"
)

type CreateAgentMCPConnectionInterface interface {
	Create(request *requests.CreateAgentMCPConnectionRequest) error
}

type CreateAgentMCPConnection struct {
	agentRepository        repositories.AgentRepositoryInterface
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface
}

func NewCreateAgentMCPConnection(
	agentRepository repositories.AgentRepositoryInterface,
	mcpConnectionRequester interfaces.MCPConnectionRequesterInterface,
) *CreateAgentMCPConnection {
	return &CreateAgentMCPConnection{
		agentRepository:        agentRepository,
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

	agent, err := c.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
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
