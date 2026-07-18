package helpers

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// IsEthAddress reports whether s is a 0x-prefixed 20-byte hex address.
// The prefix guard keeps it stricter than common.IsHexAddress alone.
func IsEthAddress(s string) bool {
	return strings.HasPrefix(s, "0x") && common.IsHexAddress(s)
}

// IsEthSignature reports whether s is a 0x-prefixed 65-byte hex signature.
func IsEthSignature(s string) bool {
	decoded, err := hexutil.Decode(s)
	return err == nil && len(decoded) == crypto.SignatureLength
}
