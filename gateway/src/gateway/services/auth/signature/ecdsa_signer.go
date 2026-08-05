package signature

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type ECDSASigner struct{}

func NewECDSASigner() *ECDSASigner {
	return &ECDSASigner{}
}

func (e *ECDSASigner) Sign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	hash := e.hashMessage(data)

	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func (e *ECDSASigner) RecoverPublicKey(data []byte, signature []byte) (*string, error) {
	address, err := e.RecoverAddress(e.hashMessage(data), signature)
	if err != nil {
		return nil, err
	}
	return &address, nil
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

func (e *ECDSASigner) hashMessage(message []byte) []byte {
	return crypto.Keccak256(message)
}
