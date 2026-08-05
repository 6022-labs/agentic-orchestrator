package interfaces

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// SiweMessageBuilderInterface renders the canonical EIP-4361 message the
// wallet signs; verify rebuilds the exact same text from the stored challenge.
type SiweMessageBuilderInterface interface {
	Build(address common.Address, chainId uint64, nonce string, issuedAt time.Time, expiresAt time.Time) string
}
