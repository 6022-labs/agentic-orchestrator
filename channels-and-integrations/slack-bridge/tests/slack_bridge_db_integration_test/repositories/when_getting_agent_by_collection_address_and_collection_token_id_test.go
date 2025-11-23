package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/repositories"
	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackAgentRepository
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach() *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite {
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

	return &WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      sut,
	}
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentByCollectionAddressAndCollectionTokenId(t *testing.T) {
	t.Parallel()

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach()
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(12345)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(collectionAddress, collectionTokenId)

			assert.NotNil(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agent exists for address but not tokenId", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach()
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(111)

			agentModel := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				Build()
			suite.sut.Save(agentModel)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(collectionAddress, collectionTokenId+1)
			assert.NotNil(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agent exists for tokenId but not address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach()
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(222)

			agentModel := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				Build()
			suite.sut.Save(agentModel)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId("non-existing-address", collectionTokenId)
			assert.NotNil(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agent exists for address and tokenId", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach()
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			address := common_tests.RandomEthAddress().Hex()
			tokenId := uint64(444)

			agentModel := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(address).
				WithCollectionTokenId(tokenId).
				Build()
			suite.sut.Save(agentModel)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(address, tokenId)
			assert.Nil(t, err)
			assert.NotNil(t, agent)
			assert.Equal(t, address, agent.CollectionAddress)
			assert.Equal(t, tokenId, agent.CollectionTokenId)
		})
	})
}
