package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const IPFSSettingsKey = "ipfs"

type IPFSSettings struct {
	IpfsNodeHttpUrl string `mapstructure:"node_http_url"`
}

func NewIPFSSettings(logger *zap.Logger, v *viper.Viper) *IPFSSettings {
	if val, ok := os.LookupEnv("IPFS__NODE_HTTP_URL"); ok {
		v.Set(IPFSSettingsKey+".node_http_url", val)
	}

	var settings IPFSSettings
	if err := v.UnmarshalKey(IPFSSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ipfs settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.IpfsNodeHttpUrl)) == 0 {
		logger.Fatal("please set your ipfs.node_http_url value in your configuration")
	}

	return &settings
}
