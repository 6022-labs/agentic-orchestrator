package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type CreateAgentFirstStepInterface interface {
	Execute(request *requests.CreateAgentFirstStepRequest) (*responses.AgentResponse, error)
}

type CreateAgentFirstStep struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewCreateAgentFirstStep(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
) *CreateAgentFirstStep {
	return &CreateAgentFirstStep{
		agentRegistryRequester: agentRegistryRequester,
	}
}

func (c *CreateAgentFirstStep) Execute(request *requests.CreateAgentFirstStepRequest) (*responses.AgentResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return c.agentRegistryRequester.CreateAgentFirstStep(&interfaces.CreateRegistryAgentFirstStepRequest{
		AuthenticatedWallet: request.AuthenticatedWallet,
		CollectionAddress:   request.CollectionAddress,
		Name:                request.Name,
		Role:                request.Role,
		WalletAddress:       request.WalletAddress,
		Mission:             request.Mission,
		CloneOfAgentId:      request.CloneOfAgentId,
		Image:               request.Image,
	})
}
