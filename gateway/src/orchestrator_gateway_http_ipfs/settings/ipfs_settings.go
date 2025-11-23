package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type IPFSSettings struct {
	IpfsNodeHttpUrl string
}

func NewIPFSSettings(logger *zap.Logger) *IPFSSettings {
	ipfsNodeHttpUrl := os.Getenv("IPFS_NODE_HTTP_URL")
	if len(strings.TrimSpace(ipfsNodeHttpUrl)) == 0 {
		logger.Fatal("IPFS_NODE_HTTP_URL is not set")
	}

	return &IPFSSettings{
		IpfsNodeHttpUrl: ipfsNodeHttpUrl,
	}
}
