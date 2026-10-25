package services

import (
	"context"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

type BlockchainRequester struct {
	AbstractRpcRequester
}

type NewBlockchainRequesterParams struct {
	dig.In

	RpcNodeSettings *settings.RpcNodeSettings
	Client          *ethclient.Client `name:"http"`
}

func NewBlockchainRequester(
	params NewBlockchainRequesterParams,
) *BlockchainRequester {
	return &BlockchainRequester{
		AbstractRpcRequester: AbstractRpcRequester{client: params.Client, rpcNodeSettings: params.RpcNodeSettings},
	}
}

func (r *BlockchainRequester) GetCurrentBlockNumber(ctx context.Context) (*uint64, error) {
	var blockNumber uint64
	err := r.retryRpc(ctx, func() error {
		var err error
		blockNumber, err = r.client.BlockNumber(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &blockNumber, nil
}
