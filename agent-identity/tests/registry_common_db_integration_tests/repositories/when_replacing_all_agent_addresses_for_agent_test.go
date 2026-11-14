package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenReplacingAllAgentAddressesForAgentTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentAddressRepository
}

func WhenReplacingAllAgentAddressesForAgentBeforeEach(
	t *testing.T,
) *WhenReplacingAllAgentAddressesForAgentTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	dbClient, _ := registry_common_db.NewDatabaseClient(
		&settings.DatabaseSettings{Url: connectionString},
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
		dbQueryMetrics,
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := repositories.NewAgentAddressRepository(dbClient)

	return &WhenReplacingAllAgentAddressesForAgentTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenReplacingAllAgentAddressesForAgentAfterEach(suite *WhenReplacingAllAgentAddressesForAgentTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenReplacingAllAgentAddressesForAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given the agent has no rows yet", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store one row per type", func(t *testing.T) {
			suite := WhenReplacingAllAgentAddressesForAgentBeforeEach(t)
			defer WhenReplacingAllAgentAddressesForAgentAfterEach(suite)

			agentId := uuid.New()
			err := suite.sut.ReplaceAllForAgent(context.Background(), agentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
				enums.AddressTypeBtc: "bc1qexampleonly",
			})

			assert.Nil(t, err)

			rows, err := suite.sut.GetAllByAgentId(context.Background(), agentId)
			assert.Nil(t, err)
			assert.Len(t, rows, 2)
		})
	})

	t.Run("Given the agent already has rows", func(t *testing.T) {
		t.Parallel()

		t.Run("Should drop rows missing from the new set", func(t *testing.T) {
			suite := WhenReplacingAllAgentAddressesForAgentBeforeEach(t)
			defer WhenReplacingAllAgentAddressesForAgentAfterEach(suite)

			agentId := uuid.New()
			err := suite.sut.ReplaceAllForAgent(context.Background(), agentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
				enums.AddressTypeBtc: "bc1qexampleonly",
			})
			assert.Nil(t, err)

			err = suite.sut.ReplaceAllForAgent(context.Background(), agentId, map[enums.AddressType]string{
				enums.AddressTypeBtc: "bc1qanothervalue",
			})
			assert.Nil(t, err)

			rows, err := suite.sut.GetAllByAgentId(context.Background(), agentId)
			assert.Nil(t, err)
			if assert.Len(t, rows, 1) {
				assert.Equal(t, enums.AddressTypeBtc, rows[0].Type)
				assert.Equal(t, "bc1qanothervalue", rows[0].Address)
			}
		})
	})
}
