package container_builders

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	sqlServerContainer       testcontainers.Container
	sqlServerConnStrTemplate string
	sqlServerContainerMu     sync.Mutex
	sqlServerDatabasesMu     sync.Mutex
	sqlServerDatabases       []string
)

// Start the SQL Server container once for all test packages
func startSQLServer() string {
	sqlServerContainerMu.Lock()
	defer sqlServerContainerMu.Unlock()

	if sqlServerContainer == nil {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "mcr.microsoft.com/mssql/server:latest",
			ExposedPorts: []string{"1433/tcp"},
			Env: map[string]string{
				"ACCEPT_EULA":       "Y",
				"MSSQL_SA_PASSWORD": "@Azertyiop123",
			},
			WaitingFor: wait.ForListeningPort("1433/tcp"),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("Failed to start SQL Server container: %v", err)
		}
		sqlServerContainer = container

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "1433")

		// SQL Server DSN format
		sqlServerConnStrTemplate = fmt.Sprintf("sqlserver://sa:@Azertyiop123@%s:%s?database=%s&encrypt=disable", host, port.Port(), "%s")

		log.Println("SQL Server container started and ready")
	}

	return sqlServerConnStrTemplate
}

// Stop the container after all tests
func stopSQLServer() {
	sqlServerContainerMu.Lock()
	defer sqlServerContainerMu.Unlock()

	if sqlServerContainer != nil {
		sqlServerContainer.Terminate(context.Background())
		sqlServerContainer = nil

		log.Println("SQL Server container stopped")
	}
}

func CreateSQLServerDatabase() (string, string) {
	// Connect as SA user to create a new database
	dnsTemplate := startSQLServer()
	mainDatabaseConnectionString := fmt.Sprintf(dnsTemplate, "master") // Always connect to master for creation

	conn, err := sql.Open("sqlserver", mainDatabaseConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to SQL Server: %v", err)
	}
	defer conn.Close()

	dbName := uuid.New().String()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE [%s];", dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}

	sqlServerDatabasesMu.Lock()
	sqlServerDatabases = append(sqlServerDatabases, dbName)
	sqlServerDatabasesMu.Unlock()

	return dbName, fmt.Sprintf(dnsTemplate, dbName)
}

func DestroySQLServerDatabase(dbName string) {
	// Just remove from the list, the container will be destroyed after all tests
	for i, db := range sqlServerDatabases {
		if db == dbName {
			sqlServerDatabasesMu.Lock()
			sqlServerDatabases = append(sqlServerDatabases[:i], sqlServerDatabases[i+1:]...)
			sqlServerDatabasesMu.Unlock()
			break
		}
	}

	log.Printf("Database %s dropped", dbName)

	// Stop the container if no more databases are in use
	if len(sqlServerDatabases) == 0 {
		stopSQLServer()
	}
}
