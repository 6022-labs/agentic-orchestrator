package services

import (
	"context"

	registry_abi "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

type AgentCollectionsManagerRequester struct {
	client                          *ethclient.Client
	agentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

type NewAgentCollectionsManagerRequesterParams struct {
	dig.In

	Client                          *ethclient.Client `name:"http"`
	AgentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

func NewAgentCollectionsManagerRequester(
	params NewAgentCollectionsManagerRequesterParams,
) *AgentCollectionsManagerRequester {
	return &AgentCollectionsManagerRequester{
		client:                          params.Client,
		agentCollectionsManagerSettings: params.AgentCollectionsManagerSettings,
	}
}

func (a *AgentCollectionsManagerRequester) GetAgentCollectionsManagerCollectionCreatedEvent(
	ctx context.Context,
	from, to uint64,
) ([]*registry_abi.AgentCollectionsManagerCollectionCreated, error) {
	agentCollectionsManager, err := registry_abi.NewAgentCollectionsManager(
		a.agentCollectionsManagerSettings.SmartContractAddress,
		a.client,
	)
	if err != nil {
		return nil, err
	}

	iter, err := agentCollectionsManager.FilterCollectionCreated(&bind.FilterOpts{
		Start:   from,
		End:     &to,
		Context: ctx,
	}, nil)
	if err != nil {
		return nil, err
	}

	var allEvents []*registry_abi.AgentCollectionsManagerCollectionCreated

	for iter.Next() {
		allEvents = append(allEvents, iter.Event)
	}

	return allEvents, nil
}
