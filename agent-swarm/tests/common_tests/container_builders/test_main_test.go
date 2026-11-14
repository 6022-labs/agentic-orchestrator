package container_builders

import (
	"os"
	"testing"
)

// TestMain stops every shared container once the test process exits.
func TestMain(m *testing.M) {
	code := m.Run()
	GetPostgresContainerBuilder().stop()
	GetRedisContainerBuilder().stop()
	GetAnvilContainerBuilder().stop()
	os.Exit(code)
}
