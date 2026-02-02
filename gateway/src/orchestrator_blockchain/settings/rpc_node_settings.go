package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const RpcNodeSettingsKey = "rpc_node"

type RpcNodeSettings struct {
	WsUrl      string `mapstructure:"ws_url"`
	HttpUrl    string `mapstructure:"http_url"`
	ENSHttpUrl string `mapstructure:"ens_http_url"`
}

func NewRpcNodeSettings(logger *zap.Logger, v *viper.Viper) *RpcNodeSettings {
	if val, ok := os.LookupEnv("RPC_NODE__WS_URL"); ok {
		v.Set(RpcNodeSettingsKey+".ws_url", val)
	}
	if val, ok := os.LookupEnv("RPC_NODE__HTTP_URL"); ok {
		v.Set(RpcNodeSettingsKey+".http_url", val)
	}
	if val, ok := os.LookupEnv("RPC_NODE__ENS_HTTP_URL"); ok {
		v.Set(RpcNodeSettingsKey+".ens_http_url", val)
	}

	var settings RpcNodeSettings
	if err := v.UnmarshalKey(RpcNodeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal rpc node settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.WsUrl)) == 0 {
		logger.Fatal("please set your rpc_node.ws_url value in your configuration")
	}

	if len(strings.TrimSpace(settings.HttpUrl)) == 0 {
		logger.Fatal("please set your rpc_node.http_url value in your configuration")
	}

	if len(strings.TrimSpace(settings.ENSHttpUrl)) == 0 {
		logger.Fatal("please set your rpc_node.ens_http_url value in your configuration")
	}

	return &settings
}
