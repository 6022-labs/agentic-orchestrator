package use_cases

import (
	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type GetAgentPromptTemplateInterface interface {
	Get(request *requests.GetAgentPromptTemplateRequest) (*responses.AgentPromptTemplateResponse, error)
}

type GetAgentPromptTemplate struct {
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface
}

func NewGetAgentPromptTemplate(
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface,
) *GetAgentPromptTemplate {
	return &GetAgentPromptTemplate{agentPromptTemplateRepository: agentPromptTemplateRepository}
}

func (uc *GetAgentPromptTemplate) Get(
	request *requests.GetAgentPromptTemplateRequest,
) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	template, err := uc.agentPromptTemplateRepository.GetByAgentAndKey(
		request.CollectionAddress,
		request.ParsedCollectionTokenId(),
		request.ParsedKey(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent prompt template: %w", err)
	}
	if template == nil {
		return nil, cb_errors.NewNotFoundError("agent prompt template", request.Key)
	}

	return responses.AgentPromptTemplateResponseFromModel(template), nil
}
