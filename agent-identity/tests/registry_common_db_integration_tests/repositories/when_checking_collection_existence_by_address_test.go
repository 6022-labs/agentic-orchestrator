package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCheckingCollectionExistenceByAddressTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentCollectionRepository
}

func WhenCheckingCollectionExistenceByAddressBeforeEach() *WhenCheckingCollectionExistenceByAddressTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := registry_common_db.NewDatabaseClient(
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

	return &WhenCheckingCollectionExistenceByAddressTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenCheckingCollectionExistenceByAddressAfterEach(suite *WhenCheckingCollectionExistenceByAddressTestingSuite) {
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
			createdCollection, err := suite.sut.Save(context.Background(), existingCollection)
			if err != nil {
				t.Errorf("Error creating collection: %v", err)
			}

			exists, err := suite.sut.ExistsByCollectionAddress(context.Background(), createdCollection.Address)

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

			exists, err := suite.sut.ExistsByCollectionAddress(context.Background(), nonExistentAddress)

			assert.Nil(t, err)
			assert.False(t, exists)
		})
	})
}
