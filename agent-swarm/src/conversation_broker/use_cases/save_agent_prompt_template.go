package use_cases

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/responses"
)

type SaveAgentPromptTemplateInterface interface {
	Save(
		ctx context.Context,
		request *requests.SaveAgentPromptTemplateRequest,
	) (*responses.AgentPromptTemplateResponse, error)
}

type SaveAgentPromptTemplate struct {
	agentRegistryRequester        interfaces.AgentRegistryRequesterInterface
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface
}

func NewSaveAgentPromptTemplate(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface,
) *SaveAgentPromptTemplate {
	return &SaveAgentPromptTemplate{
		agentRegistryRequester:        agentRegistryRequester,
		agentPromptTemplateRepository: agentPromptTemplateRepository,
	}
}

func (uc *SaveAgentPromptTemplate) Save(
	ctx context.Context,
	request *requests.SaveAgentPromptTemplateRequest,
) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	agent, err := uc.agentRegistryRequester.GetAgent(ctx, request.CollectionAddress, request.ParsedCollectionTokenId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent from registry: %w", err)
	}
	if agent == nil {
		return nil, cb_errors.NewNotFoundError("agent", request.CollectionAddress+"/"+request.CollectionTokenId)
	}

	allowed := enums.ParticipantRoleFromAttributes(agent.Attributes).AllowedPromptTemplateKeys()
	if !allowed.Contains(request.ParsedKey()) {
		return nil, cb_errors.NewValidationError(
			"key",
			fmt.Sprintf("must be one of: %s for agent role %q", allowed.String(), agent.Attributes["role"]),
		)
	}

	template, err := uc.agentPromptTemplateRepository.GetByAgentAndKey(ctx,
		request.CollectionAddress,
		request.ParsedCollectionTokenId(),
		request.ParsedKey(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent prompt template: %w", err)
	}

	if template == nil {
		template = &models.AgentPromptTemplate{
			CollectionAddress: request.CollectionAddress,
			CollectionTokenId: request.ParsedCollectionTokenId(),
			Key:               request.ParsedKey(),
		}
	}
	template.Content = request.Content

	saved, err := uc.agentPromptTemplateRepository.Save(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed to save agent prompt template: %w", err)
	}

	return responses.AgentPromptTemplateResponseFromModel(saved), nil
}
