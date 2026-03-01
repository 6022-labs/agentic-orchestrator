package db_clients

import (
	"fmt"
	"math"
	"path/filepath"
	"sync"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/utils"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AgentDatabaseClientFactory struct {
	logger                                *zap.Logger
	migrationsSettings                    *settings.MigrationsSettings
	agentDatabaseCredentialsStorage       services.AgentDatabaseCredentialsStorageInterface
	databaseConnectionRepository          repositories.DatabaseConnectionRepositoryInterface
	databaseConnectionMigrationRepository repositories.DatabaseConnectionMigrationRepositoryInterface

	migrationMutex       *sync.RWMutex
	migrationsMutex      map[uuid.UUID]*sync.RWMutex // Mutex for each agent migrations
	migrationsInFolder   map[enums.DatabaseConnectionDriver][]string
	migrationsInDatabase []uuid.UUID
}

func NewAgentDatabaseClientFactory(
	logger *zap.Logger,
	migrationsSettings *settings.MigrationsSettings,
	databaseConnectionRepository repositories.DatabaseConnectionRepositoryInterface,
	agentDatabaseCredentialsStorage services.AgentDatabaseCredentialsStorageInterface,
	databaseConnectionMigrationRepository repositories.DatabaseConnectionMigrationRepositoryInterface,
) *AgentDatabaseClientFactory {
	return &AgentDatabaseClientFactory{
		migrationMutex:       &sync.RWMutex{},
		migrationsMutex:      make(map[uuid.UUID]*sync.RWMutex),
		migrationsInFolder:   make(map[enums.DatabaseConnectionDriver][]string),
		migrationsInDatabase: make([]uuid.UUID, 0),

		logger:                                logger,
		migrationsSettings:                    migrationsSettings,
		databaseConnectionRepository:          databaseConnectionRepository,
		agentDatabaseCredentialsStorage:       agentDatabaseCredentialsStorage,
		databaseConnectionMigrationRepository: databaseConnectionMigrationRepository,
	}
}

func (factory *AgentDatabaseClientFactory) Create(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, error) {
	databaseConnection, err := factory.databaseConnectionRepository.GetByAgentId(agentId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get database connection: %w for agent ID: %s", err, agentId.String())
	}
	if databaseConnection == nil {
		return nil, fmt.Errorf("database connection not found for agent ID: %s", agentId.String())
	}

	return factory.CreateFromDatabaseConnection(databaseConnection)
}

func (factory *AgentDatabaseClientFactory) CreateFromDatabaseConnection(databaseConnection *models.DatabaseConnection) (db_clients.AgentDatabaseClientInterface, error) {
	agentDatabaseCredentials, err := factory.agentDatabaseCredentialsStorage.Get(databaseConnection.AgentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent database credentials: %w", err)
	}
	if agentDatabaseCredentials == nil {
		return nil, fmt.Errorf("agent database credentials not found")
	}

	var connectionString string

	switch databaseConnection.Driver {
	case enums.DatabaseConnectionDriverPostgres:
		connectionString = utils.BuildPostgresqlConnectionString(
			databaseConnection.Host,
			databaseConnection.AgentID.String(),
			databaseConnection.Username,
			*agentDatabaseCredentials,
			databaseConnection.Port,
		)
	case enums.DatabaseConnectionDriverMariaDB, enums.DatabaseConnectionDriverMySQL:
		connectionString = utils.BuildMySqlConnectionString(
			databaseConnection.Host,
			databaseConnection.AgentID.String(),
			databaseConnection.Username,
			*agentDatabaseCredentials,
			databaseConnection.Port,
		)
	default:
		return nil, fmt.Errorf(
			"unsupported database driver: %s",
			databaseConnection.Driver,
		)
	}

	factory.migrationMutex.Lock()
	if factory.migrationsMutex[databaseConnection.ID] == nil {
		factory.migrationsMutex[databaseConnection.ID] = &sync.RWMutex{}
	}
	factory.migrationMutex.Unlock()

	factory.migrationsMutex[databaseConnection.ID].Lock()
	defer factory.migrationsMutex[databaseConnection.ID].Unlock()

	if factory.migrationsInFolder[databaseConnection.Driver] == nil {
		err = factory.syncMigrationsInDatabaseWithMigrationsInFolder(databaseConnection.Driver)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch available migrations: %w", err)
		}
	}

	migratedMigrations, err := factory.databaseConnectionRepository.GetMigratedMigrations(databaseConnection.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get migrated migrations: %w", err)
	}

	client, err := database_connector_db_agent.NewAgentDbClient(
		connectionString,
		factory.logger,
		databaseConnection.Driver,
		factory.migrationsSettings,
	)
	if err != nil {
		return nil, err
	}

	if !lo.ElementsMatch(factory.migrationsInFolder[databaseConnection.Driver], migratedMigrations) {
		factory.logger.Info(
			"Migrating agent database",
			zap.String("agent_id", databaseConnection.AgentID.String()),
		)

		if err := client.Migrate(); err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
	}

	return client, nil
}

func (factory *AgentDatabaseClientFactory) syncMigrationsInDatabaseWithMigrationsInFolder(driver enums.DatabaseConnectionDriver) error {
	if factory.migrationsInFolder[driver] == nil {
		availableMigrations, err := factory.fetchMigrationsInFolder(driver)
		if err != nil {
			return fmt.Errorf("failed to fetch available migrations: %w", err)
		}

		factory.migrationsInFolder[driver] = availableMigrations
	}

	migrationsInDatabase, err := factory.databaseConnectionMigrationRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get migrations in database: %w", err)
	}

	factory.migrationsInDatabase = lo.Map(migrationsInDatabase, func(migration *models.DatabaseConnectionMigration, _ int) uuid.UUID {
		return migration.ID
	})

	for _, migrationInFolder := range factory.migrationsInFolder[driver] {
		_, found := lo.Find(migrationsInDatabase, func(migration *models.DatabaseConnectionMigration) bool {
			return migration.Name == migrationInFolder
		})

		if !found {
			newDatabaseConnectionMigration := &models.DatabaseConnectionMigration{
				ID:   uuid.New(),
				Name: migrationInFolder,
			}

			newDatabaseConnectionMigration, err = factory.databaseConnectionMigrationRepository.Save(newDatabaseConnectionMigration)
			if err != nil {
				return fmt.Errorf("failed to save new migration %s: %w", migrationInFolder, err)
			}

			factory.migrationsInDatabase = append(factory.migrationsInDatabase, newDatabaseConnectionMigration.ID)
		}
	}

	return nil
}

func (factory *AgentDatabaseClientFactory) fetchMigrationsInFolder(driver enums.DatabaseConnectionDriver) ([]string, error) {
	migrations, err := goose.CollectMigrations(utils.GetMigrationFolderPath(factory.migrationsSettings.MigrationsFolderPath, driver.String()), 0, math.MaxInt64)
	if err != nil {
		return nil, fmt.Errorf("failed to collect migrations for driver %s: %w", driver.String(), err)
	}

	result := make([]string, 0, migrations.Len())
	for _, migration := range migrations {
		result = append(result, filepath.Base(migration.String()))
	}

	return result, nil
}
