package interfaces

import "github.com/ethereum/go-ethereum/common"

type OrchestratorDomainResolverInterface interface {
	ResolveAIDomain(address common.Address) (*string, error)
	ResolveDappDomain(address common.Address) (*string, error)
}
