package settings

import (
	"testing"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewDatabaseSettingsLoadsProvidedValues(t *testing.T) {
	k := koanf.New(".")
	_ = k.Load(confmap.Provider(map[string]interface{}{
		"database.url":            "postgres://localhost:5432/app",
		"database.max_open_conns": 40,
	}, "."), nil)

	settings := NewDatabaseSettings(zap.NewNop(), k)

	assert.Equal(t, "postgres://localhost:5432/app", settings.URL)
	assert.Equal(t, 40, settings.MaxOpenConns)
}

func TestNewDatabaseSettingsAppliesDefaultMaxOpenConns(t *testing.T) {
	k := koanf.New(".")
	_ = k.Load(confmap.Provider(map[string]interface{}{
		"database.url":            "postgres://localhost:5432/app",
		"database.max_open_conns": 0,
	}, "."), nil)

	settings := NewDatabaseSettings(zap.NewNop(), k)

	assert.Equal(t, "postgres://localhost:5432/app", settings.URL)
	assert.Equal(t, 20, settings.MaxOpenConns)
}
