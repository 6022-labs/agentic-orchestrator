package interfaces

import "crypto/ecdsa"

type ECDSASignerInterface interface {
	Sign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error)
	RecoverPublicKey(data []byte, signature []byte) (*string, error)
}
