package interfaces

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type UpdateAgentSignerSignerInterface interface {
	Sign(tokenId big.Int, address string, nonce uuid.UUID, collectionAddress common.Address) (*string, error)
}
