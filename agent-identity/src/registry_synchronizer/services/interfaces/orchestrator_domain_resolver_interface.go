package interfaces

import "github.com/ethereum/go-ethereum/common"

type OrchestratorDomainResolverInterface interface {
	ResolveGatewayDomain(address common.Address) (*string, error)
}
