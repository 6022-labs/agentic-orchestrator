package factories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/factories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCreatingConversationBrokerDbClientTestSuite struct {
	dbName string

	migrationsSettings *settings.MigrationsSettings
	databaseSettings   *settings.DatabaseSettings
	sut                *factories.DatabaseClientFactory
}

func WhenCreatingConversationBrokerDbClientBeforeEach() *WhenCreatingConversationBrokerDbClientTestSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	databaseSettings := &settings.DatabaseSettings{
		Url: connectionString,
	}
	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: "",
	}

	sut := factories.NewDatabaseClientFactory(
		zap.NewNop(),
		migrationsSettings,
		databaseSettings,
	)

	return &WhenCreatingConversationBrokerDbClientTestSuite{
		sut:    sut,
		dbName: dbName,

		migrationsSettings: migrationsSettings,
		databaseSettings:   databaseSettings,
	}
}

func WhenCreatingConversationBrokerDbClientAfterEach(suite *WhenCreatingConversationBrokerDbClientTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCreatingConversationBrokerDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Should return a new conversation broker db client", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingConversationBrokerDbClientBeforeEach()
		defer WhenCreatingConversationBrokerDbClientAfterEach(suite)

		dbClient, err := suite.sut.CreateDatabaseClient()

		assert.Nil(t, err)
		assert.NotNil(t, dbClient)
	})
}
