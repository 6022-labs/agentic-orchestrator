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
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	agent, err := uc.agentRegistryRequester.GetAgent(ctx, request.CollectionAddress, request.CollectionTokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent from registry: %w", err)
	}
	if agent == nil {
		agentKey := fmt.Sprintf("%s/%d", request.CollectionAddress, request.CollectionTokenId)
		return nil, cb_errors.NewNotFoundError("agent", agentKey)
	}

	allowed := enums.ParticipantRoleFromAttributes(agent.Attributes).AllowedPromptTemplateKeys()
	if !allowed.Contains(request.Key) {
		return nil, cb_errors.NewValidationError(
			"key",
			fmt.Sprintf("must be one of: %s for agent role %q", allowed.String(), agent.Attributes["role"]),
		)
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
		template = &models.AgentPromptTemplate{
			CollectionAddress: request.CollectionAddress,
			CollectionTokenId: request.CollectionTokenId,
			Key:               request.Key,
		}
	}
	template.Content = request.Content

	saved, err := uc.agentPromptTemplateRepository.Save(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed to save agent prompt template: %w", err)
	}

	return responses.AgentPromptTemplateResponseFromModel(saved), nil
}
