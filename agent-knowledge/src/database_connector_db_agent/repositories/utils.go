package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
)

func quoteColumnName(driver enums.DatabaseConnectionDriver, column string) string {
	switch driver {
	case enums.DatabaseConnectionDriverMariaDB, enums.DatabaseConnectionDriverMySQL:
		return "`" + column + "`"
	case enums.DatabaseConnectionDriverPostgres:
		return `"` + column + `"`
	default:
		return column
	}
}

func ilikeContainsPattern(driver enums.DatabaseConnectionDriver, column string) string {
	switch driver {
	case enums.DatabaseConnectionDriverMariaDB, enums.DatabaseConnectionDriverMySQL:
		return "LOWER(" + column + ") LIKE LOWER(CONCAT('%', ?, '%'))"
	case enums.DatabaseConnectionDriverPostgres:
		return column + " ILIKE '%' || ? || '%'"
	default:
		return column
	}
}
