package signature

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type ECDSASigner struct{}

func NewECDSASigner() *ECDSASigner {
	return &ECDSASigner{}
}

func (e *ECDSASigner) RecoverAddress(hash []byte, signature []byte) (string, error) {
	if len(signature) != 65 {
		return "", fmt.Errorf("invalid signature length")
	}

	sig := make([]byte, 65)
	copy(sig, signature)
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}
