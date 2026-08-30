package db_clients

import "gorm.io/gorm"

type DatabaseClientInterface interface {
	Migrate() error
	// Close releases the connection pool; tests need it, a host does not.
	Close() error
	GetDB() *gorm.DB
	HealthCheck() error
	BeginTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB)
	RollbackTransaction(tx *gorm.DB)
}
