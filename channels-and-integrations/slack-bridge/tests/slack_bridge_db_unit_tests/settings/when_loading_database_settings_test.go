package settings_test

import (
	"crypto/tls"
	"crypto/x509"
	"testing"

	db_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/certificate_builders"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestWhenLoadingDatabaseSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured url", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"database.url": "postgres://user:pass@localhost:5432/db",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "postgres://user:pass@localhost:5432/db", settings.Url)
		})

		t.Run("Should leave TLS to the url", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

			assert.Nil(t, settings.TlsConfig)
		})
	})

	t.Run("Given a url with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"database.url": "  postgres://user:pass@localhost:5432/db  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

			assert.Equal(t, "postgres://user:pass@localhost:5432/db", settings.Url)
		})
	})

	t.Run("Given an invalid CA certificate", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"database.url":            "postgres://user:pass@localhost:5432/db",
			"database.ca_certificate": "not a certificate",
		}, "."), nil)

		t.Run("Should exit", func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop().WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic))

			assert.Panics(t, func() { db_settings.NewDatabaseSettings(logger, k) })
		})
	})

	t.Run("Given a TLS server name without a CA certificate", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"database.url":             "postgres://user:pass@localhost:5432/db",
			"database.tls_server_name": "postgres.internal",
		}, "."), nil)

		t.Run("Should exit", func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop().WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic))

			assert.Panics(t, func() { db_settings.NewDatabaseSettings(logger, k) })
		})
	})

	t.Run("Given a CA certificate", func(t *testing.T) {
		t.Parallel()

		caCertificatePem := certificate_builders.NewCertificateAuthority().CertificatePem
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"database.url":            "postgres://user:pass@db.example.com:5432/db",
			"database.ca_certificate": string(caCertificatePem),
		}, "."), nil)

		t.Run("Should trust only it", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

			expectedRootCAs := x509.NewCertPool()
			expectedRootCAs.AppendCertsFromPEM(caCertificatePem)
			assert.True(t, expectedRootCAs.Equal(settings.TlsConfig.RootCAs))
		})

		t.Run("Should require TLS 1.2 or later", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

			assert.Equal(t, uint16(tls.VersionTLS12), settings.TlsConfig.MinVersion)
		})

		t.Run("Should verify the url host", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

			assert.Equal(t, "db.example.com", settings.TlsConfig.ServerName)
		})
	})

	t.Run("Given a CA certificate and a TLS server name", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"database.url":             "postgres://user:pass@proxy.example.com:5432/db",
			"database.ca_certificate":  string(certificate_builders.NewCertificateAuthority().CertificatePem),
			"database.tls_server_name": "postgres.internal",
		}, "."), nil)

		t.Run("Should verify the TLS server name", func(t *testing.T) {
			t.Parallel()

			settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

			assert.Equal(t, "postgres.internal", settings.TlsConfig.ServerName)
		})
	})
}
