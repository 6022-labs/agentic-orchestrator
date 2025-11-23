package container_builders

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresContainer       testcontainers.Container
	postgresConnStrTemplate string
	postgresContainerMu     sync.Mutex
	postgresDatabasesMu     sync.Mutex
	postgresDatabases       []string
)

// Start the container once for all test packages
func startPostgres() string {
	postgresContainerMu.Lock()
	defer postgresContainerMu.Unlock()

	if postgresContainer == nil {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("Failed to start container: %v", err)
		}
		postgresContainer = container

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "5432")

		postgresConnStrTemplate = fmt.Sprintf("host=%s port=%s dbname=%s user=postgres password=postgres sslmode=disable", host, port.Port(), "%s")

		log.Println("PostgreSQL container started and ready")
	}

	return postgresConnStrTemplate
}

// Stop the container after all tests
func stopPostgres() {
	postgresContainerMu.Lock()
	defer postgresContainerMu.Unlock()

	if postgresContainer != nil {
		postgresContainer.Terminate(context.Background())
		postgresContainer = nil

		log.Println("PostgreSQL container stopped")
	}
}

func CreatePostgreSQLDatabase() (string, string) {
	// Connect as postgres user to create a new database
	dnsTemplate := startPostgres()
	mainDatabaseConnectionString := fmt.Sprintf(dnsTemplate, "postgres")

	conn, err := sql.Open("postgres", mainDatabaseConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	dbName := uuid.New().String()
	_, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}

	postgresDatabasesMu.Lock()
	postgresDatabases = append(postgresDatabases, dbName)
	postgresDatabasesMu.Unlock()

	return dbName, fmt.Sprintf(dnsTemplate, dbName)
}

func DestroyPostgreSQLDatabase(dbName string) {
	// Don't need to remove the database because it will be removed when the container is stopped

	// Just remove the database from the list
	for i, db := range postgresDatabases {
		if db == dbName {
			postgresDatabasesMu.Lock()
			postgresDatabases = append(postgresDatabases[:i], postgresDatabases[i+1:]...)
			postgresDatabasesMu.Unlock()
			break
		}
	}

	log.Printf("Database %s dropped", dbName)

	if len(postgresDatabases) == 0 {
		stopPostgres()
	}
}
