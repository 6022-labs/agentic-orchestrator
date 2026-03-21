package interfaces

import "github.com/ethereum/go-ethereum/common"

type EIP712SignatureHelperInterface interface {
	BuildDomainSeparator(collectionAddress common.Address) (*common.Hash, error)
	NormalizeSignature(signature []byte) ([]byte, error)
}
