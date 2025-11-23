package use_cases

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/use_cases/requests"
)

type CreateAgentCollectionInterface interface {
	Create(request *requests.CreateAgentCollectionRequest) error
}

type CreateAgentCollection struct {
	agentCollectionSynchronizer services.AgentCollectionSynchronizerInterface
}

func NewCreateAgentCollection(agentCollectionSynchronizer services.AgentCollectionSynchronizerInterface) *CreateAgentCollection {
	return &CreateAgentCollection{
		agentCollectionSynchronizer: agentCollectionSynchronizer,
	}
}

func (u *CreateAgentCollection) Create(request *requests.CreateAgentCollectionRequest) error {
	return u.agentCollectionSynchronizer.Sync(
		request.CollectionAddress,
		request.CollectionName,
		request.CollectionSymbol,
	)
}
