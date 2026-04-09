package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenBuildingDecentralizedDomainNameTestingSuite struct {
	sut *services.DecentralizedDomainNameBuilder

	blockchainSettings                      *settings.BlockchainSettings
	decentralizedDomainNameProviderSettings *settings.DecentralizedDomainNameProviderSettings
}

func WhenBuildingDecentralizedDomainNameBeforeEach(t *testing.T) *WhenBuildingDecentralizedDomainNameTestingSuite {
	logger := zap.NewNop()
	blockchainSettings := &settings.BlockchainSettings{ChainId: 1}
	decentralizedDomainNameProviderSettings := &settings.DecentralizedDomainNameProviderSettings{DomainExtension: "eth"}

	sut := services.NewDecentralizedDomainNameBuilder(logger, blockchainSettings, decentralizedDomainNameProviderSettings)

	return &WhenBuildingDecentralizedDomainNameTestingSuite{
		sut:                                     sut,
		blockchainSettings:                      blockchainSettings,
		decentralizedDomainNameProviderSettings: decentralizedDomainNameProviderSettings,
	}
}

func TestWhenBuildingDecentralizedDomainName(t *testing.T) {
	t.Parallel()

	t.Run("Given valid agent name and collection name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should build the correct domain name", func(t *testing.T) {
			t.Parallel()
			basics := WhenBuildingDecentralizedDomainNameBeforeEach(t)

			agentName := "my-agent"
			collectionName := "my-collection"

			result := basics.sut.BuildDomainName(agentName, collectionName)

			expected := fmt.Sprintf("%s.%s.%d.protocol6022.%s",
				agentName,
				collectionName,
				basics.blockchainSettings.ChainId,
				basics.decentralizedDomainNameProviderSettings.DomainExtension,
			)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("Given a different chain ID", func(t *testing.T) {
		t.Parallel()

		t.Run("Should include the correct chain ID in the domain", func(t *testing.T) {
			t.Parallel()
			logger := zap.NewNop()
			blockchainSettings := &settings.BlockchainSettings{ChainId: 137}
			domainSettings := &settings.DecentralizedDomainNameProviderSettings{DomainExtension: "eth"}

			sut := services.NewDecentralizedDomainNameBuilder(logger, blockchainSettings, domainSettings)

			result := sut.BuildDomainName("agent", "collection")

			assert.Equal(t, "agent.collection.137.protocol6022.eth", result)
		})
	})

	t.Run("Given a custom domain extension", func(t *testing.T) {
		t.Parallel()

		t.Run("Should use the configured domain extension", func(t *testing.T) {
			t.Parallel()
			logger := zap.NewNop()
			blockchainSettings := &settings.BlockchainSettings{ChainId: 1}
			domainSettings := &settings.DecentralizedDomainNameProviderSettings{DomainExtension: "crypto"}

			sut := services.NewDecentralizedDomainNameBuilder(logger, blockchainSettings, domainSettings)

			result := sut.BuildDomainName("agent", "collection")

			assert.Equal(t, "agent.collection.1.protocol6022.crypto", result)
		})
	})
}
