package services

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type ECDSASignature struct{}

func NewECDSASignature() *ECDSASignature {
	return &ECDSASignature{}
}

func (e *ECDSASignature) Sign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	hash := e.hashMessage(data)

	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func (e *ECDSASignature) RecoverPublicKey(data []byte, signature []byte) (*string, error) {
	hash := e.hashMessage(data)

	// Ensure signature is the correct length
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signature length")
	}

	// Recover public key
	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return nil, err
	}

	// Compute Ethereum address from public key
	address := crypto.PubkeyToAddress(*pubKey).Hex()

	return &address, nil
}

func (e *ECDSASignature) hashMessage(message []byte) []byte {
	return crypto.Keccak256(message)
}
