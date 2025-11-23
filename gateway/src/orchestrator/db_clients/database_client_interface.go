package db_clients

import "gorm.io/gorm"

type DatabaseClientInterface interface {
	Migrate() error
	GetDB() *gorm.DB
	HealthCheck() error
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
}
