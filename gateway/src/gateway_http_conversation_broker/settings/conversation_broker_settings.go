package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ConversationBrokerSettingsKey = "conversation_broker"

type ConversationBrokerSettings struct {
	BaseUrl               string `koanf:"base_url"`
	RequestTimeoutSeconds uint64 `koanf:"request_timeout_seconds"`
}

func NewConversationBrokerSettings(logger *zap.Logger, k *koanf.Koanf) *ConversationBrokerSettings {
	settings := ConversationBrokerSettings{
		RequestTimeoutSeconds: 30,
	}
	if err := k.Unmarshal(ConversationBrokerSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal conversation broker settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimRight(strings.TrimSpace(settings.BaseUrl), "/")

	if settings.BaseUrl == "" {
		logger.Fatal("please set your conversation_broker.base_url value in your configuration")
	}
	if settings.RequestTimeoutSeconds == 0 {
		settings.RequestTimeoutSeconds = 30
	}

	return &settings
}
