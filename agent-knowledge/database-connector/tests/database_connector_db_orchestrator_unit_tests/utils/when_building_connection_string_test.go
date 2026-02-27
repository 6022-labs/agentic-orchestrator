package utils_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuildingConnectionString(t *testing.T) {
	t.Parallel()

	t.Run("For mysql driver", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return mysql connection string", func(t *testing.T) {
			connectionString := utils.BuildMySqlConnectionString(
				"localhost",
				"database",
				"username",
				"password",
				3306,
			)

			assert.Equal(t, "username:password@tcp(localhost:3306)/database?parseTime=true", connectionString)
		})
	})

	t.Run("For sqlserver driver", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return sqlserver connection string", func(t *testing.T) {
			connectionString := utils.BuildSqlServerConnectionString(
				"localhost",
				"database",
				"username",
				"password",
				1443,
			)

			assert.Equal(t, "sqlserver://username:password@localhost:1443?database=database&encrypt=disable", connectionString)
		})
	})

	t.Run("For postgresql driver", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return postgresql connection string", func(t *testing.T) {
			connectionString := utils.BuildPostgresqlConnectionString(
				"localhost",
				"database",
				"username",
				"password",
				5432,
			)

			assert.Equal(t, "host=localhost port=5432 dbname=database user=username password=password", connectionString)
		})
	})
}
