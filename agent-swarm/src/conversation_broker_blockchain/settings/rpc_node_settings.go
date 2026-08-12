package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const RpcNodeSettingsKey = "rpc_node"

// RpcNodeSettings carries no ws_url: the broker only reads, it never subscribes.
type RpcNodeSettings struct {
	HttpUrl               string `koanf:"http_url"`
	RequestTimeoutSeconds uint64 `koanf:"request_timeout_seconds"`
}

func NewRpcNodeSettings(logger *zap.Logger, k *koanf.Koanf) *RpcNodeSettings {
	settings := RpcNodeSettings{
		RequestTimeoutSeconds: 10,
	}
	if err := k.Unmarshal(RpcNodeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal rpc node settings", zap.Error(err))
	}

	settings.HttpUrl = strings.TrimSpace(settings.HttpUrl)

	if len(settings.HttpUrl) == 0 {
		logger.Fatal("please set your rpc_node.http_url value in your configuration")
	}

	if settings.RequestTimeoutSeconds == 0 {
		settings.RequestTimeoutSeconds = 10
	}

	return &settings
}
