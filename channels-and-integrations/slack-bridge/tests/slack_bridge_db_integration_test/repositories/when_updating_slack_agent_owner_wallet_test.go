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

type WhenUpdatingSlackAgentOwnerWalletTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackAgentRepository
}

func WhenUpdatingSlackAgentOwnerWalletBeforeEach() *WhenUpdatingSlackAgentOwnerWalletTestingSuite {
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

	sut := repositories.NewSlackAgentRepository(nil, dbClient)

	return &WhenUpdatingSlackAgentOwnerWalletTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,

		sut: sut,
	}
}

func WhenUpdatingSlackAgentOwnerWalletAfterEach(suite *WhenUpdatingSlackAgentOwnerWalletTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenUpdatingSlackAgentOwnerWallet(t *testing.T) {
	t.Parallel()

	t.Run("Given slack agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingSlackAgentOwnerWalletBeforeEach()
			defer WhenUpdatingSlackAgentOwnerWalletAfterEach(suite)

			err := suite.sut.UpdateOwnerWallet("0x123", 1, "0x456")

			assert.NotNil(t, err)
		})
	})

	t.Run("Given slack agent exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should update owner wallet", func(t *testing.T) {
			t.Parallel()

			suite := WhenUpdatingSlackAgentOwnerWalletBeforeEach()
			defer WhenUpdatingSlackAgentOwnerWalletAfterEach(suite)

			slackAgent := model_builders.NewSlackAgentBuilder().Build()

			_, err := suite.sut.Save(slackAgent)
			if err != nil {
				t.Fatal(err)
			}

			err = suite.sut.UpdateOwnerWallet(slackAgent.CollectionAddress, slackAgent.CollectionTokenId, "0x456")
			if err != nil {
				t.Fatal(err)
			}

			updatedAgent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(slackAgent.CollectionAddress, slackAgent.CollectionTokenId)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "0x456", updatedAgent.OwnerWallet)
		})
	})
}
