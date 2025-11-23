package repositories_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenCreatingMultipleAgentContextMessagesTestCase struct {
	name            string
	clientBuilderFn func() container_builders.DatabaseContainerBuilderInterface
	driver          enums.DatabaseConnectionDriver
}

var (
	WhenCreatingMultipleAgentContextMessagesTestCases = []*WhenCreatingMultipleAgentContextMessagesTestCase{
		{
			name:   "PostgreSQL",
			driver: enums.DatabaseConnectionDriverPostgres,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetPostgresContainerBuilder()
			},
		},
		{
			name:   "MariaDB",
			driver: enums.DatabaseConnectionDriverMariaDB,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMariaDBContainerBuilder()
			},
		},
		{
			name:   "MySQL",
			driver: enums.DatabaseConnectionDriverMySQL,
			clientBuilderFn: func() container_builders.DatabaseContainerBuilderInterface {
				return container_builders.GetMySQLContainerBuilder()
			},
		},
	}
)

type WhenCreatingMultipleAgentContextMessagesTestingSuite struct {
	sut *repositories.AgentContextMessageRepository

	agentContextRepository *repositories.AgentContextRepository

	dbName                   string
	dbClient                 db_clients.AgentDatabaseClientInterface
	databaseContainerBuilder container_builders.DatabaseContainerBuilderInterface
}

func WhenCreatingMultipleAgentContextMessagesBeforeEach(testCase *WhenCreatingMultipleAgentContextMessagesTestCase) *WhenCreatingMultipleAgentContextMessagesTestingSuite {
	databaseContainerBuilder := testCase.clientBuilderFn()
	dbName, connectionString := databaseContainerBuilder.Create()

	migrationFolderPath, err := common_tests.FindAgentMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, err := database_connector_db_agent.NewAgentDbClient(
		connectionString,
		zap.NewNop(),
		testCase.driver,
		&settings.MigrationsSettings{
			MigrationsFolderPath: *migrationFolderPath,
		},
	)
	if err != nil {
		panic(err)
	}

	err = dbClient.Migrate()
	if err != nil {
		panic(err)
	}

	sut := repositories.NewAgentContextMessageRepository(dbClient)
	agentContextRepository := repositories.NewAgentContextRepository(dbClient)

	return &WhenCreatingMultipleAgentContextMessagesTestingSuite{
		sut: sut,

		agentContextRepository: agentContextRepository,

		dbName:                   dbName,
		dbClient:                 dbClient,
		databaseContainerBuilder: databaseContainerBuilder,
	}
}

func WhenCreatingMultipleAgentContextMessagesAfterEach(suite *WhenCreatingMultipleAgentContextMessagesTestingSuite) {
	suite.databaseContainerBuilder.DropDatabase(suite.dbName)
}

func TestWhenCreatingMultipleAgentContextMessages(t *testing.T) {
	t.Parallel()

	for _, testCase := range WhenCreatingMultipleAgentContextMessagesTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Given no entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should not return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingMultipleAgentContextMessagesBeforeEach(testCase)
					defer WhenCreatingMultipleAgentContextMessagesAfterEach(suite)

					context := model_builders.NewAgentContextBuilder().
						WithKey("key1").
						Build()

					_, err := suite.agentContextRepository.Save(context)
					if err != nil {
						t.Fatal(err)
					}

					messagesToCreate := []*models.AgentContextMessage{
						model_builders.NewAgentContextMessageBuilder().
							WithAgentContextId(context.ID).
							WithContent("Message 1").
							Build(),
						model_builders.NewAgentContextMessageBuilder().
							WithAgentContextId(context.ID).
							WithContent("Message 2").
							Build(),
						model_builders.NewAgentContextMessageBuilder().
							WithAgentContextId(context.ID).
							WithContent("Message 3").
							Build(),
					}

					_, err = suite.sut.CreateMultiple(messagesToCreate)
					assert.Nil(t, err)
				})
			})

			t.Run("Given entities exist", func(t *testing.T) {
				t.Parallel()

				t.Run("Should return error", func(t *testing.T) {
					t.Parallel()

					suite := WhenCreatingMultipleAgentContextMessagesBeforeEach(testCase)
					defer WhenCreatingMultipleAgentContextMessagesAfterEach(suite)

					context := model_builders.NewAgentContextBuilder().
						WithKey("key1").
						Build()

					_, err := suite.agentContextRepository.Save(context)
					if err != nil {
						t.Fatal(err)
					}

					messagesToCreate := []*models.AgentContextMessage{
						model_builders.NewAgentContextMessageBuilder().
							WithAgentContextId(context.ID).
							WithContent("Message 1").
							Build(),
						model_builders.NewAgentContextMessageBuilder().
							WithAgentContextId(context.ID).
							WithContent("Message 2").
							Build(),
						model_builders.NewAgentContextMessageBuilder().
							WithAgentContextId(context.ID).
							WithContent("Message 3").
							Build(),
					}

					_, err = suite.sut.CreateMultiple(messagesToCreate)
					if err != nil {
						t.Fatal(err)
					}

					_, err = suite.sut.CreateMultiple(messagesToCreate)

					assert.NotNil(t, err)
				})
			})
		})
	}
}
