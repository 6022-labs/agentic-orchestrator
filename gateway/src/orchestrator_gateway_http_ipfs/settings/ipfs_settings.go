package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const IPFSSettingsKey = "ipfs"

type IPFSSettings struct {
	IpfsNodeHttpUrl string `koanf:"node_http_url"`
}

func NewIPFSSettings(logger *zap.Logger, k *koanf.Koanf) *IPFSSettings {
	var settings IPFSSettings
	if err := k.Unmarshal(IPFSSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ipfs settings", zap.Error(err))
	}

	settings.IpfsNodeHttpUrl = strings.TrimSpace(settings.IpfsNodeHttpUrl)

	if len(settings.IpfsNodeHttpUrl) == 0 {
		logger.Fatal("please set your ipfs.node_http_url value in your configuration")
	}

	return &settings
}
