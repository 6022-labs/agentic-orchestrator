package repositories_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type WhenGettingAllAgentAddressesByAgentIdsTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentAddressRepository
}

func WhenGettingAllAgentAddressesByAgentIdsBeforeEach(
	t *testing.T,
) *WhenGettingAllAgentAddressesByAgentIdsTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewAgentAddressRepository(dbClient)

	return &WhenGettingAllAgentAddressesByAgentIdsTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAllAgentAddressesByAgentIdsAfterEach(suite *WhenGettingAllAgentAddressesByAgentIdsTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenGettingAllAgentAddressesByAgentIds(t *testing.T) {
	t.Parallel()

	t.Run("Given several agents have rows", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the rows of the requested agents only", func(t *testing.T) {
			suite := WhenGettingAllAgentAddressesByAgentIdsBeforeEach(t)
			defer WhenGettingAllAgentAddressesByAgentIdsAfterEach(suite)

			firstAgentId := uuid.New()
			secondAgentId := uuid.New()
			excludedAgentId := uuid.New()
			assert.Nil(t, suite.sut.ReplaceAllForAgent(context.Background(), firstAgentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
			}))
			assert.Nil(
				t,
				suite.sut.ReplaceAllForAgent(context.Background(), secondAgentId, map[enums.AddressType]string{
					enums.AddressTypeBtc: "bc1qexampleonly",
				}),
			)
			assert.Nil(
				t,
				suite.sut.ReplaceAllForAgent(context.Background(), excludedAgentId, map[enums.AddressType]string{
					enums.AddressTypeEvm: "0xead1a8eec112a20a8fc08beb5202aaa6538232df",
				}),
			)

			rows, err := suite.sut.GetAllByAgentIds(context.Background(), []uuid.UUID{firstAgentId, secondAgentId})

			assert.Nil(t, err)
			assert.Len(t, rows, 2)
			for _, row := range rows {
				assert.NotEqual(t, excludedAgentId, row.AgentId)
			}
		})
	})
}
