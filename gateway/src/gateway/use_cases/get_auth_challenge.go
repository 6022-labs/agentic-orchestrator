package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
)

type GetAuthChallengeInterface interface {
	Execute(request *requests.GetAuthChallengeRequest) (*responses.GetAuthChallengeResponse, error)
}

type GetAuthChallenge struct {
	nonceStore         interfaces.NonceStoreInterface
	siweMessageBuilder interfaces.SiweMessageBuilderInterface
}

func NewGetAuthChallenge(
	nonceStore interfaces.NonceStoreInterface,
	siweMessageBuilder interfaces.SiweMessageBuilderInterface,
) *GetAuthChallenge {
	return &GetAuthChallenge{
		nonceStore:         nonceStore,
		siweMessageBuilder: siweMessageBuilder,
	}
}

func (uc *GetAuthChallenge) Execute(
	request *requests.GetAuthChallengeRequest,
) (*responses.GetAuthChallengeResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	address := common.HexToAddress(request.Address)
	entry, err := uc.nonceStore.Generate(address, request.ChainId)
	if err != nil {
		return nil, err
	}

	message := uc.siweMessageBuilder.Build(address, entry.ChainId, entry.Nonce, entry.IssuedAt, entry.ExpiresAt)
	return &responses.GetAuthChallengeResponse{
		Message:   message,
		ExpiresAt: entry.ExpiresAt,
	}, nil
}
