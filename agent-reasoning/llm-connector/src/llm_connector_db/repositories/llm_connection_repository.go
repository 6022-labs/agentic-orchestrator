package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models/enums"
	"github.com/google/uuid"
)

type LLMConnectionRepository struct {
	client db_clients.DatabaseClientInterface
}

func NewLLMConnectionRepository(client db_clients.DatabaseClientInterface) *LLMConnectionRepository {
	return &LLMConnectionRepository{
		client: client,
	}
}

func (r *LLMConnectionRepository) Save(entity *models.LLMConnection) (*models.LLMConnection, error) {
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *LLMConnectionRepository) GetByAgentId(agentId uuid.UUID) (*models.LLMConnection, error) {
	var llmConnection models.LLMConnection
	if err := r.client.GetDB().Where("agent_id = ?", agentId).First(&llmConnection).Error; err != nil {
		return nil, err
	}

	return &llmConnection, nil
}

func (r *LLMConnectionRepository) GetSetupStatusByAgentId(agentId uuid.UUID) (*enums.LLMConnectionSetupStatus, error) {
	var setupStatus enums.LLMConnectionSetupStatus
	if err := r.client.GetDB().
		Model(&models.LLMConnection{}).
		Select("setup_status").
		Where("agent_id = ?", agentId).
		Limit(1).
		Scan(&setupStatus).
		Error; err != nil {
		return nil, err
	}

	if !enums.GetAllLLMConnectionSetupStatuses().Contains(setupStatus) {
		return nil, nil
	}

	return setupStatus.ToPointer(), nil
}

func (r *LLMConnectionRepository) DeleteByAgentId(agentId uuid.UUID) error {
	var llmConnection models.LLMConnection
	if err := r.client.GetDB().Where("agent_id = ?", agentId).Delete(&llmConnection).Error; err != nil {
		return err
	}

	return nil
}
