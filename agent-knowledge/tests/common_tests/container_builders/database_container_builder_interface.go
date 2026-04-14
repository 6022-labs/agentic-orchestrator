package container_builders

type DatabaseContainerBuilderInterface interface {
	Create() (string, string)
	DropDatabase(dbName string)
	GetHost() string
	GetPort() uint16
}
