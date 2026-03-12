package services

import (
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_blockchain/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
)

type Token6022Requester struct {
	client            *ethclient.Client
	token6022Settings *settings.Token6022Settings
}

type newToken6022RequesterParams struct {
	dig.In

	Client            *ethclient.Client `name:"http"`
	Token6022Settings *settings.Token6022Settings
}

func NewToken6022Requester(
	params newToken6022RequesterParams,
) *Token6022Requester {
	return &Token6022Requester{
		client:            params.Client,
		token6022Settings: params.Token6022Settings,
	}
}

func (r *Token6022Requester) BalanceOf(account common.Address) (*big.Int, error) {
	token6022, err := abi.NewToken6022(r.token6022Settings.SmartContractAddress, r.client)
	if err != nil {
		return nil, err
	}

	balanceOf, err := token6022.BalanceOf(nil, account)
	if err != nil {
		return nil, err
	}

	return balanceOf, nil
}
