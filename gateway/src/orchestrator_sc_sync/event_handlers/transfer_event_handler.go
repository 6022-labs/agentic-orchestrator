package event_handlers

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
)

type TransferEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1Transfer) error
}

type TransferEventHandler struct {
	updateAgentOwner use_cases.UpdateAgentOwnerInterface
}

func NewTransferEventHandler(
	updateAgentOwner use_cases.UpdateAgentOwnerInterface,
) *TransferEventHandler {
	return &TransferEventHandler{
		updateAgentOwner: updateAgentOwner,
	}
}

func (h *TransferEventHandler) Handle(event *abi.AgentCollectionV1Transfer) error {
	request := &requests.UpdateAgentOwnerRequest{
		NewOwner:          event.To.Hex(),
		CollectionAddress: event.Raw.Address.Hex(),
		CollectionTokenId: event.TokenId.Uint64(),
	}

	return h.updateAgentOwner.Update(request)
}
