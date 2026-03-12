package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const RpcNodeSettingsKey = "rpc_node"

type RpcNodeSettings struct {
	WsUrl      string `koanf:"ws_url"`
	HttpUrl    string `koanf:"http_url"`
	ENSHttpUrl string `koanf:"ens_http_url"`
}

func NewRpcNodeSettings(logger *zap.Logger, k *koanf.Koanf) *RpcNodeSettings {
	var settings RpcNodeSettings
	if err := k.Unmarshal(RpcNodeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal rpc node settings", zap.Error(err))
	}

	settings.WsUrl = strings.TrimSpace(settings.WsUrl)
	settings.HttpUrl = strings.TrimSpace(settings.HttpUrl)
	settings.ENSHttpUrl = strings.TrimSpace(settings.ENSHttpUrl)

	if len(settings.WsUrl) == 0 {
		logger.Fatal("please set your rpc_node.ws_url value in your configuration")
	}

	if len(settings.HttpUrl) == 0 {
		logger.Fatal("please set your rpc_node.http_url value in your configuration")
	}

	if len(settings.ENSHttpUrl) == 0 {
		logger.Fatal("please set your rpc_node.ens_http_url value in your configuration")
	}

	return &settings
}
