package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/db_clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_db/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCheckingCollectionExistenceByAddressTestSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentCollectionRepository
}

func WhenCheckingCollectionExistenceByAddressBeforeEach() *WhenCheckingCollectionExistenceByAddressTestSuite {
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

	sut := repositories.NewAgentCollectionRepository(dbClient)

	return &WhenCheckingCollectionExistenceByAddressTestSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenCheckingCollectionExistenceByAddressAfterEach(suite *WhenCheckingCollectionExistenceByAddressTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCheckingCollectionExistenceByAddress(t *testing.T) {
	t.Parallel()

	t.Run("Given collection exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return true", func(t *testing.T) {
			suite := WhenCheckingCollectionExistenceByAddressBeforeEach()
			defer WhenCheckingCollectionExistenceByAddressAfterEach(suite)

			existingCollection := model_builders.NewAgentCollectionBuilder().Build()
			createdCollection, err := suite.sut.Save(existingCollection)
			if err != nil {
				t.Errorf("Error creating collection: %v", err)
			}

			exists, err := suite.sut.ExistsByCollectionAddress(createdCollection.Address)

			assert.Nil(t, err)
			assert.True(t, exists)
		})
	})

	t.Run("Given collection does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return false", func(t *testing.T) {
			suite := WhenCheckingCollectionExistenceByAddressBeforeEach()
			defer WhenCheckingCollectionExistenceByAddressAfterEach(suite)

			nonExistentAddress := "0x000000000000000000000000000000000000dead"

			exists, err := suite.sut.ExistsByCollectionAddress(nonExistentAddress)

			assert.Nil(t, err)
			assert.False(t, exists)
		})
	})
}
