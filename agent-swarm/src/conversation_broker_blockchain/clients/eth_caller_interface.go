package clients

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

// EthCallerInterface is the minimal eth_call surface the CCIP-Read resolver
// needs; *ethclient.Client satisfies it.
type EthCallerInterface interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}
