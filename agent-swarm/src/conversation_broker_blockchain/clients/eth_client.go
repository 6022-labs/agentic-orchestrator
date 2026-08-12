package clients

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// EthClientInterface is the read-only RPC surface the generated contract bindings need.
// It matches bind.ContractCaller so *ethclient.Client satisfies it as-is.
type EthClientInterface interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}
