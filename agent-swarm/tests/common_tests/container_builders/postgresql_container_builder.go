package container_builders

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/certificate_builders"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresContainer   *PostgresContainerBuilder
	postgresContainerMu sync.Mutex
)

type PostgresContainerBuilder struct {
	container       testcontainers.Container
	connStrTemplate string
	containerMu     sync.Mutex

	tlsContainer       testcontainers.Container
	tlsConnStrTemplate string
	caCertificatePem   []byte
	tlsContainerMu     sync.Mutex

	databasesMu sync.Mutex
	databases   []string
}

func GetPostgresContainerBuilder() *PostgresContainerBuilder {
	postgresContainerMu.Lock()
	defer postgresContainerMu.Unlock()

	if postgresContainer == nil {
		postgresContainer = &PostgresContainerBuilder{}
	}
	return postgresContainer
}

func (b *PostgresContainerBuilder) start() string {
	b.containerMu.Lock()
	defer b.containerMu.Unlock()

	if b.container == nil {
		var connStrTemplate string
		b.container, connStrTemplate = b.run(b.request())
		b.connStrTemplate = connStrTemplate + " sslmode=disable"
		log.Println("PostgreSQL container ready")
	}
	return b.connStrTemplate
}

func (b *PostgresContainerBuilder) startTls() (string, []byte) {
	b.tlsContainerMu.Lock()
	defer b.tlsContainerMu.Unlock()

	if b.tlsContainer == nil {
		certificateAuthority := certificate_builders.NewCertificateAuthority()
		certificatePem, keyPem := certificateAuthority.IssueServerCertificate(TlsServerName)

		req := b.request()
		req.Files = []testcontainers.ContainerFile{
			{Reader: bytes.NewReader(certificatePem), ContainerFilePath: "/certs/server.crt", FileMode: 0o644},
			{Reader: bytes.NewReader(keyPem), ContainerFilePath: "/certs/server.key", FileMode: 0o644},
		}
		// postgres refuses a key file it does not own with group/other access.
		req.Entrypoint = []string{
			"sh",
			"-c",
			"install -o postgres -g postgres -m 600 /certs/server.key /tmp/server.key" +
				" && exec docker-entrypoint.sh postgres -c ssl=on" +
				" -c ssl_cert_file=/certs/server.crt -c ssl_key_file=/tmp/server.key",
		}

		var connStrTemplate string
		b.tlsContainer, connStrTemplate = b.run(req)
		b.tlsConnStrTemplate = connStrTemplate + " sslmode=require"
		b.caCertificatePem = certificateAuthority.CertificatePem
		log.Println("PostgreSQL TLS container ready")
	}
	return b.tlsConnStrTemplate, b.caCertificatePem
}

func (b *PostgresContainerBuilder) request() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        "postgres:18-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
}

// run starts req and returns the container with a connection string template missing only its sslmode.
func (b *PostgresContainerBuilder) run(req testcontainers.ContainerRequest) (testcontainers.Container, string) {
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	return container, fmt.Sprintf("host=%s port=%s dbname=%%s user=postgres password=postgres", host, port.Port())
}

func (b *PostgresContainerBuilder) stop() {
	b.containerMu.Lock()
	if b.container != nil {
		_ = b.container.Terminate(context.Background())
		b.container = nil
		log.Println("PostgreSQL container stopped")
	}
	b.containerMu.Unlock()

	b.tlsContainerMu.Lock()
	if b.tlsContainer != nil {
		_ = b.tlsContainer.Terminate(context.Background())
		b.tlsContainer = nil
		log.Println("PostgreSQL TLS container stopped")
	}
	b.tlsContainerMu.Unlock()

	postgresContainerMu.Lock()
	postgresContainer = nil
	postgresContainerMu.Unlock()
}

// Create returns (dbName, connectionString) for a scratch DB in the shared container, so tests run in parallel.
func (b *PostgresContainerBuilder) Create() (string, string) {
	dnsTemplate := b.start()
	dbName := b.createDatabase(dnsTemplate)
	return dbName, fmt.Sprintf(dnsTemplate, dbName)
}

func (b *PostgresContainerBuilder) CreateTls() (string, string, []byte) {
	dnsTemplate, caCertificatePem := b.startTls()
	dbName := b.createDatabase(dnsTemplate)
	return dbName, fmt.Sprintf(dnsTemplate, dbName), caCertificatePem
}

func (b *PostgresContainerBuilder) createDatabase(dnsTemplate string) string {
	conn, err := sql.Open("postgres", fmt.Sprintf(dnsTemplate, "postgres"))
	if err != nil {
		log.Fatalf("failed to open postgres: %v", err)
	}
	defer conn.Close()

	dbName := "test_" + uuid.New().String()[:8]
	if _, err := conn.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, dbName)); err != nil {
		log.Fatalf("failed to create db %s: %v", dbName, err)
	}

	b.databasesMu.Lock()
	b.databases = append(b.databases, dbName)
	b.databasesMu.Unlock()

	return dbName
}

func (b *PostgresContainerBuilder) DropDatabase(dbName string) {
	b.databasesMu.Lock()
	defer b.databasesMu.Unlock()
	for i, db := range b.databases {
		if db == dbName {
			b.databases = append(b.databases[:i], b.databases[i+1:]...)
			break
		}
	}
}
