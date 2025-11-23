package repositories_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WhenGettingNextAgentToProcessTestingSuite struct {
	dbName string

	dbClient *slack_bridge_db.DatabaseClient

	processingSettings     *settings.ProcessingSettings
	slackThreadRepository  *repositories.SlackThreadRepository
	slackChannelRepository *repositories.SlackChannelRepository

	sut *repositories.SlackAgentRepository
}

func WhenGettingNextAgentToProcessBeforeEach() *WhenGettingNextAgentToProcessTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := slack_bridge_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&db_settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	processingSettings := &settings.ProcessingSettings{
		AgentProcessingInterval: 20,
		FacilitatorTickInterval: 30,
	}

	sut := repositories.NewSlackAgentRepository(processingSettings, dbClient)
	slackThreadRepository := repositories.NewSlackThreadRepository(processingSettings, dbClient)
	slackChannelRepository := repositories.NewSlackChannelRepository(dbClient)

	return &WhenGettingNextAgentToProcessTestingSuite{
		sut: sut,

		dbName:   dbName,
		dbClient: dbClient,

		processingSettings:     processingSettings,
		slackThreadRepository:  slackThreadRepository,
		slackChannelRepository: slackChannelRepository,
	}
}

func WhenGettingNextAgentToProcessAfterEach(suite *WhenGettingNextAgentToProcessTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingNextAgentToProcess(t *testing.T) {
	t.Parallel()

	t.Run("Given no slack agent in database", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			agent, err := suite.sut.GetNextAgentToProcess()

			assert.Equal(t, gorm.ErrRecordNotFound, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agents are already processing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			agent := model_builders.NewSlackAgentBuilder().
				WithIsProcessing(true).
				Build()

			_, err := suite.sut.Save(agent)
			assert.NoError(t, err)

			nextAgent, err := suite.sut.GetNextAgentToProcess()

			assert.Equal(t, gorm.ErrRecordNotFound, err)
			assert.Nil(t, nextAgent)
		})
	})

	t.Run("Given slack agent exists and is not processing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return agent with the oldest processing time", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNextAgentToProcessBeforeEach()
			defer WhenGettingNextAgentToProcessAfterEach(suite)

			agent1 := model_builders.NewSlackAgentBuilder().
				WithCollectionTokenId(1).
				WithIsProcessing(false).
				WithLastProcessedAt(time.Now().UTC().Add(-time.Second * 40)).
				Build()

			agent2 := model_builders.NewSlackAgentBuilder().
				WithCollectionTokenId(2).
				WithIsProcessing(false).
				WithLastProcessedAt(time.Now().UTC().Add(-time.Second * 30)).
				Build()

			_, err := suite.sut.Save(agent1)
			assert.NoError(t, err)

			_, err = suite.sut.Save(agent2)
			assert.NoError(t, err)

			nextAgent, err := suite.sut.GetNextAgentToProcess()
			assert.NoError(t, err)

			assert.Equal(t, agent1.CollectionTokenId, nextAgent.CollectionTokenId)
		})
	})
}
