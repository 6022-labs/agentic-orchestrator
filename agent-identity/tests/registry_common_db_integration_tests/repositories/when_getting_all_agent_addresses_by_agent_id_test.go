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

type WhenGettingAllAgentAddressesByAgentIdTestingSuite struct {
	dbName string

	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.AgentAddressRepository
}

func WhenGettingAllAgentAddressesByAgentIdBeforeEach(t *testing.T) *WhenGettingAllAgentAddressesByAgentIdTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	sut := repositories.NewAgentAddressRepository(dbClient)

	return &WhenGettingAllAgentAddressesByAgentIdTestingSuite{
		sut:    sut,
		dbName: dbName,

		dbClient: dbClient,
	}
}

func WhenGettingAllAgentAddressesByAgentIdAfterEach(suite *WhenGettingAllAgentAddressesByAgentIdTestingSuite) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenGettingAllAgentAddressesByAgentId(t *testing.T) {
	t.Parallel()

	t.Run("Given the agent has rows", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return only that agent's rows", func(t *testing.T) {
			suite := WhenGettingAllAgentAddressesByAgentIdBeforeEach(t)
			defer WhenGettingAllAgentAddressesByAgentIdAfterEach(suite)

			agentId := uuid.New()
			otherAgentId := uuid.New()
			assert.Nil(t, suite.sut.ReplaceAllForAgent(context.Background(), agentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0x2f0ddd93040c19863527b0553feceb696862ea7d",
				enums.AddressTypeBtc: "bc1qexampleonly",
			}))
			assert.Nil(t, suite.sut.ReplaceAllForAgent(context.Background(), otherAgentId, map[enums.AddressType]string{
				enums.AddressTypeEvm: "0xead1a8eec112a20a8fc08beb5202aaa6538232df",
			}))

			rows, err := suite.sut.GetAllByAgentId(context.Background(), agentId)

			assert.Nil(t, err)
			assert.Len(t, rows, 2)
			for _, row := range rows {
				assert.Equal(t, agentId, row.AgentId)
			}
		})
	})

	t.Run("Given the agent has no rows", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an empty slice", func(t *testing.T) {
			suite := WhenGettingAllAgentAddressesByAgentIdBeforeEach(t)
			defer WhenGettingAllAgentAddressesByAgentIdAfterEach(suite)

			rows, err := suite.sut.GetAllByAgentId(context.Background(), uuid.New())

			assert.Nil(t, err)
			assert.Empty(t, rows)
		})
	})
}
