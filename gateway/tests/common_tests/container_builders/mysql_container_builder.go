package container_builders

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	mysqlContainer       testcontainers.Container
	mysqlConnStrTemplate string
	mysqlContainerMu     sync.Mutex
	mysqlDatabasesMu     sync.Mutex
	mysqlDatabases       []string
)

// Start the container once for all test packages
func startMySQL() string {
	mysqlContainerMu.Lock()
	defer mysqlContainerMu.Unlock()

	if mysqlContainer == nil {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "mysql:latest",
			ExposedPorts: []string{"3306/tcp"},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": "root",
				"MYSQL_DATABASE":      "test",
			},
			WaitingFor: wait.ForListeningPort("3306/tcp"),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("Failed to start MySQL container: %v", err)
		}
		mysqlContainer = container

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "3306")

		// MySQL DSN: "user:password@tcp(host:port)/dbname"
		mysqlConnStrTemplate = fmt.Sprintf("root:root@tcp(%s:%s)/%s?parseTime=true", host, port.Port(), "%s")

		log.Println("MySQL container started and ready")
	}

	return mysqlConnStrTemplate
}

// Stop the container after all tests
func stopMySQL() {
	mysqlContainerMu.Lock()
	defer mysqlContainerMu.Unlock()

	if mysqlContainer != nil {
		mysqlContainer.Terminate(context.Background())
		mysqlContainer = nil

		log.Println("MySQL container stopped")
	}
}

func CreateMySQLDatabase() (string, string) {
	// Connect as root user to create a new database
	dnsTemplate := startMySQL()
	mainDatabaseConnectionString := fmt.Sprintf(dnsTemplate, "test")

	conn, err := sql.Open("mysql", mainDatabaseConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer conn.Close()

	dbName := uuid.New().String()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE `%s`;", dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}

	mysqlDatabasesMu.Lock()
	mysqlDatabases = append(mysqlDatabases, dbName)
	mysqlDatabasesMu.Unlock()

	return dbName, fmt.Sprintf(dnsTemplate, dbName)
}

func DestroyMySQLDatabase(dbName string) {
	// Don't need to remove the database because it will be removed when the container is stopped

	// Just remove the database from the list
	for i, db := range mysqlDatabases {
		if db == dbName {
			mysqlDatabasesMu.Lock()
			mysqlDatabases = append(mysqlDatabases[:i], mysqlDatabases[i+1:]...)
			mysqlDatabasesMu.Unlock()
			break
		}
	}

	log.Printf("Database %s dropped", dbName)

	if len(mysqlDatabases) == 0 {
		stopMySQL()
	}
}
