package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentAddressRepository struct {
	*AbstractRepository[models.AgentAddress]
	client db_clients.DatabaseClientInterface
}

func NewAgentAddressRepository(client db_clients.DatabaseClientInterface) *AgentAddressRepository {
	abstractRepository := NewAbstractRepository[models.AgentAddress](client)

	return &AgentAddressRepository{
		client:             client,
		AbstractRepository: abstractRepository,
	}
}

func (r *AgentAddressRepository) GetByTypeAndAddress(addressType enums.AddressType, address string) (*models.AgentAddress, error) {
	agentAddress := &models.AgentAddress{}

	tx := r.client.GetDB().Where("type = ? AND address = ?", addressType, address).First(agentAddress)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agentAddress, nil
}

func (r *AgentAddressRepository) GetAllByAgentId(agentId uuid.UUID) ([]*models.AgentAddress, error) {
	agentAddresses := []*models.AgentAddress{}

	tx := r.client.GetDB().Where("agent_id = ?", agentId).Find(&agentAddresses)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agentAddresses, nil
}

func (r *AgentAddressRepository) GetAllByAgentIds(agentIds []uuid.UUID) ([]*models.AgentAddress, error) {
	agentAddresses := []*models.AgentAddress{}

	tx := r.client.GetDB().Where("agent_id IN ?", agentIds).Find(&agentAddresses)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agentAddresses, nil
}

func (r *AgentAddressRepository) ReplaceAllForAgent(agentId uuid.UUID, addresses map[enums.AddressType]string) error {
	return r.client.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("agent_id = ?", agentId).Delete(&models.AgentAddress{}).Error; err != nil {
			return err
		}
		for addressType, address := range addresses {
			row := models.AgentAddress{AgentId: agentId, Type: addressType, Address: address}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
