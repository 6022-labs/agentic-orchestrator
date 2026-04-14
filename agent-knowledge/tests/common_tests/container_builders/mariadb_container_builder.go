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
	mariadbContainer *MariaDBContainerBuilder
)

type MariaDBContainerBuilder struct {
	container       testcontainers.Container
	connStrTemplate string
	containerMu     sync.Mutex
	databasesMu     sync.Mutex
	databases       []string
}

func GetMariaDBContainerBuilder() *MariaDBContainerBuilder {
	if mariadbContainer == nil {
		mariadbContainer = &MariaDBContainerBuilder{}
	}
	return mariadbContainer
}

func (b *MariaDBContainerBuilder) start() string {
	b.containerMu.Lock()
	defer b.containerMu.Unlock()

	if b.container == nil {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "mariadb:latest",
			ExposedPorts: []string{"3306/tcp"},
			Env: map[string]string{
				"MARIADB_ROOT_PASSWORD": "root",
				"MARIADB_DATABASE":      "mariadb",
			},
			WaitingFor: wait.ForListeningPort("3306/tcp"),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("Failed to start MariaDB container: %v", err)
		}
		b.container = container

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "3306")

		b.connStrTemplate = fmt.Sprintf("%s:%s@tcp(%s:%s)/%%s?parseTime=true", "root", "root", host, port.Port())

		log.Println("MariaDB container started and ready")
	}

	return b.connStrTemplate
}

func (b *MariaDBContainerBuilder) stop() {
	b.containerMu.Lock()
	defer b.containerMu.Unlock()

	if b.container != nil {
		b.container.Terminate(context.Background())
		b.container = nil
		log.Println("MariaDB container stopped")
	}
	mariadbContainer = nil
}

func (b *MariaDBContainerBuilder) Create() (string, string) {
	dnsTemplate := b.start()
	mainDatabaseConnectionString := fmt.Sprintf(dnsTemplate, "mariadb")

	conn, err := sql.Open("mysql", mainDatabaseConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to MariaDB: %v", err)
	}
	defer conn.Close()

	dbName := uuid.New().String()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE `%s`;", dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}

	b.databasesMu.Lock()
	b.databases = append(b.databases, dbName)
	b.databasesMu.Unlock()

	dbConnStr := fmt.Sprintf(dnsTemplate, dbName)
	return dbName, dbConnStr
}

func (b *MariaDBContainerBuilder) DropDatabase(dbName string) {
	b.databasesMu.Lock()
	defer b.databasesMu.Unlock()

	for i, db := range b.databases {
		if db == dbName {
			b.databases = append(b.databases[:i], b.databases[i+1:]...)
			break
		}
	}
	log.Printf("Database %s dropped", dbName)
}

func (b *MariaDBContainerBuilder) GetHost() string {
	host, err := b.container.Host(context.Background())
	if err != nil {
		panic(err)
	}

	return host
}

func (b *MariaDBContainerBuilder) GetPort() uint16 {
	portNat, err := b.container.MappedPort(context.Background(), "3306")
	if err != nil {
		panic(err)
	}

	return portNat.Num()
}
