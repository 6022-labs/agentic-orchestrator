package interfaces

import "github.com/ethereum/go-ethereum/common"

// SignatureVerifierInterface recovers the signer of an EIP-191 personal_sign
// signature over the given message.
type SignatureVerifierInterface interface {
	Recover(message string, signatureHex string) (common.Address, error)
}
