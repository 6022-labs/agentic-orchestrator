package container_builders

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresTLSServerName is the only name in the TLS container's certificate, never the host it is reached on.
const PostgresTLSServerName = "postgres.internal"

var (
	postgresTLSContainer        testcontainers.Container
	postgresTLSConnStrTemplate  string
	postgresTLSCaCertificatePEM string
	postgresTLSContainerMu      sync.Mutex
)

// Postgres refuses a key file it does not own, so the entrypoint copies the mounted pair before starting.
const postgresTLSEntrypoint = "mkdir -p /etc/postgresql-tls && " +
	"cp /tls/server.crt /tls/server.key /etc/postgresql-tls/ && " +
	"chown postgres:postgres /etc/postgresql-tls/server.crt /etc/postgresql-tls/server.key && " +
	"chmod 600 /etc/postgresql-tls/server.key && " +
	"exec docker-entrypoint.sh postgres -c ssl=on " +
	"-c ssl_cert_file=/etc/postgresql-tls/server.crt -c ssl_key_file=/etc/postgresql-tls/server.key"

func startTLSPostgres() (string, string) {
	postgresTLSContainerMu.Lock()
	defer postgresTLSContainerMu.Unlock()

	if postgresTLSContainer == nil {
		ctx := context.Background()
		certificateAuthority := common_tests.NewCertificateAuthority()
		certificatePEM, privateKeyPEM := certificateAuthority.IssueServerCertificate(PostgresTLSServerName)

		req := testcontainers.ContainerRequest{
			Image:        "postgres:18-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres",
			},
			Files: []testcontainers.ContainerFile{
				{Reader: strings.NewReader(certificatePEM), ContainerFilePath: "/tls/server.crt", FileMode: 0o644},
				{Reader: strings.NewReader(privateKeyPEM), ContainerFilePath: "/tls/server.key", FileMode: 0o600},
			},
			Entrypoint: []string{"sh", "-c", postgresTLSEntrypoint},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("Failed to start TLS container: %v", err)
		}
		postgresTLSContainer = container
		postgresTLSCaCertificatePEM = certificateAuthority.CertificatePEM

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "5432")

		postgresTLSConnStrTemplate = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=postgres password=postgres",
			host,
			port.Port(),
			"%s",
		)

		log.Println("PostgreSQL TLS container started and ready")
	}

	return postgresTLSConnStrTemplate, postgresTLSCaCertificatePEM
}

func stopTLSPostgres() {
	postgresTLSContainerMu.Lock()
	defer postgresTLSContainerMu.Unlock()

	if postgresTLSContainer != nil {
		postgresTLSContainer.Terminate(context.Background())
		postgresTLSContainer = nil

		log.Println("PostgreSQL TLS container stopped")
	}
}

// CreateTLSPostgreSQLDatabase returns a database name, its connection string and the CA that signed the server.
func CreateTLSPostgreSQLDatabase() (string, string, string) {
	dnsTemplate, caCertificatePEM := startTLSPostgres()

	conn, err := sql.Open("postgres", fmt.Sprintf(dnsTemplate, "postgres")+" sslmode=require")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL over TLS: %v", err)
	}
	defer conn.Close()

	dbName := uuid.New().String()
	if _, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, dbName)); err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}

	return dbName, fmt.Sprintf(dnsTemplate, dbName), caCertificatePEM
}
