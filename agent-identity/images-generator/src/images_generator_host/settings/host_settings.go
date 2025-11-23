package settings

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type HostSettings struct {
	ApiPort int
	AppName string
}

func NewHostSettings(logger *zap.Logger) *HostSettings {
	apiPortStr := os.Getenv("API_PORT")
	if len(strings.TrimSpace(apiPortStr)) == 0 {
		logger.Warn("environment variable API_PORT is not set, defaulting to 3000")
		apiPortStr = "3000"
	}

	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		logger.Warn("invalid PORT environment variable value, defaulting to 3000")
		apiPort = 3000
	}

	appName := os.Getenv("APP_NAME")
	if len(strings.TrimSpace(appName)) == 0 {
		logger.Fatal("please set your APP_NAME value in your environment")
	}

	return &HostSettings{
		ApiPort: apiPort,
		AppName: appName,
	}
}
