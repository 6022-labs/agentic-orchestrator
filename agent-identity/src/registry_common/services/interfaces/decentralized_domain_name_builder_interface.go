package interfaces

type DecentralizedDomainNameBuilderInterface interface {
	BuildDomainName(agentName, collectionName string) string
}
