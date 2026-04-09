package interfaces

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type MintSignerInterface interface {
	Sign(
		to, name, role string,
		cloneOf *uint64,
		wallets []string,
		images map[string]string,
		attributes map[string]string,
		nonce uuid.UUID,
		collectionAddress common.Address,
	) (*string, error)
}
