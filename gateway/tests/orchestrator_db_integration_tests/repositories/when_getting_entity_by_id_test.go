package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db"
	db_repositories "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

type WhenGettingEntityByIdTestSuite struct {
	dbName string

	sut repositories.AgentRepositoryInterface
}

func WhenGettingEntityByIdBeforeEach() *WhenGettingEntityByIdTestSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := orchestrator_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := db_repositories.NewAgentRepository(dbClient)

	return &WhenGettingEntityByIdTestSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenGettingEntityByIdBeforeEachAfterEach(suite *WhenGettingEntityByIdTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingEntityById(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingEntityByIdBeforeEach()
			defer WhenGettingEntityByIdBeforeEachAfterEach(suite)

			entity := model_builders.NewAgentBuilder().Build()
			createdEntity, _ := suite.sut.Save(entity)

			fetchedEntity, err := suite.sut.GetById(createdEntity.ID)

			assert.Nil(t, err)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
		})
	})

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingEntityByIdBeforeEach()
			defer WhenGettingEntityByIdBeforeEachAfterEach(suite)

			nonExistentID := uuid.New()
			_, err := suite.sut.GetById(nonExistentID)

			assert.NotNil(t, err)
		})
	})
}
