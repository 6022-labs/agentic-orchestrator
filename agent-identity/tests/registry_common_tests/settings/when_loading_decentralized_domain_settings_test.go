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
	t.Parallel()

	t.Run("Given a configured domain extension", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"decentralized_domain_name_provider.domain_extension": "eth",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "eth", settings.DomainExtension)
		})
	})

	t.Run("Given an empty domain extension", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"decentralized_domain_name_provider.domain_extension": "",
		}, "."), nil)

		t.Run("Should fall back to the default extension", func(t *testing.T) {
			t.Parallel()

			settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "eth", settings.DomainExtension)
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to the default extension", func(t *testing.T) {
			t.Parallel()

			settings := domain_settings.NewDecentralizedDomainNameProviderSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "eth", settings.DomainExtension)
		})
	})
}
