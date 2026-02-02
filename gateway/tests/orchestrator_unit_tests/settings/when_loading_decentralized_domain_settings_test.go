package settings_test

import (
	"testing"

	domain_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDecentralizedDomainNameProviderSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("decentralized_domain_name_provider", map[string]any{
			"domain_extension": "eth",
		})

		settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "eth", settings.DomainExtension)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("DECENTRALIZED_DOMAIN__EXTENSION", "xyz")

		v := viper.New()
		settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "xyz", settings.DomainExtension)
	})
}
