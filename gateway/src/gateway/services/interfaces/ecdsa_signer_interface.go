package interfaces

type ECDSASignerInterface interface {
	// RecoverAddress recovers the signer address of a 65-byte signature over hash, normalizing legacy V (27/28).
	RecoverAddress(hash []byte, signature []byte) (string, error)
}
