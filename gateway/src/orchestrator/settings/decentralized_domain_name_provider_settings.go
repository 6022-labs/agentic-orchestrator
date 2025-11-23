package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type DecentralizedDomainNameProviderSettings struct {
	DomainExtension string
}

func NewDecentralizedDomainNameProviderSettings(logger *zap.Logger) *DecentralizedDomainNameProviderSettings {
	domainExtension := os.Getenv("DECENTRALIZED_DOMAIN_EXTENSION")
	if len(strings.TrimSpace(domainExtension)) == 0 {
		logger.Warn("DECENTRALIZED_DOMAIN_EXTENSION environment variable is not set, using default value, 'eth'")
		domainExtension = "eth"
	}

	return &DecentralizedDomainNameProviderSettings{
		DomainExtension: domainExtension,
	}
}
