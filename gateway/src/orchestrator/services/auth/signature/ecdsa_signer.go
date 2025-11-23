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

func (e *ECDSASigner) hashMessage(message []byte) []byte {
	return crypto.Keccak256(message)
}
