package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenSavingEntityTestSuite struct {
	dbName string

	dbClient *slack_bridge_db.DatabaseClient
	sut      *repositories.SlackChannelRepository
}

func WhenSavingEntityBeforeEach() *WhenSavingEntityTestSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewSlackChannelRepository(dbClient)

	return &WhenSavingEntityTestSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenSavingEntityAfterEach(suite *WhenSavingEntityTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingEntity(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach()
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, err := suite.sut.Save(entityToCreate)

			assert.Nil(t, err)
			assert.Equal(t, entityToCreate, createdEntity)
		})

		t.Run("Should store entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach()
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			suite.dbClient.GetDB().First(entityToCreate, "id = ?", createdEntity.Id)

			assert.Equal(t, entityToCreate, createdEntity)
			assert.NotEqual(t, uuid.Nil, createdEntity.Id)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return updated entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach()
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			entityToUpdate := model_builders.NewSlackChannelBuilder().WithId(createdEntity.Id).Build()
			updatedEntity, err := suite.sut.Save(entityToUpdate)

			assert.Nil(t, err)
			assert.Equal(t, entityToUpdate, updatedEntity)
		})

		t.Run("Should update entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach()
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, _ := suite.sut.Save(entityToCreate)

			entityToUpdate := model_builders.NewSlackChannelBuilder().WithId(createdEntity.Id).Build()
			suite.sut.Save(entityToUpdate)

			suite.dbClient.GetDB().First(entityToCreate, "id = ?", createdEntity.Id)

			assert.Equal(t, entityToCreate.Id, entityToUpdate.Id)
			assert.Equal(t, entityToCreate.Name, entityToUpdate.Name)
			assert.Equal(t, entityToCreate.Agents, entityToUpdate.Agents)
			assert.Equal(t, entityToCreate.Members, entityToUpdate.Members)
			assert.Equal(t, entityToCreate.Description, entityToUpdate.Description)
		})
	})
}
