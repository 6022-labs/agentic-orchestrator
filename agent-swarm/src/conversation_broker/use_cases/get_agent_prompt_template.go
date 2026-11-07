package use_cases

import (
	"context"

	"fmt"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type GetAgentPromptTemplateInterface interface {
	Get(
		ctx context.Context,
		request *requests.GetAgentPromptTemplateRequest,
	) (*responses.AgentPromptTemplateResponse, error)
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
	ctx context.Context,
	request *requests.GetAgentPromptTemplateRequest,
) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	template, err := uc.agentPromptTemplateRepository.GetByAgentAndKey(ctx,
		request.CollectionAddress,
		request.CollectionTokenId,
		request.Key,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent prompt template: %w", err)
	}
	if template == nil {
		return nil, cb_errors.NewNotFoundError(
			"prompt_template_not_found",
			"agent prompt template not found: "+request.Key.String(),
		)
	}

	return responses.AgentPromptTemplateResponseFromModel(template), nil
}
