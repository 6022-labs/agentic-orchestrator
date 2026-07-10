package use_cases

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases/responses"
	"gorm.io/gorm"
)

type SaveAgentPromptTemplateInterface interface {
	Save(request *requests.SaveAgentPromptTemplateRequest) (*responses.AgentPromptTemplateResponse, error)
}

type SaveAgentPromptTemplate struct {
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface
}

func NewSaveAgentPromptTemplate(
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface,
) *SaveAgentPromptTemplate {
	return &SaveAgentPromptTemplate{agentPromptTemplateRepository: agentPromptTemplateRepository}
}

func (uc *SaveAgentPromptTemplate) Save(request *requests.SaveAgentPromptTemplateRequest) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	template, err := uc.agentPromptTemplateRepository.GetByAgentAndKey(
		request.CollectionAddress,
		request.ParsedCollectionTokenId(),
		request.ParsedKey(),
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	saved, err := uc.agentPromptTemplateRepository.Save(template)
	if err != nil {
		return nil, fmt.Errorf("failed to save agent prompt template: %w", err)
	}

	return responses.AgentPromptTemplateResponseFromModel(saved), nil
}
