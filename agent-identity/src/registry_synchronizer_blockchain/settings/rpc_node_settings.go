package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const RpcNodeSettingsKey = "rpc_node"

// RpcNodeSettings holds the read timeout for http only; ws carries long-lived subscriptions.
type RpcNodeSettings struct {
	WsUrl                   string `koanf:"ws_url"`
	HttpUrl                 string `koanf:"http_url"`
	RequestTimeoutSeconds   uint64 `koanf:"request_timeout_seconds"`
	MaxAddressesPerLogQuery int    `koanf:"max_addresses_per_log_query"`
}

func NewRpcNodeSettings(logger *zap.Logger, k *koanf.Koanf) *RpcNodeSettings {
	settings := RpcNodeSettings{
		RequestTimeoutSeconds:   10,
		MaxAddressesPerLogQuery: 10,
	}
	if err := k.Unmarshal(RpcNodeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal rpc node settings", zap.Error(err))
	}

	settings.WsUrl = strings.TrimSpace(settings.WsUrl)
	settings.HttpUrl = strings.TrimSpace(settings.HttpUrl)

	if len(settings.WsUrl) == 0 {
		logger.Fatal("please set your rpc_node.ws_url value in your configuration")
	}

	if len(settings.HttpUrl) == 0 {
		logger.Fatal("please set your rpc_node.http_url value in your configuration")
	}

	if settings.RequestTimeoutSeconds == 0 {
		settings.RequestTimeoutSeconds = 10
	}

	if settings.MaxAddressesPerLogQuery <= 0 {
		settings.MaxAddressesPerLogQuery = 10
	}

	return &settings
}
