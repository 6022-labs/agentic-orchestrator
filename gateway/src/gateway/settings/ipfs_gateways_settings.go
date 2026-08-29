package settings

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const IpfsGatewaysSettingsKey = "ipfs"

//go:embed ipfs_gateways.json
var ipfsGatewaysJson []byte

type ipfsGatewaysFile struct {
	Templates []string `json:"templates"`
}

// IpfsGatewaysSettings lists the gateway URL templates the console tries in order.
type IpfsGatewaysSettings struct {
	GatewayUrlTemplates []string `koanf:"gateway_url_templates"`
}

func NewIpfsGatewaysSettings(logger *zap.Logger, k *koanf.Koanf) *IpfsGatewaysSettings {
	// Not seeded: mapstructure merges into an existing slice element-wise.
	var settings IpfsGatewaysSettings
	if err := k.Unmarshal(IpfsGatewaysSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ipfs settings", zap.Error(err))
	}

	settings.normalizeTemplates()
	if len(settings.GatewayUrlTemplates) == 0 {
		settings.applyBuiltInTemplates(logger)
	}

	return &settings
}

func (s *IpfsGatewaysSettings) applyBuiltInTemplates(logger *zap.Logger) {
	var file ipfsGatewaysFile
	if err := json.Unmarshal(ipfsGatewaysJson, &file); err != nil {
		logger.Fatal("failed to parse built-in ipfs_gateways.json", zap.Error(err))
	}

	s.GatewayUrlTemplates = file.Templates
	s.normalizeTemplates()
}

// Env hands the whole list through as one comma-separated string.
func (s *IpfsGatewaysSettings) normalizeTemplates() {
	normalized := make([]string, 0, len(s.GatewayUrlTemplates))
	for _, template := range s.GatewayUrlTemplates {
		for _, part := range strings.Split(template, ",") {
			if trimmed := strings.TrimSpace(part); len(trimmed) != 0 {
				normalized = append(normalized, trimmed)
			}
		}
	}

	s.GatewayUrlTemplates = normalized
}
