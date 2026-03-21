package services

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

type SignedMessageValidator struct {
}

func NewSignedMessageValidator() *SignedMessageValidator {
	return &SignedMessageValidator{}
}

func (v *SignedMessageValidator) Validate(message, signature, expectedAddress string) bool {
	// Ethereum adds "\x19Ethereum Signed Message:\n" prefix to prevent replay attacks
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// Convert hex signature to bytes
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return false
	}

	if len(sigBytes) != 65 {
		return false
	}
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	if sigBytes[64] != 0 && sigBytes[64] != 1 {
		return false
	}

	// Recover the public key from the signature
	pubKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return false
	}

	// Compute the recovered address
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	// Ensure recovered address matches provided address
	return strings.EqualFold(recoveredAddress, expectedAddress)
}
