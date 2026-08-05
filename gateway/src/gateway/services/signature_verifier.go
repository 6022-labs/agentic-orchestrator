package services

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type SignatureVerifier struct {
	ecdsaSigner interfaces.ECDSASignerInterface
}

func NewSignatureVerifier(ecdsaSigner interfaces.ECDSASignerInterface) *SignatureVerifier {
	return &SignatureVerifier{ecdsaSigner: ecdsaSigner}
}

func (v *SignatureVerifier) Recover(message string, signatureHex string) (common.Address, error) {
	// Ethereum adds "\x19Ethereum Signed Message:\n" prefix to prevent replay attacks.
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil {
		return common.Address{}, fmt.Errorf("signature is not valid hex: %w", err)
	}

	recovered, err := v.ecdsaSigner.RecoverAddress(hash.Bytes(), sigBytes)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(recovered), nil
}
