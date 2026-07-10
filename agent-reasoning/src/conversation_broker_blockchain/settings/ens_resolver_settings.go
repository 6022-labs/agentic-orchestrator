package settings

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const EnsResolverSettingsKey = "ens_resolver"

// defaultEnsRegistryAddress is the canonical ENS registry, deployed at the
// same address on mainnet and the test networks (Sepolia, Holesky).
const defaultEnsRegistryAddress = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"

type EnsResolverSettings struct {
	RpcUrl                       string `koanf:"rpc_url"`
	RegistryAddress              string `koanf:"registry_address"`
	GatewayRequestTimeoutSeconds uint64 `koanf:"gateway_request_timeout_seconds"`
}

func NewEnsResolverSettings(logger *zap.Logger, k *koanf.Koanf) *EnsResolverSettings {
	settings := EnsResolverSettings{
		RegistryAddress:              defaultEnsRegistryAddress,
		GatewayRequestTimeoutSeconds: 10,
	}
	if err := k.Unmarshal(EnsResolverSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal ens resolver settings", zap.Error(err))
	}

	settings.RpcUrl = strings.TrimSpace(settings.RpcUrl)
	settings.RegistryAddress = strings.TrimSpace(settings.RegistryAddress)
	if settings.RegistryAddress == "" {
		settings.RegistryAddress = defaultEnsRegistryAddress
	}

	if len(settings.RpcUrl) == 0 {
		logger.Fatal("please set your ens_resolver.rpc_url value in your configuration")
	}
	if !common.IsHexAddress(settings.RegistryAddress) {
		logger.Fatal("please set a valid ens_resolver.registry_address value in your configuration")
	}
	if settings.GatewayRequestTimeoutSeconds == 0 {
		settings.GatewayRequestTimeoutSeconds = 10
	}

	return &settings
}

func (s *EnsResolverSettings) RegistryCommonAddress() common.Address {
	return common.HexToAddress(s.RegistryAddress)
}
