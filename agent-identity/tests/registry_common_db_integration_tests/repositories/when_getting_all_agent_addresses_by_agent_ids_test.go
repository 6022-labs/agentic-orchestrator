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

type WhenGettingAllAgentAddressesByAgentIdsTestSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentAddressRepository
}

func WhenGettingAllAgentAddressesByAgentIdsBeforeEach() *WhenGettingAllAgentAddressesByAgentIdsTestSuite {
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

	return &WhenGettingAllAgentAddressesByAgentIdsTestSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAllAgentAddressesByAgentIdsAfterEach(suite *WhenGettingAllAgentAddressesByAgentIdsTestSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenGettingAllAgentAddressesByAgentIds(t *testing.T) {
	t.Parallel()

	t.Run("Given several agents have rows", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the rows of the requested agents only", func(t *testing.T) {
			suite := WhenGettingAllAgentAddressesByAgentIdsBeforeEach()
			defer WhenGettingAllAgentAddressesByAgentIdsAfterEach(suite)

			firstAgentId := uuid.New()
			secondAgentId := uuid.New()
			excludedAgentId := uuid.New()
			assert.Nil(t, suite.sut.ReplaceAllForAgent(firstAgentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
			}))
			assert.Nil(t, suite.sut.ReplaceAllForAgent(secondAgentId, map[enums.AddressType]string{
				enums.AddressTypeBtc: "bc1qexampleonly",
			}))
			assert.Nil(t, suite.sut.ReplaceAllForAgent(excludedAgentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0xead1a8eec112a20a8fc08beb5202aaa6538232df",
			}))

			rows, err := suite.sut.GetAllByAgentIds([]uuid.UUID{firstAgentId, secondAgentId})

			assert.Nil(t, err)
			assert.Len(t, rows, 2)
			for _, row := range rows {
				assert.NotEqual(t, excludedAgentId, row.AgentId)
			}
		})
	})
}
