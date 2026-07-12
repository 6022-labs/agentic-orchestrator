package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAgentPromptTemplateInterface interface {
	Get(request *requests.GetAgentPromptTemplateRequest) (*responses.AgentPromptTemplateResponse, error)
}

type GetAgentPromptTemplate struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewGetAgentPromptTemplate(conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface) *GetAgentPromptTemplate {
	return &GetAgentPromptTemplate{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *GetAgentPromptTemplate) Get(request *requests.GetAgentPromptTemplateRequest) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return uc.conversationBrokerRequester.GetAgentPromptTemplate(
		request.AgentCollectionAddress,
		request.ParsedCollectionTokenId(),
		request.Key,
	)
}
