package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const IPFSSettingsKey = "ipfs"

// IPFSSettings targets a Kubo node's HTTP API. BearerToken (optional) is sent as
// `Authorization: Bearer <token>` for Kubo's API.Authorizations scheme.
type IPFSSettings struct {
	NodeHttpUrl string `koanf:"node_http_url"`
	BearerToken string `koanf:"bearer_token"`
}

func NewIPFSSettings(logger *zap.Logger, k *koanf.Koanf) *IPFSSettings {
	var settings IPFSSettings
	if err := k.Unmarshal(IPFSSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ipfs settings", zap.Error(err))
	}

	settings.NodeHttpUrl = strings.TrimSpace(settings.NodeHttpUrl)
	if len(settings.NodeHttpUrl) == 0 {
		logger.Fatal("please set your ipfs.node_http_url value in your configuration")
	}

	settings.BearerToken = strings.TrimSpace(settings.BearerToken)

	return &settings
}
