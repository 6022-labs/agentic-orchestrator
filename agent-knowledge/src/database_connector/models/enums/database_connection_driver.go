package enums

import "strings"

type DatabaseConnectionDriver string
type DatabaseConnectionDriverList []DatabaseConnectionDriver

const (
	DatabaseConnectionDriverMySQL    DatabaseConnectionDriver = "mysql"
	DatabaseConnectionDriverMariaDB  DatabaseConnectionDriver = "mariadb"
	DatabaseConnectionDriverPostgres DatabaseConnectionDriver = "postgres"
)

func GetAllDatabaseConnectionDrivers() DatabaseConnectionDriverList {
	return DatabaseConnectionDriverList{
		DatabaseConnectionDriverMySQL,
		DatabaseConnectionDriverMariaDB,
		DatabaseConnectionDriverPostgres,
	}
}

func (k DatabaseConnectionDriver) String() string {
	return string(k)
}

func (k DatabaseConnectionDriver) ToGooseDialect() string {
	switch k {
	case DatabaseConnectionDriverMySQL:
		return "mysql"
	case DatabaseConnectionDriverMariaDB:
		return "mysql"
	case DatabaseConnectionDriverPostgres:
		return "postgres"
	default:
		return ""
	}
}

func (k DatabaseConnectionDriverList) String() string {
	var drivers []string
	for _, v := range k {
		drivers = append(drivers, string(v))
	}
	return strings.Join(drivers, ", ")
}

func (k DatabaseConnectionDriverList) Contains(driver DatabaseConnectionDriver) bool {
	for _, v := range k {
		if v == driver {
			return true
		}
	}
	return false
}
