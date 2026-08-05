package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type RegisterBridgeInterface interface {
	// Register upserts a bridge. The second return value is true when the
	// caller-supplied id already existed (caller can surface 200 vs 201).
	Register(request *requests.RegisterBridgeRequest) (*responses.BridgeResponse, bool, error)
}

type RegisterBridge struct {
	bridgeRepository repositories.BridgeRepositoryInterface
}

func NewRegisterBridge(bridgeRepository repositories.BridgeRepositoryInterface) *RegisterBridge {
	return &RegisterBridge{bridgeRepository: bridgeRepository}
}

func (uc *RegisterBridge) Register(request *requests.RegisterBridgeRequest) (*responses.BridgeResponse, bool, error) {
	if err := request.Validate(); err != nil {
		return nil, false, err
	}

	if request.ID != nil {
		existing, err := uc.bridgeRepository.GetById(*request.ID)
		if err != nil {
			return nil, false, fmt.Errorf("failed to load existing bridge: %w", err)
		}
		if existing != nil {
			if existing.BaseUrl == request.BaseUrl && existing.Kind == request.Kind {
				return responses.BridgeResponseFromModel(existing), true, nil
			}
			existing.BaseUrl = request.BaseUrl
			existing.Kind = request.Kind
			saved, err := uc.bridgeRepository.Save(existing)
			if err != nil {
				return nil, false, fmt.Errorf("failed to update existing bridge: %w", err)
			}
			return responses.BridgeResponseFromModel(saved), true, nil
		}
	}

	bridge := &models.Bridge{
		Kind:    request.Kind,
		BaseUrl: request.BaseUrl,
	}
	if request.ID != nil {
		bridge.ID = *request.ID
	}

	saved, err := uc.bridgeRepository.Save(bridge)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save bridge: %w", err)
	}

	return responses.BridgeResponseFromModel(saved), false, nil
}
