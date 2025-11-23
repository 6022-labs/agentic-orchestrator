package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type AgentCollectionSettings struct {
	SigningDomain string
}

func NewAgentCollectionSettings(logger *zap.Logger) *AgentCollectionSettings {
	signingDomain := os.Getenv("AGENT_COLLECTION_SIGNING_DOMAIN")
	if len(strings.TrimSpace(signingDomain)) == 0 {
		logger.Fatal("please set your AGENT_COLLECTION_SIGNING_DOMAIN value in your environment")
	}

	return &AgentCollectionSettings{
		SigningDomain: signingDomain,
	}
}
