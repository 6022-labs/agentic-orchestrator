package models_test

import (
	"encoding/json"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/models"
	"github.com/stretchr/testify/assert"
)

// decodeExternal mirrors the client: struct-decode the body, which is then mapped using those same bytes.
func decodeExternal(t *testing.T, body string) *models.ExternalWellKnownAgent6022 {
	t.Helper()

	var external models.ExternalWellKnownAgent6022
	assert.NoError(t, json.Unmarshal([]byte(body), &external))
	return &external
}
