package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/helpers"
	"github.com/samber/lo"
)

// ExternalParticipantResponse is the broker's generated participant wire shape, named locally to carry mapping.
type ExternalParticipantResponse clients.ResponsesParticipantResponse

func (participant *ExternalParticipantResponse) ToUseCaseResponse() *responses.BrokerParticipantResponse {
	if participant == nil {
		return nil
	}

	var role *string
	if participant.Role != nil {
		role = lo.ToPtr(string(*participant.Role))
	}

	var collectionTokenId *uint64
	if participant.CollectionTokenId != nil {
		collectionTokenId = lo.ToPtr(uint64(*participant.CollectionTokenId))
	}

	return &responses.BrokerParticipantResponse{
		ID:                helpers.ParseUuid(participant.Id),
		CollectionAddress: participant.CollectionAddress,
		CollectionTokenId: collectionTokenId,
		EnsDomain:         participant.EnsDomain,
		Role:              role,
		DisplayName:       participant.DisplayName,
	}
}
