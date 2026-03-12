package services

import (
	"fmt"
	"math/big"
	"strings"

	registry_abi "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/dig"
)

type AgentCollectionsManagerRequester struct {
	client                          *ethclient.Client
	agentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

type newAgentCollectionsManagerRequesterParams struct {
	dig.In

	Client                          *ethclient.Client `name:"http"`
	AgentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

func NewAgentCollectionsManagerRequester(
	params newAgentCollectionsManagerRequesterParams,
) *AgentCollectionsManagerRequester {
	return &AgentCollectionsManagerRequester{
		client:                          params.Client,
		agentCollectionsManagerSettings: params.AgentCollectionsManagerSettings,
	}
}

func (a *AgentCollectionsManagerRequester) GetCollectionAddresses() ([]common.Address, error) {
	agentCollectionsManager, err := registry_abi.NewAgentCollectionsManager(a.agentCollectionsManagerSettings.SmartContractAddress, a.client)
	if err != nil {
		return nil, err
	}

	nextCollectionId, err := agentCollectionsManager.NextCollectionId(nil)
	if err != nil {
		return nil, err
	}

	return a.getCollectionAddresses(1, nextCollectionId.Uint64())
}

func (a *AgentCollectionsManagerRequester) getCollectionAddresses(from, to uint64) ([]common.Address, error) {
	agentCollectionABI, err := abi.JSON(strings.NewReader(registry_abi.AgentCollectionsManagerABI)) // You must expose raw ABI string
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent collection ABI: %w", err)
	}

	batch := make([]rpc.BatchElem, 0, to-from)
	results := make([]hexutil.Bytes, to-from)

	for i := from; i < to; i++ {
		data, err := agentCollectionABI.Pack("collections", big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to pack collections call: %w", err)
		}

		msg := map[string]interface{}{
			"to":   a.agentCollectionsManagerSettings.SmartContractAddress.Hex(),
			"data": "0x" + common.Bytes2Hex(data),
		}

		batch = append(batch, rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{msg, "latest"},
			Result: &results[i-from],
		})
	}

	if err := a.client.Client().BatchCall(batch); err != nil {
		return nil, fmt.Errorf("failed to batch call collections: %w", err)
	}

	collectionAddresses := make([]common.Address, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue
		}

		unpacked, err := agentCollectionABI.Unpack("collections", result)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack collections result: %w", err)
		}
		collectionAddress := *abi.ConvertType(unpacked[0], new(common.Address)).(*common.Address)

		collectionAddresses = append(collectionAddresses, collectionAddress)
	}

	return collectionAddresses, nil
}
