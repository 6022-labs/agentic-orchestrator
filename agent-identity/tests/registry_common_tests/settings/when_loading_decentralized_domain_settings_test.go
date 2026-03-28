package settings_test

import (
	"testing"

	domain_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDecentralizedDomainNameProviderSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"decentralized_domain_name_provider.domain_extension": "eth",
		}, "."), nil)

		settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "eth", settings.DomainExtension)
	})

	t.Run("Loads settings with default domain extension when empty", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"decentralized_domain_name_provider.domain_extension": "",
		}, "."), nil)

		settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "eth", settings.DomainExtension)
	})

	t.Run("Loads settings with default domain extension when not set", func(t *testing.T) {
		k := koanf.New(".")

		settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "eth", settings.DomainExtension)
	})
}
