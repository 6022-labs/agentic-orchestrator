package services

import (
	"context"
	"fmt"

	registry_abi "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

type AgentCollectionsManagerRequester struct {
	AbstractRpcRequester

	contractAbi                     *abi.ABI
	filterer                        *registry_abi.AgentCollectionsManagerFilterer
	agentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

type NewAgentCollectionsManagerRequesterParams struct {
	dig.In

	RpcNodeSettings                 *settings.RpcNodeSettings
	Client                          *ethclient.Client `name:"http"`
	AgentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

func NewAgentCollectionsManagerRequester(
	params NewAgentCollectionsManagerRequesterParams,
) (*AgentCollectionsManagerRequester, error) {
	contractAbi, err := registry_abi.AgentCollectionsManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent collections manager ABI: %w", err)
	}

	filterer, err := registry_abi.NewAgentCollectionsManagerFilterer(common.Address{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent collections manager filterer: %w", err)
	}

	return &AgentCollectionsManagerRequester{
		AbstractRpcRequester: AbstractRpcRequester{
			client:          params.Client,
			rpcNodeSettings: params.RpcNodeSettings,
		},
		contractAbi:                     contractAbi,
		filterer:                        filterer,
		agentCollectionsManagerSettings: params.AgentCollectionsManagerSettings,
	}, nil
}

func (a *AgentCollectionsManagerRequester) GetAgentCollectionsManagerCollectionCreatedEvent(
	ctx context.Context,
	from, to uint64,
) ([]*registry_abi.AgentCollectionsManagerCollectionCreated, error) {
	logs, err := a.filterEventLogs(ctx,
		from,
		to,
		[]common.Address{a.agentCollectionsManagerSettings.SmartContractAddress},
		a.contractAbi,
		"CollectionCreated",
	)
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionsManagerCollectionCreated, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseCollectionCreated(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CollectionCreated log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}
