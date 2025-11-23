package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/google/uuid"
)

type AgentSearchParams struct {
	IsInternal *bool
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
	GetAll() ([]*models.Agent, error)
	Search(params AgentSearchParamsWithPagination) ([]*models.Agent, error)
	SearchCount(params AgentSearchParams) (uint, error)
	Remove(agent *models.Agent) error
	GetById(id uuid.UUID) (*models.Agent, error)
	GetByWalletAddress(walletAddress string) (*models.Agent, error)
	GetByEnsDomain(ensDomain string) (*models.Agent, error)
	GetByEnsDomains(ensDomains []string) ([]*models.Agent, error)
	GetByCollectionAddressAndCollectionTokenId(collectionAddress string, collectionTokenId uint64) (*models.Agent, error)
	GetNotYetMintedByCollectionAddressAndWalletAddress(collectionAddress string, walletAddress string) (*models.Agent, error)
	GetNotYetMintedByCollectionAddressAndProposalId(collectionAddress string, proposalId uint64) (*models.Agent, error)
	ExistsByCollectionAddressAndCollectionTokenId(address string, tokenId uint64) (bool, error)
}
