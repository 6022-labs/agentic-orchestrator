package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const GatewaySettingsKey = "gateway"

const defaultPublicPrefix = "/api"

type GatewaySettings struct {
	// PublicPrefix is the path a reverse proxy mounts the gateway under, as the browser sees it.
	PublicPrefix string `koanf:"public_prefix"`
}

func NewGatewaySettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *GatewaySettings {
	settings := GatewaySettings{PublicPrefix: defaultPublicPrefix}
	if err := k.Unmarshal(GatewaySettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal gateway settings", zap.Error(err))
	}

	settings.PublicPrefix = "/" + strings.Trim(strings.TrimSpace(settings.PublicPrefix), "/")

	return &settings
}
