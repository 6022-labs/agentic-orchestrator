package use_cases

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"gorm.io/gorm"
)

type RegisterBridgeInterface interface {
	// Register upserts a bridge. The second return value is true when the
	// caller-supplied id already existed (caller can surface 200 vs 201).
	Register(request *requests.RegisterBridgeRequest) (*models.Bridge, bool, error)
}

type RegisterBridge struct {
	bridgeRepository repositories.BridgeRepositoryInterface
}

func NewRegisterBridge(bridgeRepository repositories.BridgeRepositoryInterface) *RegisterBridge {
	return &RegisterBridge{bridgeRepository: bridgeRepository}
}

func (uc *RegisterBridge) Register(request *requests.RegisterBridgeRequest) (*models.Bridge, bool, error) {
	if err := request.Validate(); err != nil {
		return nil, false, err
	}

	if request.ID != nil {
		existing, err := uc.bridgeRepository.GetById(*request.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, fmt.Errorf("failed to load existing bridge: %w", err)
		}
		if existing != nil {
			if existing.CallbackUrl == request.CallbackUrl && existing.Kind == request.Kind {
				return existing, true, nil
			}
			existing.CallbackUrl = request.CallbackUrl
			existing.Kind = request.Kind
			saved, err := uc.bridgeRepository.Save(existing)
			if err != nil {
				return nil, false, fmt.Errorf("failed to update existing bridge: %w", err)
			}
			return saved, true, nil
		}
	}

	bridge := &models.Bridge{
		Kind:        request.Kind,
		CallbackUrl: request.CallbackUrl,
	}
	if request.ID != nil {
		bridge.ID = *request.ID
	}

	saved, err := uc.bridgeRepository.Save(bridge)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save bridge: %w", err)
	}

	return saved, false, nil
}
