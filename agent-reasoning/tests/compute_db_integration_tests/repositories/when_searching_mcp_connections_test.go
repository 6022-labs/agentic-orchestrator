package repositories_test

import (
	"testing"

	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/repositories"
	orchestrator_agent_compute_db "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db"
	db_repositories "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_db/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenSearchingMCPConnectionsTestingSuite struct {
	sut *db_repositories.MCPConnectionRepository

	dbName   string
	dbClient db_clients.DatabaseClientInterface
}

func WhenSearchingMCPConnectionsBeforeEach() *WhenSearchingMCPConnectionsTestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	dbClient, _ := orchestrator_agent_compute_db.NewOrchestratorDbClient(
		connectionString,
		zap.NewNop(),
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	dbClient.Migrate()

	sut := db_repositories.NewMCPConnectionRepository(dbClient)

	return &WhenSearchingMCPConnectionsTestingSuite{
		sut:      sut,
		dbName:   dbName,
		dbClient: dbClient,
	}
}

func WhenSearchingMCPConnectionsAfterEach(suite *WhenSearchingMCPConnectionsTestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhenSearchingMCPConnections(t *testing.T) {
	t.Parallel()

	t.Run("Given AgentId parameter", func(t *testing.T) {
		t.Parallel()

		t.Run("And multiple connections for different agents exist", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return only connections for the specified agent", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingMCPConnectionsBeforeEach()
				defer WhenSearchingMCPConnectionsAfterEach(suite)

				agentId := uuid.New()
				otherAgentId := uuid.New()
				conn1 := model_builders.NewMCPConnectionBuilder().WithAgentId(agentId).Build()
				conn2 := model_builders.NewMCPConnectionBuilder().WithAgentId(agentId).Build()
				conn3 := model_builders.NewMCPConnectionBuilder().WithAgentId(otherAgentId).Build()
				suite.sut.Save(conn1)
				suite.sut.Save(conn2)
				suite.sut.Save(conn3)

				params := repositories.MCPConnectionSearchParams{AgentId: &agentId}
				results, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, results, 2)
				for _, c := range results {
					assert.Equal(t, agentId, c.AgentID)
				}
			})
		})
	})

	t.Run("Given no AgentId parameter", func(t *testing.T) {
		t.Parallel()

		t.Run("And multiple connections for different agents exist", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return all connections", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingMCPConnectionsBeforeEach()
				defer WhenSearchingMCPConnectionsAfterEach(suite)

				agentId := uuid.New()
				otherAgentId := uuid.New()
				conn1 := model_builders.NewMCPConnectionBuilder().WithAgentId(agentId).Build()
				conn2 := model_builders.NewMCPConnectionBuilder().WithAgentId(otherAgentId).Build()
				suite.sut.Save(conn1)
				suite.sut.Save(conn2)

				params := repositories.MCPConnectionSearchParams{}
				results, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, results, 2)
			})
		})
	})

	t.Run("Given Search parameter", func(t *testing.T) {
		t.Parallel()

		t.Run("And connections with matching names or urls exist", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return only connections whose name or url matches (ILIKE)", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingMCPConnectionsBeforeEach()
				defer WhenSearchingMCPConnectionsAfterEach(suite)

				agentId := uuid.New()
				otherAgentId := uuid.New()

				// Matches by name
				connByName := model_builders.NewMCPConnectionBuilder().
					WithAgentId(agentId).
					WithName("alpha-connector").
					WithUrl("https://service-one.example.com").
					Build()

				// Matches by URL
				connByURL := model_builders.NewMCPConnectionBuilder().
					WithAgentId(agentId).
					WithName("irrelevant").
					WithUrl("https://alpha-services.example.com/api").
					Build()

				// Doesn't match
				nonMatching := model_builders.NewMCPConnectionBuilder().
					WithAgentId(otherAgentId).
					WithName("beta").
					WithUrl("https://beta.example.com").
					Build()

				suite.sut.Save(connByName)
				suite.sut.Save(connByURL)
				suite.sut.Save(nonMatching)

				term := "alpha"
				params := repositories.MCPConnectionSearchParams{Search: &term}
				results, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.Len(t, results, 2)
				for _, c := range results {
					containsInName := strings.Contains(strings.ToLower(c.Name), strings.ToLower(term))
					containsInURL := strings.Contains(strings.ToLower(c.Url), strings.ToLower(term))
					assert.True(t, containsInName || containsInURL)
				}
			})
		})
	})

	t.Run("Given Limit parameter", func(t *testing.T) {
		t.Parallel()

		t.Run("And more connections exist than the limit", func(t *testing.T) {
			t.Parallel()

			t.Run("Should return at most 'limit' connections", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingMCPConnectionsBeforeEach()
				defer WhenSearchingMCPConnectionsAfterEach(suite)

				c1 := model_builders.NewMCPConnectionBuilder().Build()
				c2 := model_builders.NewMCPConnectionBuilder().Build()
				c3 := model_builders.NewMCPConnectionBuilder().Build()
				suite.sut.Save(c1)
				suite.sut.Save(c2)
				suite.sut.Save(c3)

				var limit uint = 2
				params := repositories.MCPConnectionSearchParams{Limit: &limit}
				results, err := suite.sut.Search(params)

				assert.Nil(t, err)
				assert.LessOrEqual(t, len(results), int(limit))

				// With 3 rows present, expect exactly 2 returned.
				assert.Equal(t, int(limit), len(results))
			})
		})
	})

	t.Run("Given Offset parameter", func(t *testing.T) {
		t.Parallel()

		t.Run("And multiple connections exist", func(t *testing.T) {
			t.Parallel()

			t.Run("Should skip the first N results", func(t *testing.T) {
				t.Parallel()
				suite := WhenSearchingMCPConnectionsBeforeEach()
				defer WhenSearchingMCPConnectionsAfterEach(suite)

				// Create 3 connections
				c1 := model_builders.NewMCPConnectionBuilder().Build()
				c2 := model_builders.NewMCPConnectionBuilder().Build()
				c3 := model_builders.NewMCPConnectionBuilder().Build()
				suite.sut.Save(c1)
				suite.sut.Save(c2)
				suite.sut.Save(c3)

				var offset uint = 2
				params := repositories.MCPConnectionSearchParams{Offset: &offset}
				results, err := suite.sut.Search(params)

				assert.Nil(t, err)

				// With 3 rows and offset 2, we should get at most 1 back.
				assert.LessOrEqual(t, len(results), 1)
			})
		})
	})

	t.Run("Given Offset and Limit parameters together", func(t *testing.T) {
		t.Parallel()

		t.Run("And multiple connections exist", func(t *testing.T) {
			t.Parallel()
			suite := WhenSearchingMCPConnectionsBeforeEach()
			defer WhenSearchingMCPConnectionsAfterEach(suite)

			for i := 0; i < 5; i++ {
				suite.sut.Save(model_builders.NewMCPConnectionBuilder().Build())
			}

			var offset uint = 1
			var limit uint = 2
			params := repositories.MCPConnectionSearchParams{Offset: &offset, Limit: &limit}
			results, err := suite.sut.Search(params)

			assert.Nil(t, err)
			assert.LessOrEqual(t, len(results), int(limit))
			assert.Equal(t, int(limit), len(results))
		})
	})

	t.Run("Given OrderBy parameter", func(t *testing.T) {
		t.Parallel()

		t.Run("And multiple connections with varying values exist", func(t *testing.T) {
			t.Parallel()
			suite := WhenSearchingMCPConnectionsBeforeEach()
			defer WhenSearchingMCPConnectionsAfterEach(suite)

			// Insert in non-sorted order (by name)
			c1 := model_builders.NewMCPConnectionBuilder().WithName("delta").Build()
			c2 := model_builders.NewMCPConnectionBuilder().WithName("alpha").Build()
			c3 := model_builders.NewMCPConnectionBuilder().WithName("charlie").Build()
			suite.sut.Save(c1)
			suite.sut.Save(c2)
			suite.sut.Save(c3)

			orderBy := "name"
			params := repositories.MCPConnectionSearchParams{OrderBy: &orderBy}
			results, err := suite.sut.Search(params)

			assert.Nil(t, err)
			assert.Len(t, results, 3)

			expectedAsc := []string{"alpha", "charlie", "delta"}
			for i, r := range results {
				assert.Equal(t, expectedAsc[i], r.Name)
			}
		})
	})

	t.Run("Given OrderBy and OrderByDirection parameters", func(t *testing.T) {
		t.Parallel()

		t.Run("And multiple connections with varying values exist", func(t *testing.T) {
			t.Parallel()
			suite := WhenSearchingMCPConnectionsBeforeEach()
			defer WhenSearchingMCPConnectionsAfterEach(suite)

			c1 := model_builders.NewMCPConnectionBuilder().WithName("delta").Build()
			c2 := model_builders.NewMCPConnectionBuilder().WithName("alpha").Build()
			c3 := model_builders.NewMCPConnectionBuilder().WithName("charlie").Build()
			suite.sut.Save(c1)
			suite.sut.Save(c2)
			suite.sut.Save(c3)

			orderBy := "name"
			orderByDirection := "DESC"
			params := repositories.MCPConnectionSearchParams{OrderBy: &orderBy, OrderByDirection: &orderByDirection}
			results, err := suite.sut.Search(params)

			assert.Nil(t, err)
			assert.Len(t, results, 3)

			expectedDesc := []string{"delta", "charlie", "alpha"}
			for i, r := range results {
				assert.Equal(t, expectedDesc[i], r.Name)
			}
		})
	})
}
