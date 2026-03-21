package services_test

import (
	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenBuildingDecentralizedDomainNameTestingSuite struct {
	blockchainSettings                      *settings.BlockchainSettings
	decentralizedDomainNameProviderSettings *settings.DecentralizedDomainNameProviderSettings

	sut *services.DecentralizedDomainNameBuilder
}

func WhenBuildingDecentralizedDomainNameBeforeEach() *WhenBuildingDecentralizedDomainNameTestingSuite {
	blockchainSettings := &settings.BlockchainSettings{
		ChainId: 1,
	}

	decentralizedDomainNameProviderSettings := &settings.DecentralizedDomainNameProviderSettings{
		DomainExtension: "eth",
	}

	sut := services.NewDecentralizedDomainNameBuilder(
		zap.NewNop(),
		blockchainSettings,
		decentralizedDomainNameProviderSettings,
	)

	return &WhenBuildingDecentralizedDomainNameTestingSuite{
		blockchainSettings:                      blockchainSettings,
		decentralizedDomainNameProviderSettings: decentralizedDomainNameProviderSettings,

		sut: sut,
	}
}

func TestWhenBuildingDecentralizedDomainName(t *testing.T) {
	t.Parallel()

	t.Run("Should return domain provider based on settings", func(t *testing.T) {
		t.Parallel()

		suite := WhenBuildingDecentralizedDomainNameBeforeEach()

		suite.decentralizedDomainNameProviderSettings.DomainExtension = "eth"
		firstDomain := suite.sut.BuildDomainName("john", "doe")
		assert.True(t, strings.HasSuffix(firstDomain, ".eth"))

		suite.decentralizedDomainNameProviderSettings.DomainExtension = "btc"
		secondDomain := suite.sut.BuildDomainName("john", "doe")
		assert.True(t, strings.HasSuffix(secondDomain, ".btc"))
	})

	t.Run("Should append domain name with chain id from settings", func(t *testing.T) {
		t.Parallel()

		suite := WhenBuildingDecentralizedDomainNameBeforeEach()

		suite.blockchainSettings.ChainId = 1
		firstDomain := suite.sut.BuildDomainName("john", "doe")
		assert.True(t, strings.Contains(firstDomain, ".1."))

		suite.blockchainSettings.ChainId = 2
		secondDomain := suite.sut.BuildDomainName("john", "doe")
		assert.True(t, strings.Contains(secondDomain, ".2."))
	})

	t.Run("Should be prefixed by agent name", func(t *testing.T) {
		t.Parallel()

		suite := WhenBuildingDecentralizedDomainNameBeforeEach()

		firstDomain := suite.sut.BuildDomainName("john", "doe")
		assert.True(t, strings.HasPrefix(firstDomain, "john."))

		secondDomain := suite.sut.BuildDomainName("jane", "doe")
		assert.True(t, strings.HasPrefix(secondDomain, "jane."))
	})

	t.Run("Should append domain name with agent collection name from parameters", func(t *testing.T) {
		t.Parallel()

		suite := WhenBuildingDecentralizedDomainNameBeforeEach()

		suite.blockchainSettings.ChainId = 1
		firstDomain := suite.sut.BuildDomainName("john", "doe")
		assert.True(t, strings.Contains(firstDomain, ".doe."))

		suite.blockchainSettings.ChainId = 2
		secondDomain := suite.sut.BuildDomainName("john", "dare")
		assert.True(t, strings.Contains(secondDomain, ".dare."))
	})
}
