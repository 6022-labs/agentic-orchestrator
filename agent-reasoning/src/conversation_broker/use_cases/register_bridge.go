package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
)

type RegisterBridgeInterface interface {
	Register(request *requests.RegisterBridgeRequest) (*models.Bridge, error)
}

type RegisterBridge struct {
	bridgeRepository repositories.BridgeRepositoryInterface
}

func NewRegisterBridge(bridgeRepository repositories.BridgeRepositoryInterface) *RegisterBridge {
	return &RegisterBridge{bridgeRepository: bridgeRepository}
}

func (uc *RegisterBridge) Register(request *requests.RegisterBridgeRequest) (*models.Bridge, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	bridge := &models.Bridge{
		Kind:        request.Kind,
		CallbackUrl: request.CallbackUrl,
	}

	saved, err := uc.bridgeRepository.Save(bridge)
	if err != nil {
		return nil, fmt.Errorf("failed to save bridge: %w", err)
	}

	return saved, nil
}
