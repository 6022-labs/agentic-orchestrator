package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/google/uuid"
)

type AgentSearchParams struct {
	// IsMinted filters on mint state: true = minted only, false = proposals only, nil = all.
	IsMinted *bool
}

type AgentSearchParamsWithPagination struct {
	AgentSearchParams
	Page            *uint
	ElementsPerPage *uint
	SortBy          *string
	SortDirection   *string
}

type AgentRepositoryInterface interface {
	RepositoryInterface[models.Agent]
	GetAll(ctx context.Context) ([]*models.Agent, error)
	Search(ctx context.Context, params AgentSearchParamsWithPagination) ([]*models.Agent, error)
	SearchCount(ctx context.Context, params AgentSearchParams) (uint, error)
	Remove(ctx context.Context, agent *models.Agent) error
	GetById(ctx context.Context, id uuid.UUID) (*models.Agent, error)
	GetMintedByEnsDomain(ctx context.Context, ensDomain string) (*models.Agent, error)
	GetByEnsDomains(ctx context.Context, ensDomains []string) ([]*models.Agent, error)
	GetByCollectionAddressAndCollectionTokenId(ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
	) (*models.Agent, error)
	GetNotYetMintedByCollectionAddressAndProposalId(
		ctx context.Context,
		collectionAddress string,
		proposalId uint64,
	) (*models.Agent, error)
	ExistsByCollectionAddressAndCollectionTokenId(ctx context.Context, address string, tokenId uint64) (bool, error)
}
