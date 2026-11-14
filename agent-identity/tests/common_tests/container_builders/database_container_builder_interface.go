package container_builders

// DatabaseContainerBuilderInterface gives every engine the same scratch-DB flow: create, connect, drop on cleanup.
type DatabaseContainerBuilderInterface interface {
	Create() (string, string)
	DropDatabase(dbName string)
}

// TlsServerName is the only name the TLS containers' server certificates are valid for.
const TlsServerName = "db.test.internal"

// TlsDatabaseContainerBuilderInterface also serves scratch DBs from a TLS container whose certificate its own CA signs.
type TlsDatabaseContainerBuilderInterface interface {
	DatabaseContainerBuilderInterface
	// CreateTls returns (dbName, connectionString, caCertificatePem) for a scratch DB in the TLS container.
	CreateTls() (string, string, []byte)
}
