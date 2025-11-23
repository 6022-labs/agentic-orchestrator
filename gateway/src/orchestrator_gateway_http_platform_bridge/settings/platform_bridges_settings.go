package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type PlatformBridgesSettings struct {
	SlackBridgeUrl string
}

func NewPlatformBridgesSettings(logger *zap.Logger) *PlatformBridgesSettings {
	slackBridgeUrl := os.Getenv("SLACK_BRIDGE_URL")
	if len(strings.TrimSpace(slackBridgeUrl)) == 0 {
		logger.Fatal("SLACK_BRIDGE_URL is not set")
	}

	return &PlatformBridgesSettings{
		SlackBridgeUrl: slackBridgeUrl,
	}
}
