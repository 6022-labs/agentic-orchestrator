package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type OrchestratorSlackBridgeDatabaseSettings struct {
	Url string
}

func NewOrchestratorSlackBridgeDatabaseSettings(
	logger *zap.Logger,
) *OrchestratorSlackBridgeDatabaseSettings {
	databaseUrl := os.Getenv("DATABASE_URL")

	if len(strings.TrimSpace(databaseUrl)) == 0 {
		logger.Fatal("please set your DATABASE_URL value in your environment")
	}

	return &OrchestratorSlackBridgeDatabaseSettings{
		Url: databaseUrl,
	}
}
