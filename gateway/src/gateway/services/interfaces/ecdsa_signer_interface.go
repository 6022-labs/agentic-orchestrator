package interfaces

import "crypto/ecdsa"

type ECDSASignerInterface interface {
	Sign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error)
	RecoverPublicKey(data []byte, signature []byte) (*string, error)
	// RecoverAddress recovers the signer address of a 65-byte signature over hash, normalizing legacy V (27/28).
	RecoverAddress(hash []byte, signature []byte) (string, error)
}
