package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenGettingSlackUserByIdTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackUserRepository
}

func WhenGettingSlackUserByIdBeforeEach(t *testing.T) *WhenGettingSlackUserByIdTestingSuite {
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
		&db_settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
		dbQueryMetrics,
	)
	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewSlackUserRepository(dbClient)

	return &WhenGettingSlackUserByIdTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenGettingSlackUserByIdAfterEach(suite *WhenGettingSlackUserByIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingSlackUserById(t *testing.T) {
	t.Parallel()

	t.Run("Given slack user does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackUserByIdBeforeEach(t)
			defer WhenGettingSlackUserByIdAfterEach(suite)

			id := "non-existing-slack-user-id"

			slackUser, err := suite.sut.GetById(context.Background(), id)

			assert.NoError(t, err)
			assert.Nil(t, slackUser)
		})
	})

	t.Run("Given slack user exists for this id", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return slack user", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackUserByIdBeforeEach(t)
			defer WhenGettingSlackUserByIdAfterEach(suite)

			slackUser := model_builders.NewSlackUserBuilder().Build()

			_, err := suite.sut.Save(context.Background(), slackUser)
			if err != nil {
				t.Fatal(err)
			}

			userWallet, err := suite.sut.GetById(context.Background(), slackUser.Id)

			assert.Nil(t, err)
			assert.Equal(t, slackUser, userWallet)
		})
	})
}
