package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapEnvEntries(t *testing.T) {
	mapped := mapEnvEntries([]string{
		"DATABASE__URL=postgres://localhost:5432/app",
		"HOST__API_PORT=3000",
		"NO_SEPARATOR=ignored",
		"INVALID_ENTRY",
	})

	assert.Equal(t, "postgres://localhost:5432/app", mapped["database.url"])
	assert.Equal(t, "3000", mapped["host.api_port"])
	_, exists := mapped["no_separator"]
	assert.False(t, exists)
}

func TestLoadKoanfConfigWithFileAndEnvOverride(t *testing.T) {
	t.Setenv("DATABASE__URL", "postgres://env-host:5432/app")

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "appsettings.json")
	content := `{"database":{"url":"postgres://file-host:5432/app","max_open_conns":10}}`
	require.NoError(t, os.WriteFile(configFile, []byte(content), 0o600))

	k, err := LoadKoanfConfig(configFile)
	require.NoError(t, err)

	assert.Equal(t, "postgres://env-host:5432/app", k.String("database.url"))
	assert.Equal(t, int64(10), k.Int64("database.max_open_conns"))
}
