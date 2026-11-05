package jwt_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth/jwt"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/stretchr/testify/assert"
)

type WhenGeneratingRefreshTokenIdTestingSuite struct {
	sut *jwt.RefreshTokenManager
}

func WhenGeneratingRefreshTokenIdBeforeEach() *WhenGeneratingRefreshTokenIdTestingSuite {
	return &WhenGeneratingRefreshTokenIdTestingSuite{
		sut: jwt.NewRefreshTokenManager(&settings.JwtSettings{Hash: "test-signing-hash-at-least-32-characters"}),
	}
}

func TestWhenGeneratingRefreshTokenId(t *testing.T) {
	t.Parallel()

	t.Run("Given a working entropy source", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a 32-byte hex id", func(t *testing.T) {
			t.Parallel()

			suite := WhenGeneratingRefreshTokenIdBeforeEach()

			id, err := suite.sut.NewTokenId()

			assert.Nil(t, err)
			assert.Len(t, id, 64)
		})

		t.Run("Should return a different id on every call", func(t *testing.T) {
			t.Parallel()

			suite := WhenGeneratingRefreshTokenIdBeforeEach()

			first, _ := suite.sut.NewTokenId()
			second, _ := suite.sut.NewTokenId()

			assert.NotEqual(t, first, second)
		})
	})
}
