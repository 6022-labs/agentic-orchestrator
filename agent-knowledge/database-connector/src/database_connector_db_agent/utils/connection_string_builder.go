package utils

import "fmt"

func BuildMySqlConnectionString(host, database, username, password string, port uint16) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, host, port, database)
}

func BuildSqlServerConnectionString(host, database, username, password string, port uint16) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable", username, password, host, port, database)
}

func BuildPostgresqlConnectionString(host, database, username, password string, port uint16) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, database, username, password)
}
