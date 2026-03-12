package repositories

import (
	"gorm.io/gorm"
)

type RepositoryInterface[T any] interface {
	Save(entity *T) (*T, error)
	SaveUsingTransaction(entity *T, tx *gorm.DB) (*T, error)
}
