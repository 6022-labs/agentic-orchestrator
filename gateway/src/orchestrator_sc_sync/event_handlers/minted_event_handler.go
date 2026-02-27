package event_handlers

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases/requests"
)

type MintedEventHandlerInterface interface {
	Handle(event *abi.AgentCollectionV1Minted) error
}

type MintedEventHandler struct {
	finishAgentCreation use_cases.FinishAgentCreationInterface
}

func NewMintedEventHandler(
	finishAgentCreation use_cases.FinishAgentCreationInterface,
) *MintedEventHandler {
	return &MintedEventHandler{
		finishAgentCreation: finishAgentCreation,
	}
}

func (h *MintedEventHandler) Handle(event *abi.AgentCollectionV1Minted) error {
	request := &requests.FinishAgentCreationRequest{
		Signer:            event.Signer.Hex(),
		CollectionAddress: event.Raw.Address.Hex(),
		CollectionTokenId: event.TokenId.Uint64(),
	}

	return h.finishAgentCreation.Finish(request)
}
