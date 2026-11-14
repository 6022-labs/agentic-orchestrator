package settings

import (
	"crypto/tls"
	"crypto/x509"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const DatabaseSettingsKey = "database"

type DatabaseSettings struct {
	Url string
	// TlsConfig verifies the server against ca_certificate; nil leaves TLS to the url.
	TlsConfig *tls.Config
}

type databaseSettingsRaw struct {
	Url           string `koanf:"url"`
	CaCertificate string `koanf:"ca_certificate"`
	TlsServerName string `koanf:"tls_server_name"`
}

func NewConversationBrokerDatabaseSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *DatabaseSettings {
	var raw databaseSettingsRaw
	if err := k.Unmarshal(DatabaseSettingsKey, &raw); err != nil {
		logger.Fatal("failed to unmarshal database settings", zap.Error(err))
	}

	raw.Url = strings.TrimSpace(raw.Url)
	raw.CaCertificate = strings.TrimSpace(raw.CaCertificate)
	raw.TlsServerName = strings.TrimSpace(raw.TlsServerName)

	if len(raw.Url) == 0 {
		logger.Fatal("please set your database.url value in your configuration")
	}

	return &DatabaseSettings{
		Url:       raw.Url,
		TlsConfig: raw.tlsConfig(logger),
	}
}

func (r *databaseSettingsRaw) tlsConfig(logger *zap.Logger) *tls.Config {
	if len(r.CaCertificate) == 0 {
		if len(r.TlsServerName) > 0 {
			logger.Fatal("database.tls_server_name requires database.ca_certificate")
		}
		return nil
	}

	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM([]byte(r.CaCertificate)) {
		logger.Fatal("database.ca_certificate holds no valid PEM certificate")
	}

	serverName := r.TlsServerName
	if len(serverName) == 0 {
		config, err := pgconn.ParseConfig(r.Url)
		if err != nil {
			logger.Fatal("failed to read the host from database.url", zap.Error(err))
		}
		serverName = config.Host
	}

	return &tls.Config{RootCAs: rootCAs, ServerName: serverName, MinVersion: tls.VersionTLS12}
}
