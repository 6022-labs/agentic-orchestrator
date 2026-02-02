package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const HostSettingsKey = "host"

type HostSettings struct {
	UseApi         bool   `mapstructure:"use_api"`
	UseListeners   bool   `mapstructure:"use_listeners"`
	ApiPort        int    `mapstructure:"api_port"`
	AppName        string `mapstructure:"app_name"`
	BodyLimitBytes int    `mapstructure:"body_limit_bytes"`
}

func NewHostSettings(logger *zap.Logger, v *viper.Viper) *HostSettings {
	v.SetDefault(HostSettingsKey+".use_api", true)
	v.SetDefault(HostSettingsKey+".use_listeners", true)
	v.SetDefault(HostSettingsKey+".api_port", 3000)
	v.SetDefault(HostSettingsKey+".body_limit_bytes", 4*1024*1024)

	if val, ok := os.LookupEnv("HOST__USE_API"); ok {
		v.Set(HostSettingsKey+".use_api", val)
	}
	if val, ok := os.LookupEnv("HOST__USE_LISTENERS"); ok {
		v.Set(HostSettingsKey+".use_listeners", val)
	}
	if val, ok := os.LookupEnv("HOST__API_PORT"); ok {
		v.Set(HostSettingsKey+".api_port", val)
	}
	if val, ok := os.LookupEnv("HOST__APP_NAME"); ok {
		v.Set(HostSettingsKey+".app_name", val)
	}
	if val, ok := os.LookupEnv("HOST__BODY_LIMIT_BYTES"); ok {
		v.Set(HostSettingsKey+".body_limit_bytes", val)
	}

	var settings HostSettings
	if err := v.UnmarshalKey(HostSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal host settings", zap.Error(err))
	}

	if settings.ApiPort <= 0 {
		logger.Warn("invalid host.api_port value, defaulting to 3000")
		settings.ApiPort = v.GetInt(HostSettingsKey + ".api_port")
		if settings.ApiPort <= 0 {
			settings.ApiPort = 3000
		}
	}

	if len(strings.TrimSpace(settings.AppName)) == 0 {
		logger.Fatal("please set your host.app_name value in your configuration")
	}

	if settings.BodyLimitBytes <= 0 {
		logger.Warn("invalid host.body_limit_bytes value, defaulting to 4MB")
		settings.BodyLimitBytes = v.GetInt(HostSettingsKey + ".body_limit_bytes")
		if settings.BodyLimitBytes <= 0 {
			settings.BodyLimitBytes = 4 * 1024 * 1024
		}
	}

	return &settings
}
