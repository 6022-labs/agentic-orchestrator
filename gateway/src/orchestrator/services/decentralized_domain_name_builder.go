package services

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"go.uber.org/zap"
)

type DecentralizedDomainNameBuilder struct {
	logger                                  *zap.Logger
	blockchainSettings                      *settings.BlockchainSettings
	decentralizedDomainNameProviderSettings *settings.DecentralizedDomainNameProviderSettings
}

func NewDecentralizedDomainNameBuilder(
	logger *zap.Logger,
	blockchainSettings *settings.BlockchainSettings,
	decentralizedDomainNameProviderSettings *settings.DecentralizedDomainNameProviderSettings,
) *DecentralizedDomainNameBuilder {
	return &DecentralizedDomainNameBuilder{
		logger:                                  logger,
		blockchainSettings:                      blockchainSettings,
		decentralizedDomainNameProviderSettings: decentralizedDomainNameProviderSettings,
	}
}

func (b *DecentralizedDomainNameBuilder) BuildDomainName(agentName, collectionName string) string {
	domainName := fmt.Sprintf(
		"%s.%s.%d.protocol6022.%s",
		agentName,
		collectionName,
		b.blockchainSettings.ChainId,
		b.decentralizedDomainNameProviderSettings.DomainExtension,
	)

	b.logger.Debug(
		"Builded domain name",
		zap.String("agentName", agentName),
		zap.String("collectionName", collectionName),
		zap.String("domainName", domainName),
	)

	return domainName
}
