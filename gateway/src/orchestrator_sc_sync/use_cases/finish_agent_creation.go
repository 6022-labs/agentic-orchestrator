package use_cases

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
)

type FinishAgentCreationInterface interface {
	Finish(request *requests.FinishAgentCreationRequest) error
}

type FinishAgentCreation struct {
	agentSynchronizer services.AgentSynchronizerInterface
}

func NewFinishAgentCreation(
	agentSynchronizer services.AgentSynchronizerInterface,
) *FinishAgentCreation {
	return &FinishAgentCreation{
		agentSynchronizer: agentSynchronizer,
	}
}

func (f *FinishAgentCreation) Finish(request *requests.FinishAgentCreationRequest) error {
	return f.agentSynchronizer.Sync(
		request.CollectionAddress,
		request.CollectionTokenId,
	)
}
