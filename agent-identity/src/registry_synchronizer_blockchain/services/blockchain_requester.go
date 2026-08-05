package services

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

type BlockchainRequester struct {
	client *ethclient.Client
}

type NewBlockchainRequesterParams struct {
	dig.In

	Client *ethclient.Client `name:"http"`
}

func NewBlockchainRequester(
	params NewBlockchainRequesterParams,
) *BlockchainRequester {
	return &BlockchainRequester{
		client: params.Client,
	}
}

func (r *BlockchainRequester) GetCurrentBlockNumber() (*uint64, error) {
	blockNumber, err := r.client.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}

	return &blockNumber, nil
}
