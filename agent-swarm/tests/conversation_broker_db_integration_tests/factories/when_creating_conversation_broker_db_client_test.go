package factories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/factories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenCreatingConversationBrokerDbClientTestingSuite struct {
	dbName string

	migrationsSettings *settings.MigrationsSettings
	databaseSettings   *settings.DatabaseSettings
	sut                *factories.DatabaseClientFactory
}

func WhenCreatingConversationBrokerDbClientBeforeEach(
	t *testing.T,
) *WhenCreatingConversationBrokerDbClientTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	databaseSettings := &settings.DatabaseSettings{
		Url: connectionString,
	}
	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: "",
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	sut := factories.NewDatabaseClientFactory(
		zap.NewNop(),
		dbQueryMetrics,
		migrationsSettings,
		databaseSettings,
	)

	return &WhenCreatingConversationBrokerDbClientTestingSuite{
		sut:    sut,
		dbName: dbName,

		migrationsSettings: migrationsSettings,
		databaseSettings:   databaseSettings,
	}
}

func WhenCreatingConversationBrokerDbClientAfterEach(suite *WhenCreatingConversationBrokerDbClientTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenCreatingConversationBrokerDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Should return a new conversation broker db client", func(t *testing.T) {
		t.Parallel()

		suite := WhenCreatingConversationBrokerDbClientBeforeEach(t)
		defer WhenCreatingConversationBrokerDbClientAfterEach(suite)

		dbClient, err := suite.sut.CreateDatabaseClient()

		assert.Nil(t, err)
		assert.NotNil(t, dbClient)
	})
}
