package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	orchestrator_agent_database_connector_db_orchestrator "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenSavingDatabaseConnectionTestingSuite struct {
	sut *repositories.DatabaseConnectionRepository

	dbName                   string
	dbClient                 db_clients.OrchestratorDatabaseClientInterface
	postgresContainerBuilder *container_builders.PostgresContainerBuilder
}

func WhenSavingDatabaseConnectionBeforeEach() *WhenSavingDatabaseConnectionTestingSuite {
	postgresContainerBuilder := container_builders.GetPostgresContainerBuilder()
	dbName, connectionString := postgresContainerBuilder.Create()

	migrationFolderPath, err := common_tests.FindOrchestratorMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, err := orchestrator_agent_database_connector_db_orchestrator.NewOrchestratorDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err != nil {
		panic(err)
	}

	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewDatabaseConnectionRepository(dbClient)

	return &WhenSavingDatabaseConnectionTestingSuite{
		sut: sut,

		dbName:                   dbName,
		dbClient:                 dbClient,
		postgresContainerBuilder: postgresContainerBuilder,
	}
}

func WhenSavingDatabaseConnectionAfterEach(suite *WhenSavingDatabaseConnectionTestingSuite) {
	suite.postgresContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenSavingDatabaseConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingDatabaseConnectionBeforeEach()
			defer WhenSavingDatabaseConnectionAfterEach(suite)

			entityToCreate := model_builders.NewDatabaseConnectionBuilder().Build()
			createdEntity, err := suite.sut.Save(entityToCreate)

			assert.Nil(t, err)
			assert.Equal(t, entityToCreate, createdEntity)
		})

		t.Run("Should store entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingDatabaseConnectionBeforeEach()
			defer WhenSavingDatabaseConnectionAfterEach(suite)

			entityToCreate := model_builders.NewDatabaseConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

			assert.Equal(t, entityToCreate, createdEntity)
			assert.NotEqual(t, uuid.Nil, createdEntity.ID)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return updated entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingDatabaseConnectionBeforeEach()
			defer WhenSavingDatabaseConnectionAfterEach(suite)

			entityToCreate := model_builders.NewDatabaseConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			entityToUpdate := model_builders.NewDatabaseConnectionBuilder().WithID(createdEntity.ID).Build()
			updatedEntity, err := suite.sut.Save(entityToUpdate)

			assert.Nil(t, err)
			assert.Equal(t, entityToUpdate, updatedEntity)
		})

		t.Run("Should update entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingDatabaseConnectionBeforeEach()
			defer WhenSavingDatabaseConnectionAfterEach(suite)

			entityToCreate := model_builders.NewDatabaseConnectionBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			entityToUpdate := model_builders.NewDatabaseConnectionBuilder().WithID(createdEntity.ID).Build()
			suite.sut.Save(entityToUpdate)

			suite.dbClient.GetDB().First(&entityToCreate, createdEntity.ID)

			assert.Equal(t, entityToCreate.ID, entityToUpdate.ID)
			assert.Equal(t, entityToCreate.Host, entityToUpdate.Host)
			assert.Equal(t, entityToCreate.Port, entityToUpdate.Port)
			assert.Equal(t, entityToCreate.Driver, entityToUpdate.Driver)
		})
	})
}
