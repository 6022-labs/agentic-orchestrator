package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenSavingEntityTestingSuite struct {
	dbName string

	dbClient *slack_bridge_db.DatabaseClient
	sut      *repositories.SlackChannelRepository
}

func WhenSavingEntityBeforeEach(t *testing.T) *WhenSavingEntityTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
		dbQueryMetrics,
	)
	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewSlackChannelRepository(dbClient)

	return &WhenSavingEntityTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenSavingEntityAfterEach(suite *WhenSavingEntityTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingEntity(t *testing.T) {
	t.Parallel()

	t.Run("Given entity does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return created entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, err := suite.sut.Save(context.Background(), entityToCreate)

			assert.Nil(t, err)
			assert.Equal(t, entityToCreate, createdEntity)
		})

		t.Run("Should store entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			suite.dbClient.GetDB().First(entityToCreate, "id = ?", createdEntity.Id)

			assert.Equal(t, entityToCreate, createdEntity)
			assert.NotEqual(t, uuid.Nil, createdEntity.Id)
		})
	})

	t.Run("Given entity exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return updated entity", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			entityToUpdate := model_builders.NewSlackChannelBuilder().WithId(createdEntity.Id).Build()
			updatedEntity, err := suite.sut.Save(context.Background(), entityToUpdate)

			assert.Nil(t, err)
			assert.Equal(t, entityToUpdate, updatedEntity)
		})

		t.Run("Should update entity in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingEntityBeforeEach(t)
			defer WhenSavingEntityAfterEach(suite)

			entityToCreate := model_builders.NewSlackChannelBuilder().Build()
			createdEntity, _ := suite.sut.Save(context.Background(), entityToCreate)

			entityToUpdate := model_builders.NewSlackChannelBuilder().WithId(createdEntity.Id).Build()
			suite.sut.Save(context.Background(), entityToUpdate)

			suite.dbClient.GetDB().First(entityToCreate, "id = ?", createdEntity.Id)

			assert.Equal(t, entityToCreate.Id, entityToUpdate.Id)
			assert.Equal(t, entityToCreate.Name, entityToUpdate.Name)
			assert.Equal(t, entityToCreate.Agents, entityToUpdate.Agents)
			assert.Equal(t, entityToCreate.Description, entityToUpdate.Description)
		})
	})
}
