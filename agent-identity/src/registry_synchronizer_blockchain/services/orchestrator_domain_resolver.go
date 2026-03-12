package services

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

type OrchestratorDomainResolver struct {
	client *ethclient.Client
}

type newOrchestratorDomainResolverParams struct {
	dig.In

	Client *ethclient.Client `name:"ens"`
}

func NewOrchestratorDomainResolver(params newOrchestratorDomainResolverParams) *OrchestratorDomainResolver {
	return &OrchestratorDomainResolver{
		client: params.Client,
	}
}

func (e *OrchestratorDomainResolver) ResolveGatewayDomain(address common.Address) (*string, error) {
	return nil, fmt.Errorf("not implemented")
}
