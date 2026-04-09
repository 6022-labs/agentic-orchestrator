package db_clients

import "gorm.io/gorm"

type OrchestratorDatabaseClientInterface interface {
	Migrate() error
	GetDB() *gorm.DB
	HealthCheck() error
	BeginTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB)
	RollbackTransaction(tx *gorm.DB)
}
