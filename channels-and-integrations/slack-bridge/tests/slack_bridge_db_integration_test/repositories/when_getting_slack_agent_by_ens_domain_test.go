package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingSlackAgentByEnsDomainTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackAgentRepository
}

func WhenGettingSlackAgentByEnsDomainBeforeEach() *WhenGettingSlackAgentByEnsDomainTestingSuite {
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

	sut := repositories.NewSlackAgentRepository(dbClient)

	return &WhenGettingSlackAgentByEnsDomainTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,

		sut: sut,
	}
}

func WhenGettingSlackAgentByEnsDomainAfterEach(suite *WhenGettingSlackAgentByEnsDomainTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingSlackAgentByEnsDomain(t *testing.T) {
	t.Parallel()

	t.Run("Given slack agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackAgentByEnsDomainBeforeEach()
			defer WhenGettingSlackAgentByEnsDomainAfterEach(suite)

			ensDomain := "nonexistent.eth"

			slackAgent, err := suite.sut.GetByEnsDomain(ensDomain)

			assert.NotNil(t, err)
			assert.Nil(t, slackAgent)
		})
	})

	t.Run("Given slack agent exists for this ens domain", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return slack agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingSlackAgentByEnsDomainBeforeEach()
			defer WhenGettingSlackAgentByEnsDomainAfterEach(suite)

			slackAgent := model_builders.NewSlackAgentBuilder().
				WithEnsDomain("existing.eth").
				Build()

			_, err := suite.sut.Save(slackAgent)
			if err != nil {
				t.Fatal(err)
			}

			result, err := suite.sut.GetByEnsDomain(slackAgent.EnsDomain)

			assert.Nil(t, err)
			assert.Equal(t, slackAgent, result)
		})
	})
}
