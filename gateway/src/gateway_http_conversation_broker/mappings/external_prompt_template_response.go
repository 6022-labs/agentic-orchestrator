package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/clients"
	"github.com/samber/lo"
)

// ExternalPromptTemplateResponse is the broker's generated template wire shape, named locally to carry mapping.
type ExternalPromptTemplateResponse clients.ResponsesAgentPromptTemplateResponse

func (template *ExternalPromptTemplateResponse) ToUseCaseResponse() *responses.AgentPromptTemplateResponse {
	if template == nil {
		return nil
	}

	var key string
	if template.Key != nil {
		key = string(*template.Key)
	}

	var collectionTokenId uint64
	if template.CollectionTokenId != nil {
		collectionTokenId = uint64(*template.CollectionTokenId)
	}

	return &responses.AgentPromptTemplateResponse{
		CollectionAddress: lo.FromPtr(template.CollectionAddress),
		CollectionTokenId: collectionTokenId,
		Key:               key,
		Content:           lo.FromPtr(template.Content),
	}
}
