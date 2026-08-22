package repositories

import (
	"context"

	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentRepository struct {
	*AbstractRepository[models.Agent]
	client db_clients.DatabaseClientInterface
}

func NewAgentRepository(client db_clients.DatabaseClientInterface) *AgentRepository {
	abstractRepository := NewAbstractRepository[models.Agent](client)

	return &AgentRepository{
		client:             client,
		AbstractRepository: abstractRepository,
	}
}

func (r *AgentRepository) Remove(ctx context.Context, agent *models.Agent) error {
	tx := r.client.GetDB().WithContext(ctx).Delete(agent)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *AgentRepository) GetAll(ctx context.Context) ([]*models.Agent, error) {
	agents := []*models.Agent{}

	tx := r.client.GetDB().WithContext(ctx).Find(&agents)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agents, nil
}

func (r *AgentRepository) Search(
	ctx context.Context,
	params repositories.AgentSearchParamsWithPagination,
) ([]*models.Agent, error) {
	agents := []*models.Agent{}

	db := r.applySearchParams(
		r.client.GetDB().WithContext(ctx).Model(&models.Agent{}),
		params.AgentSearchParams,
	)

	if params.ElementsPerPage != nil {
		offset := 0
		if params.Page != nil && *params.Page > 0 {
			offset = int((*params.Page - 1) * *params.ElementsPerPage)
		}
		db = db.Offset(offset).Limit(int(*params.ElementsPerPage))
	}

	sortDirection := "asc"
	if params.SortDirection != nil {
		sortDirection = *params.SortDirection
	}

	sortBy := "created_at"
	if params.SortBy != nil {
		sortBy = *params.SortBy
	}

	db = db.Order(sortBy + " " + sortDirection)

	tx := db.Find(&agents)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agents, nil
}

func (r *AgentRepository) SearchCount(ctx context.Context, params repositories.AgentSearchParams) (uint, error) {
	var count int64

	db := r.applySearchParams(r.client.GetDB().WithContext(ctx).Model(&models.Agent{}), params)

	tx := db.Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return uint(count), nil
}

// applySearchParams translates AgentSearchParams into query filters.
func (r *AgentRepository) applySearchParams(db *gorm.DB, params repositories.AgentSearchParams) *gorm.DB {
	if params.IsMinted != nil {
		if *params.IsMinted {
			db = db.Where("collection_token_id IS NOT NULL")
		} else {
			db = db.Where("collection_token_id IS NULL")
		}
	}
	return db
}

func (r *AgentRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Agent, error) {
	entity := new(models.Agent)
	tx := r.client.GetDB().WithContext(ctx).First(entity, id)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil
}

func (r *AgentRepository) GetByCollectionAddressAndCollectionTokenId(ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
) (*models.Agent, error) {
	agent := &models.Agent{}

	tx := r.client.GetDB().WithContext(ctx).
		Where("collection_address = ? AND collection_token_id = ?", collectionAddress, collectionTokenId).
		First(agent)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agent, nil
}

func (r *AgentRepository) GetMintedByEnsDomain(ctx context.Context, ensDomain string) (*models.Agent, error) {
	agent := &models.Agent{}

	tx := r.client.GetDB().
		WithContext(ctx).
		Where("ens_domain = ? AND collection_token_id IS NOT NULL", ensDomain).
		First(agent)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agent, nil
}

func (r *AgentRepository) GetByEnsDomains(ctx context.Context, ensDomains []string) ([]*models.Agent, error) {
	agents := []*models.Agent{}

	tx := r.client.GetDB().WithContext(ctx).Where("ens_domain IN ?", ensDomains).Find(&agents)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agents, nil
}

func (r *AgentRepository) GetNotYetMintedByCollectionAddressAndProposalId(ctx context.Context,
	collectionAddress string,
	proposalId uint64,
) (*models.Agent, error) {
	agent := &models.Agent{}

	tx := r.client.GetDB().WithContext(ctx).Where(
		"collection_token_id IS NULL AND collection_address = ? AND proposal_id = ?",
		collectionAddress,
		proposalId,
	).First(agent)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return agent, nil
}

func (r *AgentRepository) ExistsByCollectionAddressAndCollectionTokenId(
	ctx context.Context,
	address string,
	tokenId uint64,
) (bool, error) {
	var count int64

	tx := r.client.GetDB().WithContext(ctx).
		Model(&models.Agent{}).
		Where("collection_address = ? AND collection_token_id = ?", address, tokenId).
		Count(&count)
	if tx.Error != nil {
		return false, tx.Error
	}

	return count > 0, nil
}
