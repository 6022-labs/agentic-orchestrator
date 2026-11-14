package factories_test

import (
	"crypto/x509"
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/factories"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/certificate_builders"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	db_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_db_mocks/metrics_mocks/interfaces_mocks"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite struct {
	dbName string

	sut *factories.DatabaseClientFactory
}

func WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(
	t *testing.T,
	dbName string,
	databaseConfiguration map[string]any,
) *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
	k := koanf.New(".")
	_ = k.Load(confmap.Provider(databaseConfiguration, "."), nil)
	databaseSettings := settings.NewDatabaseSettings(zap.NewNop(), k)

	migrationFolderPath, err := common_tests.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}
	migrationsSettings := &settings.MigrationsSettings{
		MigrationsFolderPath: *migrationFolderPath,
	}

	dbQueryMetrics := db_metrics_mocks.NewMockDbQueryMetricsInterface(gomock.NewController(t))
	dbQueryMetrics.EXPECT().
		RecordQuery(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes()

	sut := factories.NewDatabaseClientFactory(
		zap.NewNop(),
		dbQueryMetrics,
		migrationsSettings,
		databaseSettings,
	)

	return &WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite{
		sut:    sut,
		dbName: dbName,
	}
}

func WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(
	suite *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite,
) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

func TestWhenCreatingOrchestratorSlackBridgeDbClient(t *testing.T) {
	t.Parallel()

	t.Run("Given no CA certificate", func(t *testing.T) {
		t.Parallel()

		initSuite := func(t *testing.T) *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
			dbName, connectionString := container_builders.GetPostgresContainerBuilder().Create()
			return WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(t, dbName, map[string]any{
				"database.url": connectionString,
			})
		}

		t.Run("Should return a new orchestrator db client", func(t *testing.T) {
			t.Parallel()

			suite := initSuite(t)
			defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

			dbClient, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()

			assert.Nil(t, err)
			assert.NotNil(t, dbClient)
		})
	})

	t.Run("Given a database without TLS", func(t *testing.T) {
		t.Parallel()

		initSuite := func(t *testing.T) *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
			dbName, connectionString := container_builders.GetPostgresContainerBuilder().Create()
			return WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(t, dbName, map[string]any{
				// prefer is the sslmode whose pgx fallback would retry in plaintext.
				"database.url":             strings.Replace(connectionString, "sslmode=disable", "sslmode=prefer", 1),
				"database.ca_certificate":  string(certificate_builders.NewCertificateAuthority().CertificatePem),
				"database.tls_server_name": container_builders.TlsServerName,
			})
		}

		t.Run("Should refuse to connect in plaintext", func(t *testing.T) {
			t.Parallel()

			suite := initSuite(t)
			defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

			_, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()

			assert.ErrorContains(t, err, "server refused TLS connection")
		})
	})

	t.Run("Given another CA certificate", func(t *testing.T) {
		t.Parallel()

		initSuite := func(t *testing.T) *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
			dbName, connectionString, _ := container_builders.GetPostgresContainerBuilder().CreateTls()
			return WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(t, dbName, map[string]any{
				"database.url":             connectionString,
				"database.ca_certificate":  string(certificate_builders.NewCertificateAuthority().CertificatePem),
				"database.tls_server_name": container_builders.TlsServerName,
			})
		}

		t.Run("Should refuse the server certificate", func(t *testing.T) {
			t.Parallel()

			suite := initSuite(t)
			defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

			_, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()

			var unknownAuthorityError x509.UnknownAuthorityError
			assert.ErrorAs(t, err, &unknownAuthorityError)
		})
	})

	t.Run("Given direct TLS negotiation", func(t *testing.T) {
		t.Parallel()

		initSuite := func(t *testing.T) *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
			dbName, connectionString, caCertificatePem := container_builders.GetPostgresContainerBuilder().CreateTls()
			return WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(t, dbName, map[string]any{
				// direct negotiation is refused without the postgresql ALPN protocol.
				"database.url":             connectionString + " sslnegotiation=direct",
				"database.ca_certificate":  string(caCertificatePem),
				"database.tls_server_name": container_builders.TlsServerName,
			})
		}

		t.Run("Should connect over TLS", func(t *testing.T) {
			t.Parallel()

			suite := initSuite(t)
			defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

			dbClient, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()
			require.NoError(t, err)

			var ssl bool
			dbClient.GetDB().Raw("SELECT ssl FROM pg_stat_ssl WHERE pid = pg_backend_pid()").Scan(&ssl)
			assert.True(t, ssl)
		})
	})

	t.Run("Given the CA certificate of the database", func(t *testing.T) {
		t.Parallel()

		initSuite := func(t *testing.T) *WhenCreatingOrchestratorSlackBridgeDbClientTestingSuite {
			dbName, connectionString, caCertificatePem := container_builders.GetPostgresContainerBuilder().CreateTls()
			return WhenCreatingOrchestratorSlackBridgeDbClientBeforeEach(t, dbName, map[string]any{
				"database.url":             connectionString,
				"database.ca_certificate":  string(caCertificatePem),
				"database.tls_server_name": container_builders.TlsServerName,
			})
		}

		t.Run("Should connect over TLS", func(t *testing.T) {
			t.Parallel()

			suite := initSuite(t)
			defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

			dbClient, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()
			require.NoError(t, err)

			var ssl bool
			dbClient.GetDB().Raw("SELECT ssl FROM pg_stat_ssl WHERE pid = pg_backend_pid()").Scan(&ssl)
			assert.True(t, ssl)
		})

		t.Run("Should apply the migrations", func(t *testing.T) {
			t.Parallel()

			suite := initSuite(t)
			defer WhenCreatingOrchestratorSlackBridgeDbClientAfterEach(suite)

			dbClient, err := suite.sut.CreateOrchestratorSlackBridgeDBClient()
			require.NoError(t, err)

			assert.NoError(t, dbClient.Migrate())
		})
	})
}
