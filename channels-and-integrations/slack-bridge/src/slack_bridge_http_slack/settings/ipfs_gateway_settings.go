package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const IpfsGatewaySettingsKey = "ipfs_gateway"

type IpfsGatewaySettings struct {
	GatewayUrlTemplate string `mapstructure:"url_template"`
}

func NewIpfsGatewaySettings(logger *zap.Logger, v *viper.Viper) *IpfsGatewaySettings {
	v.SetDefault(IpfsGatewaySettingsKey+".url_template", "https://ipfs.io/ipfs/%s")

	if val, ok := os.LookupEnv("IPFS_GATEWAY__URL_TEMPLATE"); ok {
		v.Set(IpfsGatewaySettingsKey+".url_template", val)
	}

	var settings IpfsGatewaySettings
	if err := v.UnmarshalKey(IpfsGatewaySettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ipfs gateway settings", zap.Error(err))
	}

	return &settings
}
