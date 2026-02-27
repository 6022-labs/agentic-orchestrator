package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const IpfsGatewaySettingsKey = "ipfs_gateway"

type IpfsGatewaySettings struct {
	GatewayUrlTemplate string `koanf:"url_template"`
}

func NewIpfsGatewaySettings(logger *zap.Logger, k *koanf.Koanf) *IpfsGatewaySettings {
	settings := IpfsGatewaySettings{
		GatewayUrlTemplate: "https://ipfs.io/ipfs/%s",
	}
	if err := k.Unmarshal(IpfsGatewaySettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ipfs gateway settings", zap.Error(err))
	}

	settings.GatewayUrlTemplate = strings.TrimSpace(settings.GatewayUrlTemplate)

	if len(settings.GatewayUrlTemplate) == 0 {
		settings.GatewayUrlTemplate = "https://ipfs.io/ipfs/%s"
	}

	return &settings
}
