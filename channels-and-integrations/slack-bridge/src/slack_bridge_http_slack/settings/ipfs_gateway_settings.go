package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const IpfsGatewaySettingsKey = "ipfs_gateway"

// Path form: agents pin CIDv0 hashes, which subdomain gateways cannot serve.
const defaultGatewayUrlTemplate = "https://w3s.link/ipfs/{cid}"

// IpfsGatewaySettings is the gateway Slack fetches agent icons through.
type IpfsGatewaySettings struct {
	GatewayUrlTemplate string `koanf:"url_template"`
}

func NewIpfsGatewaySettings(logger *zap.Logger, k *koanf.Koanf) *IpfsGatewaySettings {
	settings := IpfsGatewaySettings{
		GatewayUrlTemplate: defaultGatewayUrlTemplate,
	}
	if err := k.Unmarshal(IpfsGatewaySettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ipfs gateway settings", zap.Error(err))
	}

	settings.GatewayUrlTemplate = strings.TrimSpace(settings.GatewayUrlTemplate)

	if len(settings.GatewayUrlTemplate) == 0 {
		settings.GatewayUrlTemplate = defaultGatewayUrlTemplate
	}

	return &settings
}

// ImageUrl resolves a bare CID, an ipfs:// URI or an already-displayable URL.
func (s *IpfsGatewaySettings) ImageUrl(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}

	// A path suffix can never go into a subdomain template's host.
	cid, path, hasPath := strings.Cut(strings.TrimPrefix(trimmed, "ipfs://"), "/")
	if hasPath {
		path = "/" + path
	}

	switch {
	case strings.Contains(s.GatewayUrlTemplate, "{cid}"):
		return strings.ReplaceAll(s.GatewayUrlTemplate, "{cid}", cid) + path
	case strings.Contains(s.GatewayUrlTemplate, "%s"):
		return strings.Replace(s.GatewayUrlTemplate, "%s", cid, 1) + path
	default:
		return s.GatewayUrlTemplate + cid + path
	}
}
