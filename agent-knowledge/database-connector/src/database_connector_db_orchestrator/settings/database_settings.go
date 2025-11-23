package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type DatabaseSettings struct {
	Url string
}

func NewDatabaseSettings(
	logger *zap.Logger,
) *DatabaseSettings {
	databaseUrl := os.Getenv("DATABASE_URL")

	if len(strings.TrimSpace(databaseUrl)) == 0 {
		logger.Fatal("please set your DATABASE_URL value in your environment")
	}

	return &DatabaseSettings{
		Url: databaseUrl,
	}
}
