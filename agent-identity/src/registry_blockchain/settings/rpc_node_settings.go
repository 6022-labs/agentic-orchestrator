package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const RpcNodeSettingsKey = "rpc_node"

type RpcNodeSettings struct {
	HttpUrl string `koanf:"http_url"`
}

func NewRpcNodeSettings(logger *zap.Logger, k *koanf.Koanf) *RpcNodeSettings {
	var settings RpcNodeSettings
	if err := k.Unmarshal(RpcNodeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal rpc node settings", zap.Error(err))
	}

	settings.HttpUrl = strings.TrimSpace(settings.HttpUrl)

	if len(settings.HttpUrl) == 0 {
		logger.Fatal("please set your rpc_node.http_url value in your configuration")
	}

	return &settings
}
