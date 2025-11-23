package factories_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/factories/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/database_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type WhenCreatingAgentDbClientTestCase struct {
	name             string
	clientBuilderFn  func() container_builders.DatabaseContainerBuilderInterface
	createDatabaseFn func(databaseConnection *models.DatabaseConnection, connectionString string) error
	driver           enums.DatabaseConnectionDriver
}

var (
	WhenCreatingAgentDbClientTestCases = []WhenCreatingAgentDbClientTestCase{
		{
			name:   "PostgreSQL",
			driver: enums.DatabaseConnectionDriverPostgres,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetPostgresContainerBuilder()
			},
			createDatabaseFn: createPostgresAgentDatabase,
		},
		{
			name:   "MariaDB",
			driver: enums.DatabaseConnectionDriverMariaDB,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMariaDBContainerBuilder()
			},
			createDatabaseFn: createMariaDBAgentDatabase,
		},
		{
			name:   "MySQL",
			driver: enums.DatabaseConnectionDriverMySQL,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMySQLContainerBuilder()
			},
			createDatabaseFn: createMySQLAgentDatabase,
		},
	}
)

type WhenCreatingAgentDbClientTestingSuite struct {
	dbName                   string
	connectionString         string
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface

	sut *db_clients.AgentDatabaseClientFactory

	migrationsSettings                        *settings.MigrationsSettings
	mockDatabaseConnectionRepository          *repositories_mocks.MockDatabaseConnectionRepositoryInterface
	mockAgentDatabaseCredentialsStorage       *services_mocks.MockAgentDatabaseCredentialsStorageInterface
	mockDatabaseConnectionMigrationRepository *repositories_mocks.MockDatabaseConnectionMigrationRepositoryInterface
}

func WhenCreatingAgentDbClientBeforeEach(t *testing.T, clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface) *WhenCreatingAgentDbClientTestingSuite {
	databaseContainerBuilder := clientBuilderFn()
	dbName, connectionString := databaseContainerBuilder.Create()

	mockController := gomock.NewController(t)

	mockDatabaseConnectionRepository := repositories_mocks.NewMockDatabaseConnectionRepositoryInterface(mockController)
	mockAgentDatabaseCredentialsStorage := services_mocks.NewMockAgentDatabaseCredentialsStorageInterface(mockController)
	mockDatabaseConnectionMigrationRepository := repositories_mocks.NewMockDatabaseConnectionMigrationRepositoryInterface(mockController)

	migrationFolderPath, err := common_tests.FindAgentMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: *migrationFolderPath,
	}

	sut := db_clients.NewAgentDatabaseClientFactory(
		zap.NewNop(),
		migrationsSettings,
		mockDatabaseConnectionRepository,
		mockAgentDatabaseCredentialsStorage,
		mockDatabaseConnectionMigrationRepository,
	)

	return &WhenCreatingAgentDbClientTestingSuite{
		sut: sut,

		dbName:                   dbName,
		connectionString:         connectionString,
		databaseContainerBuilder: databaseContainerBuilder,

		migrationsSettings: migrationsSettings,

		mockDatabaseConnectionRepository:          mockDatabaseConnectionRepository,
		mockAgentDatabaseCredentialsStorage:       mockAgentDatabaseCredentialsStorage,
		mockDatabaseConnectionMigrationRepository: mockDatabaseConnectionMigrationRepository,
	}
}

func WhenCreatingAgentDbClientAfterEach(suite *WhenCreatingAgentDbClientTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func createPostgresAgentDatabase(databaseConnection *models.DatabaseConnection, connectionString string) error {
	// Create the database and role for the agent
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB: %w", err)
	}
	defer sqlDB.Close()

	createDatabaseQuery := fmt.Sprintf(`CREATE DATABASE "%s";`, databaseConnection.AgentID)
	if err := db.Exec(createDatabaseQuery).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	createRoleQuery := fmt.Sprintf(`CREATE ROLE "%s" WITH LOGIN PASSWORD '%s';`, databaseConnection.Username, databaseConnection.AgentID)
	if err := db.Exec(createRoleQuery).Error; err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	grantPrivilegesQuery := fmt.Sprintf(`GRANT ALL PRIVILEGES ON DATABASE "%s" TO "%s";`, databaseConnection.AgentID, databaseConnection.Username)
	if err := db.Exec(grantPrivilegesQuery).Error; err != nil {
		return fmt.Errorf("failed to grant privileges on database: %w", err)
	}

	re := regexp.MustCompile(`(dbname=)([a-fA-F0-9-]+)`)
	newDatabaseConnectionString := re.ReplaceAllString(connectionString, fmt.Sprintf("${1}%s", databaseConnection.AgentID))

	db, err = gorm.Open(postgres.Open(newDatabaseConnectionString), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to new database: %w", err)
	}

	grantSchemaPrivilegesQuery := fmt.Sprintf(`GRANT ALL PRIVILEGES ON SCHEMA public TO "%s";`, databaseConnection.Username)
	if err := db.Exec(grantSchemaPrivilegesQuery).Error; err != nil {
		return fmt.Errorf("failed to grant privileges on schema: %w", err)
	}

	return nil
}

func createMySQLAgentDatabase(databaseConnection *models.DatabaseConnection, connectionString string) error {
	// Create the database and user for the agent
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB: %w", err)
	}
	defer sqlDB.Close()

	createDatabaseQuery := fmt.Sprintf("CREATE DATABASE `%s`;", databaseConnection.AgentID)
	if err := db.Exec(createDatabaseQuery).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	createUserQuery := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s';", databaseConnection.Username, databaseConnection.AgentID)
	if err := db.Exec(createUserQuery).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	grantPrivilegesQuery := fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%%';", databaseConnection.AgentID, databaseConnection.Username)
	if err := db.Exec(grantPrivilegesQuery).Error; err != nil {
		return fmt.Errorf("failed to grant privileges: %w", err)
	}

	flushPrivilegesQuery := "FLUSH PRIVILEGES;"
	if err := db.Exec(flushPrivilegesQuery).Error; err != nil {
		return fmt.Errorf("failed to flush privileges: %w", err)
	}

	return nil
}

func createMariaDBAgentDatabase(databaseConnection *models.DatabaseConnection, connectionString string) error {
	return createMySQLAgentDatabase(databaseConnection, connectionString) // MariaDB is compatible with MySQL syntax
}

func TestWhenCreatingAgentDbClient(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenCreatingAgentDbClientTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given database connection not found", func(t *testing.T) {
				t.Parallel()

				initSuite := func(suite *WhenCreatingAgentDbClientTestingSuite, agentId uuid.UUID) {
					suite.mockDatabaseConnectionRepository.EXPECT().
						GetByAgentId(gomock.Any()).
						Return(nil, nil)
				}

				t.Run("Should return error", func(t *testing.T) {

					t.Parallel()

					suite := WhenCreatingAgentDbClientBeforeEach(t, testCase.clientBuilderFn)
					defer WhenCreatingAgentDbClientAfterEach(suite)

					agentId := uuid.New()

					initSuite(suite, agentId)

					_, err := suite.sut.Create(agentId)

					assert.Error(t, err)
				})
			})

			t.Run("Given error occurs while getting agent database credentials", func(t *testing.T) {
				t.Parallel()

				initSuite := func(suite *WhenCreatingAgentDbClientTestingSuite, agentId uuid.UUID) {
					databaseConnection := model_builders.NewDatabaseConnectionBuilder().
						WithAgentId(agentId).
						Build()

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetByAgentId(agentId).
						Return(databaseConnection, nil)

					suite.mockAgentDatabaseCredentialsStorage.EXPECT().
						Get(agentId).
						Return(nil, assert.AnError)
				}

				t.Run("Should return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingAgentDbClientBeforeEach(t, testCase.clientBuilderFn)
					defer WhenCreatingAgentDbClientAfterEach(suite)

					agentId := uuid.New()

					initSuite(suite, agentId)

					_, err := suite.sut.Create(agentId)

					assert.Error(t, err)
				})
			})

			t.Run("Given agent database credentials not found", func(t *testing.T) {
				t.Parallel()

				initSuite := func(suite *WhenCreatingAgentDbClientTestingSuite, agentId uuid.UUID) {
					databaseConnection := model_builders.NewDatabaseConnectionBuilder().
						WithAgentId(agentId).
						Build()

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetByAgentId(agentId).
						Return(databaseConnection, nil)

					suite.mockAgentDatabaseCredentialsStorage.EXPECT().
						Get(agentId).
						Return(nil, nil)
				}

				t.Run("Should return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingAgentDbClientBeforeEach(t, testCase.clientBuilderFn)
					defer WhenCreatingAgentDbClientAfterEach(suite)

					agentId := uuid.New()

					initSuite(suite, agentId)

					_, err := suite.sut.Create(agentId)

					assert.Error(t, err)
				})
			})

			t.Run("Given error occurs while fetching agent migrations in orchestrator database", func(t *testing.T) {
				t.Parallel()

				initSuite := func(suite *WhenCreatingAgentDbClientTestingSuite, agentId uuid.UUID) {
					databaseConnection := model_builders.NewDatabaseConnectionBuilder().
						WithAgentId(agentId).
						WithHost(suite.databaseContainerBuilder.GetHost()).
						WithPort(uint16(suite.databaseContainerBuilder.GetPort())).
						WithDriver(testCase.driver).
						Build()

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetByAgentId(agentId).
						Return(databaseConnection, nil)

					password := agentId.String()
					suite.mockAgentDatabaseCredentialsStorage.EXPECT().
						Get(agentId).
						Return(&password, nil)

					suite.mockDatabaseConnectionMigrationRepository.EXPECT().
						GetAll().
						Return(nil, assert.AnError)
				}

				t.Run("Should return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingAgentDbClientBeforeEach(t, testCase.clientBuilderFn)
					defer WhenCreatingAgentDbClientAfterEach(suite)

					agentId := uuid.New()

					initSuite(suite, agentId)

					_, err := suite.sut.Create(agentId)

					assert.Error(t, err)
				})
			})

			t.Run("Given error occurs while saving migrations in orchestrator database", func(t *testing.T) {
				t.Parallel()

				initSuite := func(suite *WhenCreatingAgentDbClientTestingSuite, agentId uuid.UUID) {
					databaseConnection := model_builders.NewDatabaseConnectionBuilder().
						WithAgentId(agentId).
						WithHost(suite.databaseContainerBuilder.GetHost()).
						WithPort(uint16(suite.databaseContainerBuilder.GetPort())).
						WithDriver(testCase.driver).
						Build()

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetByAgentId(agentId).
						Return(databaseConnection, nil)

					password := agentId.String()
					suite.mockAgentDatabaseCredentialsStorage.EXPECT().
						Get(agentId).
						Return(&password, nil)

					suite.mockDatabaseConnectionMigrationRepository.EXPECT().
						GetAll().
						Return([]*models.DatabaseConnectionMigration{
							{ID: uuid.New(), Name: "20250729144720_initial-migration.sql"},
						}, nil)

					suite.mockDatabaseConnectionMigrationRepository.EXPECT().
						Save(gomock.Any()).
						Return(nil, assert.AnError)
				}

				t.Run("Should return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingAgentDbClientBeforeEach(t, testCase.clientBuilderFn)
					defer WhenCreatingAgentDbClientAfterEach(suite)

					agentId := uuid.New()

					initSuite(suite, agentId)

					_, err := suite.sut.Create(agentId)

					assert.Error(t, err)
				})
			})

			t.Run("Given error occurs while getting migrated migrations", func(t *testing.T) {
				t.Parallel()

				initSuite := func(suite *WhenCreatingAgentDbClientTestingSuite, agentId uuid.UUID) {
					databaseConnection := model_builders.NewDatabaseConnectionBuilder().
						WithAgentId(agentId).
						WithHost(suite.databaseContainerBuilder.GetHost()).
						WithPort(uint16(suite.databaseContainerBuilder.GetPort())).
						WithDriver(testCase.driver).
						Build()

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetByAgentId(agentId).
						Return(databaseConnection, nil)

					password := agentId.String()
					suite.mockAgentDatabaseCredentialsStorage.EXPECT().
						Get(agentId).
						Return(&password, nil)

					suite.mockDatabaseConnectionMigrationRepository.EXPECT().
						GetAll().
						Return([]*models.DatabaseConnectionMigration{
							{ID: uuid.New(), Name: "20250729144720_initial-migration.sql"},
						}, nil)

					databaseConnectionMigration := model_builders.NewDatabaseConnectionMigrationBuilder().Build()
					suite.mockDatabaseConnectionMigrationRepository.EXPECT().
						Save(gomock.Any()).
						Return(databaseConnectionMigration, nil).
						AnyTimes()

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetMigratedMigrations(databaseConnection.ID).
						Return(nil, assert.AnError)
				}

				t.Run("Should return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingAgentDbClientBeforeEach(t, testCase.clientBuilderFn)
					defer WhenCreatingAgentDbClientAfterEach(suite)

					agentId := uuid.New()

					initSuite(suite, agentId)

					_, err := suite.sut.Create(agentId)

					assert.Error(t, err)
				})
			})

			t.Run("Given successful connection to the database", func(t *testing.T) {
				t.Parallel()

				initSuite := func(suite *WhenCreatingAgentDbClientTestingSuite, agentId uuid.UUID) {
					databaseConnection := model_builders.NewDatabaseConnectionBuilder().
						WithAgentId(agentId).
						WithHost(suite.databaseContainerBuilder.GetHost()).
						WithPort(uint16(suite.databaseContainerBuilder.GetPort())).
						WithDriver(testCase.driver).
						Build()

					err := testCase.createDatabaseFn(databaseConnection, suite.connectionString)
					if err != nil {
						panic(fmt.Sprintf("failed to create database: %v", err))
					}

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetByAgentId(agentId).
						Return(databaseConnection, nil)

					password := agentId.String()
					suite.mockAgentDatabaseCredentialsStorage.EXPECT().
						Get(agentId).
						Return(&password, nil)

					suite.mockDatabaseConnectionMigrationRepository.EXPECT().
						GetAll().
						Return([]*models.DatabaseConnectionMigration{
							{ID: uuid.New(), Name: "20250729144720_initial-migration.sql"},
						}, nil)

					databaseConnectionMigration := model_builders.NewDatabaseConnectionMigrationBuilder().Build()
					suite.mockDatabaseConnectionMigrationRepository.EXPECT().
						Save(gomock.Any()).
						Return(databaseConnectionMigration, nil).
						AnyTimes()

					suite.mockDatabaseConnectionRepository.EXPECT().
						GetMigratedMigrations(databaseConnection.ID).
						Return([]string{"20250729144720_initial-migration.sql"}, nil)
				}

				t.Run("Should return a new agent db client", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingAgentDbClientBeforeEach(t, testCase.clientBuilderFn)
					defer WhenCreatingAgentDbClientAfterEach(suite)

					agentId := uuid.New()
					initSuite(suite, agentId)

					dbClient, err := suite.sut.Create(agentId)

					assert.Nil(t, err)
					assert.NotNil(t, dbClient)
				})
			})
		})
	}
}
