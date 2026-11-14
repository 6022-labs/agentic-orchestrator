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

type WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite struct {
	dbName   string
	dbClient *slack_bridge_db.DatabaseClient

	sut *repositories.SlackAgentRepository
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(
	t *testing.T,
) *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite {
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
		&db_settings.DatabaseSettings{Url: connectionString},
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

	sut := repositories.NewSlackAgentRepository(dbClient)

	return &WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      sut,
	}
}

func WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(
	suite *WhenGettingAgentByCollectionAddressAndCollectionTokenIdTestingSuite,
) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentByCollectionAddressAndCollectionTokenId(t *testing.T) {
	t.Parallel()

	t.Run("Given agent does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(t)
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(12345)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(
				context.Background(),
				collectionAddress,
				collectionTokenId,
			)

			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agent exists for address but not tokenId", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(t)
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(111)

			agentModel := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				Build()
			suite.sut.Save(context.Background(), agentModel)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(
				context.Background(),
				collectionAddress,
				collectionTokenId+1,
			)
			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agent exists for tokenId but not address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(t)
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			collectionAddress := common_tests.RandomEthAddress().Hex()
			collectionTokenId := uint64(222)

			agentModel := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(collectionAddress).
				WithCollectionTokenId(collectionTokenId).
				Build()
			suite.sut.Save(context.Background(), agentModel)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(context.Background(),
				"non-existing-address",
				collectionTokenId,
			)
			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agent exists for address and tokenId", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return agent", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingAgentByCollectionAddressAndCollectionTokenIdBeforeEach(t)
			defer WhenGettingAgentByCollectionAddressAndCollectionTokenIdAfterEach(suite)

			address := common_tests.RandomEthAddress().Hex()
			tokenId := uint64(444)

			agentModel := model_builders.NewSlackAgentBuilder().
				WithCollectionAddress(address).
				WithCollectionTokenId(tokenId).
				Build()
			suite.sut.Save(context.Background(), agentModel)

			agent, err := suite.sut.GetByCollectionAddressAndCollectionTokenId(context.Background(), address, tokenId)
			assert.Nil(t, err)
			assert.NotNil(t, agent)
			assert.Equal(t, address, agent.CollectionAddress)
			assert.Equal(t, tokenId, agent.CollectionTokenId)
		})
	})
}
