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
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingMintedAgentsByEnsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentRepository
}

func WhenGettingMintedAgentsByEnsBeforeEach() *WhenGettingMintedAgentsByEnsTestingSuite {
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

	sut := repositories.NewAgentRepository(dbClient)

	return &WhenGettingMintedAgentsByEnsTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingMintedAgentsByEnsAfterEach(suite *WhenGettingMintedAgentsByEnsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingMintedAgentsByEns(t *testing.T) {
	t.Parallel()

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMintedAgentsByEnsBeforeEach()
			defer WhenGettingMintedAgentsByEnsAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().
				WithCollectionAddress(common_tests.RandomEthAddress().Hex()).
				WithCollectionTokenId(lo.ToPtr(uint64(1))).
				Build()
			createdEntity, err := suite.sut.Save(context.Background(), existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntity, err := suite.sut.GetMintedByEnsDomain(context.Background(), createdEntity.EnsDomain)

			assert.Nil(t, err)
			assert.NotNil(t, fetchedEntity)
			assert.Equal(t, createdEntity.ID, fetchedEntity.ID)
			assert.Equal(t, createdEntity.Name, fetchedEntity.Name)
			assert.Equal(t, createdEntity.Images, fetchedEntity.Images)
			assert.Equal(t, createdEntity.Creator, fetchedEntity.Creator)
			assert.Equal(t, createdEntity.EnsDomain, fetchedEntity.EnsDomain)
			assert.Equal(t, createdEntity.Attributes, fetchedEntity.Attributes)
		})
	})

	t.Run("Given entity exists but not minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMintedAgentsByEnsBeforeEach()
			defer WhenGettingMintedAgentsByEnsAfterEach(suite)

			existingEntity := model_builders.NewAgentBuilder().
				WithCollectionAddress(common_tests.RandomEthAddress().Hex()).
				WithCollectionTokenId(nil).
				Build()
			createdEntity, err := suite.sut.Save(context.Background(), existingEntity)
			if err != nil {
				t.Errorf("Error creating entity: %v", err)
			}

			fetchedEntity, err := suite.sut.GetMintedByEnsDomain(context.Background(), createdEntity.EnsDomain)

			assert.Nil(t, err)
			assert.Nil(t, fetchedEntity)
		})
	})

	t.Run("Given no entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingMintedAgentsByEnsBeforeEach()
			defer WhenGettingMintedAgentsByEnsAfterEach(suite)

			fetchedEntity, err := suite.sut.GetMintedByEnsDomain(context.Background(), "nonexistent.ens")

			assert.Nil(t, err)
			assert.Nil(t, fetchedEntity)
		})
	})
}
