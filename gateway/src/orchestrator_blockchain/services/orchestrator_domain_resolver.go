package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"
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

func (e *OrchestratorDomainResolver) ResolveAIDomain(address common.Address) (*string, error) {
	return e.getEnsTextRecodeValue(address, "ai_domain")
}

func (e *OrchestratorDomainResolver) ResolveDappDomain(address common.Address) (*string, error) {
	return e.getEnsTextRecodeValue(address, "dapp_domain")
}

func (e *OrchestratorDomainResolver) getEnsTextRecodeValue(address common.Address, key string) (*string, error) {
	reverseResolver, err := ens.NewReverseResolver(e.client)
	if err != nil {
		return nil, err
	}

	ensDomain, err := reverseResolver.Name(address)
	if err != nil {
		return nil, err
	}

	resolver, err := ens.NewResolver(e.client, ensDomain)
	if err != nil {
		return nil, err
	}

	domain, err := resolver.Text(key)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}
