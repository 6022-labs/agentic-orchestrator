package services

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/ethereum/go-ethereum/crypto"
)

type SignedMessageValidator struct {
	ecdsaSigner interfaces.ECDSASignerInterface
}

func NewSignedMessageValidator(ecdsaSigner interfaces.ECDSASignerInterface) *SignedMessageValidator {
	return &SignedMessageValidator{ecdsaSigner: ecdsaSigner}
}

func (v *SignedMessageValidator) Validate(message, signature, expectedAddress string) bool {
	// Ethereum adds "\x19Ethereum Signed Message:\n" prefix to prevent replay attacks
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return false
	}

	recoveredAddress, err := v.ecdsaSigner.RecoverAddress(hash.Bytes(), sigBytes)
	if err != nil {
		return false
	}

	return strings.EqualFold(recoveredAddress, expectedAddress)
}
