package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAuthChallengeInterface interface {
	Get() (*responses.GetAuthChallengeResponse, error)
}

type GetAuthChallenge struct {
	nonceStore interfaces.NonceStoreInterface
}

func NewGetAuthChallenge(nonceStore interfaces.NonceStoreInterface) *GetAuthChallenge {
	return &GetAuthChallenge{nonceStore: nonceStore}
}

func (uc *GetAuthChallenge) Get() (*responses.GetAuthChallengeResponse, error) {
	nonce, expiresAt, err := uc.nonceStore.Generate()
	if err != nil {
		return nil, err
	}
	return &responses.GetAuthChallengeResponse{Nonce: nonce, ExpiresAt: expiresAt}, nil
}
