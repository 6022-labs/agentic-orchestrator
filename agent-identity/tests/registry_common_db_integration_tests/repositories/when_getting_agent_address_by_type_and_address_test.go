package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenGettingAgentAddressByTypeAndAddressTestSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentAddressRepository
}

func WhenGettingAgentAddressByTypeAndAddressBeforeEach() *WhenGettingAgentAddressByTypeAndAddressTestSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, _ := registry_common_db.NewDatabaseClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}

	sut := repositories.NewAgentAddressRepository(dbClient)

	return &WhenGettingAgentAddressByTypeAndAddressTestSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAgentAddressByTypeAndAddressAfterEach(suite *WhenGettingAgentAddressByTypeAndAddressTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAgentAddressByTypeAndAddress(t *testing.T) {
	t.Parallel()

	t.Run("Given the row exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the row with its agent id", func(t *testing.T) {
			suite := WhenGettingAgentAddressByTypeAndAddressBeforeEach()
			defer WhenGettingAgentAddressByTypeAndAddressAfterEach(suite)

			agentId := uuid.New()
			err := suite.sut.ReplaceAllForAgent(agentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
				enums.AddressTypeBtc: "bc1qexampleonly",
			})
			assert.Nil(t, err)

			row, err := suite.sut.GetByTypeAndAddress(enums.AddressTypeBtc, "bc1qexampleonly")

			assert.Nil(t, err)
			assert.Equal(t, agentId, row.AgentId)
			assert.Equal(t, enums.AddressTypeBtc, row.Type)
		})
	})

	t.Run("Given the row does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			suite := WhenGettingAgentAddressByTypeAndAddressBeforeEach()
			defer WhenGettingAgentAddressByTypeAndAddressAfterEach(suite)

			row, err := suite.sut.GetByTypeAndAddress(enums.AddressTypeEvm, "0xunknown")

			assert.Nil(t, row)
			assert.NotNil(t, err)
		})
	})
}
