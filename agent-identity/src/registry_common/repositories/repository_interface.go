package repositories

import (
	"context"

	"gorm.io/gorm"
)

type RepositoryInterface[T any] interface {
	Save(ctx context.Context, entity *T) (*T, error)
	SaveUsingTransaction(ctx context.Context, entity *T, tx *gorm.DB) (*T, error)
}
