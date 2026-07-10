package prompts

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"gorm.io/gorm"
)

type promptTemplateProvider struct {
	agentPromptTemplateRepository repositories.AgentPromptTemplateRepositoryInterface
}

func (p *promptTemplateProvider) templateFor(agent *models.Participant, key enums.AgentPromptTemplateKey) (string, error) {
	if agent.CollectionAddress != nil && agent.CollectionTokenId != nil {
		stored, err := p.agentPromptTemplateRepository.GetByAgentAndKey(*agent.CollectionAddress, *agent.CollectionTokenId, key)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("failed to load prompt template %s: %w", key, err)
		}
		if stored != nil {
			return stored.Content, nil
		}
	}

	tpl, ok := p.defaultTemplateForKey(key)
	if !ok {
		return "", fmt.Errorf("no default template for key %s", key)
	}
	return tpl, nil
}
