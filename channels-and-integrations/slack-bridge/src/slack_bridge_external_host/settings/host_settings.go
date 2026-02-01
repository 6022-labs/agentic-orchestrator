package settings

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type HostSettings struct {
	ApiPort        int
	AppName        string
	BodyLimitBytes int
}

func NewHostSettings(logger *zap.Logger) *HostSettings {
	apiPortStr := os.Getenv("PORT")
	if len(strings.TrimSpace(apiPortStr)) == 0 {
		logger.Warn("environment variable PORT is not set, defaulting to 3000")
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

	bodyLimitBytes := 4 * 1024 * 1024
	bodyLimitBytesStr := os.Getenv("BODY_LIMIT_BYTES")
	if len(strings.TrimSpace(bodyLimitBytesStr)) == 0 {
		logger.Warn("environment variable BODY_LIMIT_BYTES is not set, defaulting to 4MB")
	} else {
		parsedBytes, err := strconv.Atoi(bodyLimitBytesStr)
		if err != nil || parsedBytes <= 0 {
			logger.Warn("invalid BODY_LIMIT_BYTES environment variable value, defaulting to 4MB")
		} else {
			bodyLimitBytes = parsedBytes
		}
	}

	return &HostSettings{
		ApiPort:        apiPort,
		AppName:        appName,
		BodyLimitBytes: bodyLimitBytes,
	}
}
