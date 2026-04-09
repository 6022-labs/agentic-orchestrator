package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type CreateAgentSecondStepInterface interface {
	Execute(request *requests.AgentByIdRequest) (*responses.CreateAgentSecondStepResponse, error)
}

type CreateAgentSecondStep struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewCreateAgentSecondStep(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
) *CreateAgentSecondStep {
	return &CreateAgentSecondStep{
		agentRegistryRequester: agentRegistryRequester,
	}
}

func (f *CreateAgentSecondStep) Execute(request *requests.AgentByIdRequest) (*responses.CreateAgentSecondStepResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return f.agentRegistryRequester.CreateAgentSecondStep(&interfaces.CreateRegistryAgentSecondStepRequest{
		AuthenticatedWallet: request.AuthenticatedWallet,
		AgentId:             request.AgentId,
	})
}
