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

type WhenGettingAgentAddressByTypeAndAddressTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentAddressRepository
}

func WhenGettingAgentAddressByTypeAndAddressBeforeEach(
	t *testing.T,
) *WhenGettingAgentAddressByTypeAndAddressTestingSuite {
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

	return &WhenGettingAgentAddressByTypeAndAddressTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAgentAddressByTypeAndAddressAfterEach(suite *WhenGettingAgentAddressByTypeAndAddressTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentAddressByTypeAndAddress(t *testing.T) {
	t.Parallel()

	t.Run("Given the row exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the row with its agent id", func(t *testing.T) {
			suite := WhenGettingAgentAddressByTypeAndAddressBeforeEach(t)
			defer WhenGettingAgentAddressByTypeAndAddressAfterEach(suite)

			agentId := uuid.New()
			err := suite.sut.ReplaceAllForAgent(context.Background(), agentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
				enums.AddressTypeBtc: "bc1qexampleonly",
			})
			assert.Nil(t, err)

			row, err := suite.sut.GetByTypeAndAddress(context.Background(), enums.AddressTypeBtc, "bc1qexampleonly")

			assert.Nil(t, err)
			assert.Equal(t, agentId, row.AgentId)
			assert.Equal(t, enums.AddressTypeBtc, row.Type)
		})
	})

	t.Run("Given the row does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil without error", func(t *testing.T) {
			suite := WhenGettingAgentAddressByTypeAndAddressBeforeEach(t)
			defer WhenGettingAgentAddressByTypeAndAddressAfterEach(suite)

			row, err := suite.sut.GetByTypeAndAddress(context.Background(), enums.AddressTypeEvm, "0xunknown")

			assert.Nil(t, row)
			assert.Nil(t, err)
		})
	})
}
