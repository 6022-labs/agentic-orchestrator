package container_builders

import (
	"os"
	"testing"
)

// TestMain stops every shared container once the test process exits.
func TestMain(m *testing.M) {
	code := m.Run()
	stopRedis()
	stopPostgres()
	os.Exit(code)
}
