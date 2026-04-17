package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
)

type WhenSavingBridgeTestingSuite struct {
	dbName   string
	dbClient db_clients.DatabaseClientInterface
	sut      *repositories.BridgeRepository
}

func WhenSavingBridgeBeforeEach() *WhenSavingBridgeTestingSuite {
	dbName, dbClient := newMigratedDatabase()
	return &WhenSavingBridgeTestingSuite{
		dbName:   dbName,
		dbClient: dbClient,
		sut:      repositories.NewBridgeRepository(dbClient),
	}
}

func WhenSavingBridgeAfterEach(suite *WhenSavingBridgeTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSavingBridge(t *testing.T) {
	t.Parallel()

	t.Run("Given bridge does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store bridge in database", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingBridgeBeforeEach()
			defer WhenSavingBridgeAfterEach(suite)

			entity := model_builders.NewBridgeBuilder().Build()
			saved, err := suite.sut.Save(entity)

			assert.NoError(t, err)
			assert.Equal(t, entity.ID, saved.ID)

			byId, err := suite.sut.GetById(entity.ID)
			assert.NoError(t, err)
			assert.Equal(t, entity.BaseUrl, byId.BaseUrl)
		})
	})

	t.Run("Given multiple bridges of the same kind", func(t *testing.T) {
		t.Parallel()

		t.Run("Should list all of them via ListByKind", func(t *testing.T) {
			t.Parallel()

			suite := WhenSavingBridgeBeforeEach()
			defer WhenSavingBridgeAfterEach(suite)

			first := model_builders.NewBridgeBuilder().Build()
			second := model_builders.NewBridgeBuilder().Build()
			_, err := suite.sut.Save(first)
			assert.NoError(t, err)
			_, err = suite.sut.Save(second)
			assert.NoError(t, err)

			list, err := suite.sut.ListByKind(enums.BridgeKindSlack)
			assert.NoError(t, err)
			assert.Len(t, list, 2)
		})
	})
}
