package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type IpfsGatewaySettings struct {
	GatewayUrlTemplate string
}

func NewIpfsGatewaySettings(logger *zap.Logger) *IpfsGatewaySettings {
	ipfsGatewayUrlTemplate := os.Getenv("IPFS_GATEWAY_URL_TEMPLATE")
	if len(strings.TrimSpace(ipfsGatewayUrlTemplate)) == 0 {
		logger.Info("IPFS_GATEWAY_URL_TEMPLATE is not set, using default")
		ipfsGatewayUrlTemplate = "https://ipfs.io/ipfs/%s"
	}

	return &IpfsGatewaySettings{
		GatewayUrlTemplate: ipfsGatewayUrlTemplate,
	}
}
