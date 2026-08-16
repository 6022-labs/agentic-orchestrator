package settings

import (
	"strings"

	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ConversationBrokerSettingsKey = "conversation_broker"

type ConversationBrokerSettings struct {
	BaseUrl               string    `koanf:"base_url"`
	CallbackUrl           string    `koanf:"callback_url"`
	BridgeID              uuid.UUID `koanf:"bridge_id"`
	RequestTimeoutSeconds uint64    `koanf:"request_timeout_seconds"`
}

func NewConversationBrokerSettings(logger *zap.Logger, k *koanf.Koanf) *ConversationBrokerSettings {
	settings := ConversationBrokerSettings{
		RequestTimeoutSeconds: 30,
	}
	if err := k.Unmarshal(ConversationBrokerSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal conversation_broker settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)
	settings.CallbackUrl = strings.TrimSpace(settings.CallbackUrl)

	if settings.BaseUrl == "" {
		logger.Fatal("please set your conversation_broker.base_url value in your configuration")
	}
	if settings.CallbackUrl == "" {
		logger.Fatal("please set your conversation_broker.callback_url value in your configuration")
	}
	if settings.BridgeID == uuid.Nil {
		logger.Fatal("please set your conversation_broker.bridge_id value in your configuration (any stable UUID)")
	}
	if settings.RequestTimeoutSeconds == 0 {
		settings.RequestTimeoutSeconds = 30
	}
	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")
	return &settings
}
