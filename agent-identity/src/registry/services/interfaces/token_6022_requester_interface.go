package interfaces

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Token6022RequesterInterface interface {
	BalanceOf(account common.Address) (*big.Int, error)
}
